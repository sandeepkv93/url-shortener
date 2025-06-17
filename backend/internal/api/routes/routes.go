package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	
	"url-shortener/internal/api/handlers"
	"url-shortener/internal/api/middleware"
	"url-shortener/internal/core/ports"
)

type Config struct {
	// Handlers
	AuthHandler      *handlers.AuthHandler
	URLHandler       *handlers.URLHandler
	AnalyticsHandler *handlers.AnalyticsHandler
	QRHandler        *handlers.QRHandler
	
	// Middleware
	AuthMiddleware     *middleware.AuthMiddleware
	CORSMiddleware     *middleware.CORSMiddleware
	LoggingMiddleware  *middleware.LoggingMiddleware
	
	// Services (for rate limiting)
	CacheService ports.CacheService
	
	// Configuration
	EnableCORS   bool
	EnableLogging bool
	AllowedOrigins []string
}

type Router struct {
	config *Config
	chi    *chi.Mux
}

func NewRouter(config *Config) *Router {
	return &Router{
		config: config,
		chi:    chi.NewRouter(),
	}
}

func (r *Router) SetupRoutes() http.Handler {
	// Global middleware
	r.chi.Use(chimiddleware.Recoverer)
	r.chi.Use(chimiddleware.RequestID)
	r.chi.Use(chimiddleware.RealIP)
	
	// CORS middleware
	if r.config.EnableCORS {
		if r.config.CORSMiddleware != nil {
			r.chi.Use(r.config.CORSMiddleware.Handler)
		} else {
			// Use default CORS if no custom middleware provided
			r.chi.Use(middleware.CORS(r.config.AllowedOrigins...))
		}
	}
	
	// Logging middleware
	if r.config.EnableLogging {
		if r.config.LoggingMiddleware != nil {
			r.chi.Use(r.config.LoggingMiddleware.Handler)
		} else {
			// Use default logging if no custom middleware provided
			r.chi.Use(middleware.RequestLogging())
		}
	}
	
	// Global rate limiting
	if r.config.CacheService != nil {
		r.chi.Use(middleware.GlobalRateLimit(r.config.CacheService))
	}
	
	// Health check endpoint
	r.chi.Get("/health", r.healthCheck)
	
	// API versioning
	r.chi.Route("/api/v1", func(apiRouter chi.Router) {
		r.setupV1Routes(apiRouter)
	})
	
	// Short URL redirection (no API prefix)
	if r.config.URLHandler != nil {
		r.chi.Get("/{shortCode}", r.config.URLHandler.RedirectURL)
	}
	
	return r.chi
}

func (r *Router) setupV1Routes(apiRouter chi.Router) {
	// Authentication routes (no auth required)
	if r.config.AuthHandler != nil {
		apiRouter.Route("/auth", func(authRouter chi.Router) {
			// Rate limiting for auth endpoints
			if r.config.CacheService != nil {
				authRouter.Use(middleware.AuthRateLimit(r.config.CacheService))
			}
			
			authRouter.Post("/register", r.config.AuthHandler.Register)
			authRouter.Post("/login", r.config.AuthHandler.Login)
			authRouter.Post("/refresh", r.config.AuthHandler.RefreshToken)
			
			// Routes requiring authentication
			if r.config.AuthMiddleware != nil {
				authRouter.Group(func(protectedRouter chi.Router) {
					protectedRouter.Use(r.config.AuthMiddleware.RequireAuth)
					protectedRouter.Post("/logout", r.config.AuthHandler.Logout)
					protectedRouter.Get("/profile", r.config.AuthHandler.GetProfile)
					protectedRouter.Put("/profile", r.config.AuthHandler.UpdateProfile)
					protectedRouter.Post("/change-password", r.config.AuthHandler.ChangePassword)
					protectedRouter.Get("/validate", r.config.AuthHandler.ValidateToken)
				})
			}
		})
	}
	
	// URL management routes
	if r.config.URLHandler != nil {
		apiRouter.Route("/urls", func(urlRouter chi.Router) {
			// Public endpoints
			urlRouter.Get("/popular", r.config.URLHandler.GetPopularURLs)
			
			// Protected endpoints
			if r.config.AuthMiddleware != nil {
				urlRouter.Group(func(protectedRouter chi.Router) {
					protectedRouter.Use(r.config.AuthMiddleware.RequireAuth)
					
					// URL creation with rate limiting
					protectedRouter.Group(func(createRouter chi.Router) {
						if r.config.CacheService != nil {
							createRouter.Use(middleware.URLCreationRateLimit(r.config.CacheService))
						}
						createRouter.Post("/", r.config.URLHandler.CreateShortURL)
					})
					
					// URL management
					protectedRouter.Get("/", r.config.URLHandler.GetUserURLs)
					protectedRouter.Get("/{id}", r.config.URLHandler.GetURL)
					protectedRouter.Put("/{id}", r.config.URLHandler.UpdateURL)
					protectedRouter.Delete("/{id}", r.config.URLHandler.DeleteURL)
				})
			}
		})
	}
	
	// Analytics routes
	if r.config.AnalyticsHandler != nil && r.config.AuthMiddleware != nil {
		apiRouter.Route("/analytics", func(analyticsRouter chi.Router) {
			analyticsRouter.Use(r.config.AuthMiddleware.RequireAuth)
			
			// Dashboard analytics
			analyticsRouter.Get("/dashboard", r.config.AnalyticsHandler.GetDashboard)
			analyticsRouter.Get("/global", r.config.AnalyticsHandler.GetGlobalStats)
			analyticsRouter.Get("/top-urls", r.config.AnalyticsHandler.GetTopPerformingURLs)
			analyticsRouter.Get("/export", r.config.AnalyticsHandler.ExportAnalytics)
			
			// URL-specific analytics
			analyticsRouter.Route("/urls/{id}", func(urlAnalyticsRouter chi.Router) {
				urlAnalyticsRouter.Get("/", r.config.AnalyticsHandler.GetURLAnalytics)
				urlAnalyticsRouter.Get("/timeline", r.config.AnalyticsHandler.GetClickTimeline)
				urlAnalyticsRouter.Get("/geo", r.config.AnalyticsHandler.GetGeographicStats)
				urlAnalyticsRouter.Get("/devices", r.config.AnalyticsHandler.GetDeviceStats)
				urlAnalyticsRouter.Get("/referrers", r.config.AnalyticsHandler.GetReferrerStats)
			})
		})
	}
	
	// QR code routes
	if r.config.QRHandler != nil {
		apiRouter.Route("/qr", func(qrRouter chi.Router) {
			// Public endpoints
			qrRouter.Get("/formats", r.config.QRHandler.GetQRCodeFormats)
			qrRouter.Get("/sizes", r.config.QRHandler.GetQRCodeSizes)
			qrRouter.Post("/validate", r.config.QRHandler.ValidateQRCodeOptions)
			qrRouter.Post("/preview", r.config.QRHandler.GetQRCodePreview)
			
			// QR code generation (with optional auth)
			if r.config.AuthMiddleware != nil {
				qrRouter.Group(func(qrGenRouter chi.Router) {
					qrGenRouter.Use(r.config.AuthMiddleware.OptionalAuth)
					qrGenRouter.Post("/generate", r.config.QRHandler.GenerateQRCode)
					qrGenRouter.Get("/{shortCode}", r.config.QRHandler.GenerateQRCodeForURL)
				})
			}
		})
	}
	
	// Admin routes (future expansion)
	if r.config.AuthMiddleware != nil && r.config.AnalyticsHandler != nil {
		apiRouter.Route("/admin", func(adminRouter chi.Router) {
			adminRouter.Use(r.config.AuthMiddleware.RequireAuth)
			// adminRouter.Use(r.config.AuthMiddleware.AdminOnly) // Enable when admin functionality is added
			
			// Placeholder for admin endpoints
			adminRouter.Get("/stats", r.config.AnalyticsHandler.GetGlobalStats)
		})
	}
}

func (r *Router) healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{
		"status": "healthy",
		"service": "url-shortener",
		"version": "1.0.0",
		"timestamp": "` + req.Header.Get("Date") + `"
	}`))
}

// Route helpers for testing and debugging

func (r *Router) PrintRoutes() {
	chi.Walk(r.chi, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		// This can be used for debugging routes
		return nil
	})
}

func (r *Router) GetHandler() http.Handler {
	return r.chi
}

// Builder pattern for easier configuration

type RouterBuilder struct {
	config *Config
}

func NewRouterBuilder() *RouterBuilder {
	return &RouterBuilder{
		config: &Config{
			EnableCORS:     true,
			EnableLogging:  true,
			AllowedOrigins: []string{"*"},
		},
	}
}

func (b *RouterBuilder) WithAuthHandler(handler *handlers.AuthHandler) *RouterBuilder {
	b.config.AuthHandler = handler
	return b
}

func (b *RouterBuilder) WithURLHandler(handler *handlers.URLHandler) *RouterBuilder {
	b.config.URLHandler = handler
	return b
}

func (b *RouterBuilder) WithAnalyticsHandler(handler *handlers.AnalyticsHandler) *RouterBuilder {
	b.config.AnalyticsHandler = handler
	return b
}

func (b *RouterBuilder) WithQRHandler(handler *handlers.QRHandler) *RouterBuilder {
	b.config.QRHandler = handler
	return b
}

func (b *RouterBuilder) WithAuthMiddleware(middleware *middleware.AuthMiddleware) *RouterBuilder {
	b.config.AuthMiddleware = middleware
	return b
}

func (b *RouterBuilder) WithCORSMiddleware(middleware *middleware.CORSMiddleware) *RouterBuilder {
	b.config.CORSMiddleware = middleware
	return b
}

func (b *RouterBuilder) WithLoggingMiddleware(middleware *middleware.LoggingMiddleware) *RouterBuilder {
	b.config.LoggingMiddleware = middleware
	return b
}

func (b *RouterBuilder) WithCacheService(service ports.CacheService) *RouterBuilder {
	b.config.CacheService = service
	return b
}

func (b *RouterBuilder) WithCORS(enabled bool, origins ...string) *RouterBuilder {
	b.config.EnableCORS = enabled
	if len(origins) > 0 {
		b.config.AllowedOrigins = origins
	}
	return b
}

func (b *RouterBuilder) WithLogging(enabled bool) *RouterBuilder {
	b.config.EnableLogging = enabled
	return b
}

func (b *RouterBuilder) Build() *Router {
	return NewRouter(b.config)
}