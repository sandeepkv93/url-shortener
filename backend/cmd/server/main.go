package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"url-shortener/internal/config"
	"url-shortener/internal/infrastructure/cache"
	"url-shortener/internal/infrastructure/database"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Ensure database is migrated
	if err := db.AutoMigrate(); err != nil {
		log.Printf("Warning: Failed to run auto migrations: %v", err)
	}

	// Connect to Redis
	redisAddr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	redisClient, err := cache.NewRedisClient(redisAddr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Create cache service
	cacheService := cache.NewCacheService(redisClient)

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Compress(5))
	r.Use(middleware.Timeout(60 * time.Second))

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		// Check database health
		dbErr := db.Health()
		
		// Check Redis health
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		redisErr := cacheService.Ping(ctx)
		
		// Determine overall health
		if dbErr != nil || redisErr != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			status := map[string]string{
				"status": "unhealthy",
			}
			if dbErr != nil {
				status["database"] = "down"
			} else {
				status["database"] = "up"
			}
			if redisErr != nil {
				status["redis"] = "down"
			} else {
				status["redis"] = "up"
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"status":"%s","database":"%s","redis":"%s"}`, 
				status["status"], status["database"], status["redis"])
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"healthy","database":"up","redis":"up"}`))
	})

	// Debug endpoints (development only)
	if cfg.IsDevelopment() {
		r.Get("/debug/db-stats", func(w http.ResponseWriter, r *http.Request) {
			stats := db.GetStats()
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "%+v", stats)
		})
		
		r.Get("/debug/redis-info", func(w http.ResponseWriter, r *http.Request) {
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
	}

	// Default route
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"URL Shortener API","version":"1.0.0"}`))
	})

	// Create server
	server := &http.Server{
		Addr:         cfg.GetServerAddress(),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s", cfg.GetServerAddress())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}