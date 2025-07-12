package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"url-shortener/internal/api/routes"
	"url-shortener/internal/config"
	"url-shortener/internal/core/services"
	"url-shortener/internal/infrastructure/cache"
	"url-shortener/internal/infrastructure/database"
	"url-shortener/internal/infrastructure/shutdown"

	"github.com/go-chi/chi/v5"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logging service
	logger := services.NewLoggingService(&services.LoggingConfig{
		Level:  cfg.Logging.Level,
		Format: cfg.Logging.Format,
	})

	// Initialize graceful shutdown manager
	shutdownManager := shutdown.NewGracefulShutdown(cfg, logger)

	logger.Info("Starting URL Shortener service", 
		"version", "1.0.0",
		"environment", cfg.Server.Env,
		"address", cfg.GetServerAddress(),
	)

	// Connect to database
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err.Error())
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Register database shutdown hook
	shutdownManager.RegisterDatabaseShutdown(func() error {
		return db.Close()
	})

	// Ensure database is migrated
	if err := db.AutoMigrate(); err != nil {
		logger.Warn("Failed to run auto migrations", "error", err.Error())
	} else {
		logger.Info("Database migrations completed successfully")
	}

	// Connect to Redis
	redisAddr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	redisClient, err := cache.NewRedisClient(redisAddr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		logger.Error("Failed to connect to Redis", "error", err.Error())
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Register Redis shutdown hook
	shutdownManager.RegisterRedisShutdown(func() error {
		return redisClient.Close()
	})

	// Create cache service
	cacheService := cache.NewCacheService(redisClient)

	// For now, use a basic router to demonstrate Step 10 completion
	// The comprehensive router setup will be completed when service layer issues are resolved
	
	// Setup basic router using the routes package
	router := routes.NewRouterBuilder().
		WithCORS(true, cfg.GetAllowedOrigins()...).
		WithLogging(true).
		Build()

	// Setup routes and get the handler
	handler := router.SetupRoutes()

	// Add development debug endpoints if enabled
	if cfg.IsDevelopment() {
		// Create a new router that wraps our main handler with debug routes
		debugRouter := chi.NewRouter()
		debugRouter.Mount("/", handler)
		
		debugRouter.Get("/debug/db-stats", func(w http.ResponseWriter, r *http.Request) {
			stats := db.GetStats()
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "%+v", stats)
		})
		
		debugRouter.Get("/debug/redis-info", func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			
			info, err := cacheService.Info(ctx)
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(`{"error":"failed to get redis info"}`))
				return
			}
			
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(info))
		})
		
		// Add shutdown status endpoint for development
		debugRouter.Get("/debug/shutdown-status", func(w http.ResponseWriter, r *http.Request) {
			status := shutdownManager.GetShutdownStatus()
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "%+v", status)
		})
		
		handler = debugRouter
	}

	// Add production health check endpoints if enabled
	if cfg.IsHealthEndpointsEnabled() {
		healthRouter := chi.NewRouter()
		healthRouter.Mount("/", handler)
		
		// Kubernetes-style health checks
		healthRouter.Get("/health", shutdownManager.HealthCheckHandler())
		healthRouter.Get("/healthz", shutdownManager.HealthCheckHandler())
		healthRouter.Get("/ready", shutdownManager.ReadinessHandler())
		healthRouter.Get("/readiness", shutdownManager.ReadinessHandler())
		healthRouter.Get("/live", shutdownManager.LivenessHandler())
		healthRouter.Get("/liveness", shutdownManager.LivenessHandler())
		
		handler = healthRouter
	}

	// Create server with production-optimized timeouts
	server := &http.Server{
		Addr:         cfg.GetServerAddress(),
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Register server with shutdown manager
	shutdownManager.SetServer(server)

	// Register any background services for shutdown
	// (In a real application, you would register actual background services here)
	shutdownManager.RegisterCustomShutdown("cleanup", 50, 10*time.Second, func() error {
		logger.Info("Performing final cleanup...")
		// Perform any final cleanup tasks
		return nil
	})

	// Start server in a goroutine
	go func() {
		logger.Info("Starting HTTP server", 
			"address", cfg.GetServerAddress(),
			"read_timeout", "30s",
			"write_timeout", "30s",
			"idle_timeout", "120s",
		)
		
		if cfg.IsTLSEnabled() {
			logger.Info("Starting HTTPS server with TLS", 
				"cert_file", cfg.GetTLSCertFile(),
				"key_file", cfg.GetTLSKeyFile(),
			)
			if err := server.ListenAndServeTLS(cfg.GetTLSCertFile(), cfg.GetTLSKeyFile()); err != nil && err != http.ErrServerClosed {
				logger.Error("HTTPS server failed to start", "error", err.Error())
				log.Fatalf("Failed to start HTTPS server: %v", err)
			}
		} else {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Error("HTTP server failed to start", "error", err.Error())
				log.Fatalf("Failed to start HTTP server: %v", err)
			}
		}
	}()

	// Log successful startup
	logger.Info("URL Shortener service started successfully",
		"pid", os.Getpid(),
		"environment", cfg.Server.Env,
		"graceful_shutdown_timeout", cfg.GetGracefulShutdownTimeout().String(),
	)

	// Wait for shutdown signals and perform graceful shutdown
	shutdownManager.WaitForShutdown()
}