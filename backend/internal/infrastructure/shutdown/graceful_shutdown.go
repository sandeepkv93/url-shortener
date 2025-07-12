package shutdown

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"url-shortener/internal/config"
	"url-shortener/internal/core/services"
)

// GracefulShutdown manages the graceful shutdown of all application components
type GracefulShutdown struct {
	logger        *services.LoggingService
	config        *config.Config
	server        *http.Server
	shutdownHooks []ShutdownHook
	mutex         sync.RWMutex
	isShuttingDown bool
}

// ShutdownHook represents a function that should be called during shutdown
type ShutdownHook struct {
	Name     string
	Priority int // Lower numbers have higher priority
	Timeout  time.Duration
	Func     func(ctx context.Context) error
}

// NewGracefulShutdown creates a new graceful shutdown manager
func NewGracefulShutdown(cfg *config.Config, logger *services.LoggingService) *GracefulShutdown {
	return &GracefulShutdown{
		logger:        logger,
		config:        cfg,
		shutdownHooks: make([]ShutdownHook, 0),
	}
}

// SetServer sets the HTTP server to be shutdown
func (gs *GracefulShutdown) SetServer(server *http.Server) {
	gs.server = server
}

// AddShutdownHook adds a function to be called during shutdown
func (gs *GracefulShutdown) AddShutdownHook(name string, priority int, timeout time.Duration, fn func(ctx context.Context) error) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	
	if gs.isShuttingDown {
		gs.logger.Warn("Attempted to add shutdown hook during shutdown", "name", name)
		return
	}

	hook := ShutdownHook{
		Name:     name,
		Priority: priority,
		Timeout:  timeout,
		Func:     fn,
	}

	gs.shutdownHooks = append(gs.shutdownHooks, hook)
	gs.logger.Debug("Added shutdown hook", "name", name, "priority", priority)
}

// RegisterDatabaseShutdown registers database shutdown hooks
func (gs *GracefulShutdown) RegisterDatabaseShutdown(dbCloser func() error) {
	gs.AddShutdownHook("database", 100, 30*time.Second, func(ctx context.Context) error {
		gs.logger.Info("Closing database connections...")
		if err := dbCloser(); err != nil {
			gs.logger.Error("Error closing database", "error", err.Error())
			return err
		}
		gs.logger.Info("Database connections closed successfully")
		return nil
	})
}

// RegisterRedisShutdown registers Redis shutdown hooks
func (gs *GracefulShutdown) RegisterRedisShutdown(redisCloser func() error) {
	gs.AddShutdownHook("redis", 90, 15*time.Second, func(ctx context.Context) error {
		gs.logger.Info("Closing Redis connections...")
		if err := redisCloser(); err != nil {
			gs.logger.Error("Error closing Redis", "error", err.Error())
			return err
		}
		gs.logger.Info("Redis connections closed successfully")
		return nil
	})
}

// RegisterBackgroundServiceShutdown registers background service shutdown hooks
func (gs *GracefulShutdown) RegisterBackgroundServiceShutdown(serviceName string, serviceCloser func() error) {
	gs.AddShutdownHook(serviceName, 80, 20*time.Second, func(ctx context.Context) error {
		gs.logger.Info("Stopping background service", "service", serviceName)
		if err := serviceCloser(); err != nil {
			gs.logger.Error("Error stopping background service", "service", serviceName, "error", err.Error())
			return err
		}
		gs.logger.Info("Background service stopped successfully", "service", serviceName)
		return nil
	})
}

// RegisterCustomShutdown registers a custom shutdown hook
func (gs *GracefulShutdown) RegisterCustomShutdown(name string, priority int, timeout time.Duration, fn func() error) {
	gs.AddShutdownHook(name, priority, timeout, func(ctx context.Context) error {
		return fn()
	})
}

// WaitForShutdown waits for shutdown signals and initiates graceful shutdown
func (gs *GracefulShutdown) WaitForShutdown() {
	// Create channel to receive OS signals
	signalChan := make(chan os.Signal, 1)
	
	// Register signals to catch
	signal.Notify(signalChan, 
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGTERM, // Termination signal
		syscall.SIGQUIT, // Quit signal
		syscall.SIGHUP,  // Hangup signal
	)

	// Wait for signal
	sig := <-signalChan
	gs.logger.Info("Received shutdown signal", "signal", sig.String())

	// Start graceful shutdown
	if err := gs.Shutdown(); err != nil {
		gs.logger.Error("Graceful shutdown failed", "error", err.Error())
		os.Exit(1)
	}

	gs.logger.Info("Graceful shutdown completed successfully")
}

// Shutdown performs the graceful shutdown
func (gs *GracefulShutdown) Shutdown() error {
	gs.mutex.Lock()
	if gs.isShuttingDown {
		gs.mutex.Unlock()
		return fmt.Errorf("shutdown already in progress")
	}
	gs.isShuttingDown = true
	gs.mutex.Unlock()

	gs.logger.Info("Starting graceful shutdown process...")
	
	startTime := time.Now()
	shutdownTimeout := gs.config.GetGracefulShutdownTimeout()

	// Create context with timeout for the entire shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Shutdown HTTP server first to stop accepting new requests
	if gs.server != nil {
		gs.logger.Info("Shutting down HTTP server...")
		
		serverCtx, serverCancel := context.WithTimeout(ctx, 30*time.Second)
		defer serverCancel()
		
		if err := gs.server.Shutdown(serverCtx); err != nil {
			gs.logger.Error("Error shutting down HTTP server", "error", err.Error())
			// Continue with other shutdown tasks even if server shutdown fails
		} else {
			gs.logger.Info("HTTP server shutdown successfully")
		}
	}

	// Sort shutdown hooks by priority (lower numbers first)
	gs.sortShutdownHooks()

	// Execute shutdown hooks
	for _, hook := range gs.shutdownHooks {
		gs.logger.Info("Executing shutdown hook", "name", hook.Name, "priority", hook.Priority)
		
		hookCtx, hookCancel := context.WithTimeout(ctx, hook.Timeout)
		
		// Execute hook with timeout
		done := make(chan error, 1)
		go func() {
			done <- hook.Func(hookCtx)
		}()

		select {
		case err := <-done:
			hookCancel()
			if err != nil {
				gs.logger.Error("Shutdown hook failed", "name", hook.Name, "error", err.Error())
				// Continue with other hooks even if one fails
			} else {
				gs.logger.Info("Shutdown hook completed successfully", "name", hook.Name)
			}
		case <-hookCtx.Done():
			hookCancel()
			gs.logger.Warn("Shutdown hook timed out", "name", hook.Name, "timeout", hook.Timeout)
			// Continue with other hooks even if one times out
		}

		// Check if main shutdown context is cancelled
		select {
		case <-ctx.Done():
			gs.logger.Warn("Shutdown process timed out, forcing exit")
			return fmt.Errorf("shutdown process timed out after %v", shutdownTimeout)
		default:
			// Continue with next hook
		}
	}

	shutdownDuration := time.Since(startTime)
	gs.logger.Info("Graceful shutdown process completed", 
		"duration", shutdownDuration.String(),
		"timeout", shutdownTimeout.String(),
	)

	return nil
}

// ForceShutdown forces immediate shutdown without graceful cleanup
func (gs *GracefulShutdown) ForceShutdown() {
	gs.logger.Warn("Forcing immediate shutdown...")
	
	if gs.server != nil {
		if err := gs.server.Close(); err != nil {
			gs.logger.Error("Error during forced server shutdown", "error", err.Error())
		}
	}
	
	os.Exit(1)
}

// IsShuttingDown returns true if shutdown is in progress
func (gs *GracefulShutdown) IsShuttingDown() bool {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.isShuttingDown
}

// GetShutdownStatus returns current shutdown status information
func (gs *GracefulShutdown) GetShutdownStatus() map[string]interface{} {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	
	hookNames := make([]string, len(gs.shutdownHooks))
	for i, hook := range gs.shutdownHooks {
		hookNames[i] = hook.Name
	}
	
	return map[string]interface{}{
		"is_shutting_down": gs.isShuttingDown,
		"registered_hooks": hookNames,
		"hook_count":      len(gs.shutdownHooks),
		"shutdown_timeout": gs.config.GetGracefulShutdownTimeout().String(),
	}
}

// sortShutdownHooks sorts hooks by priority (lower numbers first)
func (gs *GracefulShutdown) sortShutdownHooks() {
	// Simple bubble sort for small arrays
	n := len(gs.shutdownHooks)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if gs.shutdownHooks[j].Priority > gs.shutdownHooks[j+1].Priority {
				gs.shutdownHooks[j], gs.shutdownHooks[j+1] = gs.shutdownHooks[j+1], gs.shutdownHooks[j]
			}
		}
	}
}

// HealthCheckHandler returns a handler for health checks that considers shutdown status
func (gs *GracefulShutdown) HealthCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if gs.IsShuttingDown() {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"shutting_down","message":"Server is shutting down"}`))
			return
		}
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","message":"Server is running"}`))
	}
}

// ReadinessHandler returns a handler for readiness checks that considers shutdown status
func (gs *GracefulShutdown) ReadinessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if gs.IsShuttingDown() {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"not_ready","message":"Server is shutting down"}`))
			return
		}
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ready","message":"Server is ready to accept requests"}`))
	}
}

// LivenessHandler returns a handler for liveness checks
func (gs *GracefulShutdown) LivenessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Liveness check should only fail if the process is completely unresponsive
		// Even during shutdown, the process is still "alive"
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"alive","message":"Server process is alive"}`))
	}
}

// SetupShutdownHandlers adds shutdown-related handlers to a router
func (gs *GracefulShutdown) SetupShutdownHandlers(router interface{}) {
	// This would be implemented based on the specific router type
	// For now, we'll provide the handlers as methods that can be manually registered
}

// ShutdownMetrics provides metrics about the shutdown process
type ShutdownMetrics struct {
	StartTime        time.Time `json:"start_time"`
	Duration         string    `json:"duration"`
	HooksExecuted    int       `json:"hooks_executed"`
	HooksFailed      int       `json:"hooks_failed"`
	HooksTimedOut    int       `json:"hooks_timed_out"`
	ShutdownComplete bool      `json:"shutdown_complete"`
}

// GetShutdownMetrics returns metrics about the shutdown process
func (gs *GracefulShutdown) GetShutdownMetrics() *ShutdownMetrics {
	// This would track actual metrics during shutdown
	// For now, return basic structure
	return &ShutdownMetrics{
		StartTime:        time.Now(),
		Duration:         "0s",
		HooksExecuted:    len(gs.shutdownHooks),
		HooksFailed:      0,
		HooksTimedOut:    0,
		ShutdownComplete: false,
	}
}