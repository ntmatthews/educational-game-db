package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"educational-game-db/internal/database"
	"educational-game-db/internal/handlers"
	"educational-game-db/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	db          *database.Database
	router      *gin.Engine
	port        string
	rateLimiter *middleware.RateLimiter
}

func NewServer(db *database.Database, port string) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	server := &Server{
		db:          db,
		router:      router,
		port:        port,
		rateLimiter: middleware.NewRateLimiter(),
	}

	server.setupMiddleware()
	server.setupRoutes()
	return server
}

func (s *Server) setupMiddleware() {
	// CORS configuration
	config := cors.Config{
		AllowOrigins:     []string{"*"}, // In production, specify allowed origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-API-Key", "X-Request-ID"},
		ExposeHeaders:    []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	s.router.Use(cors.New(config))

	// Security middleware
	s.router.Use(middleware.SecurityHeaders())
	s.router.Use(middleware.RequestIDMiddleware())
	s.router.Use(middleware.LoggingMiddleware())

	// Rate limiting - 100 requests per minute per IP
	s.router.Use(s.rateLimiter.RateLimit(100, 20))
}

func (s *Server) setupRoutes() {
	handler := handlers.NewHandler(s.db)

	// Serve static files
	s.router.Static("/static", "./web/static")
	s.router.LoadHTMLGlob("web/templates/*")

	// Web routes
	s.router.GET("/", s.serveIndex)
	s.router.GET("/admin", s.serveAdmin)

	// PWA routes
	s.router.GET("/manifest.json", func(c *gin.Context) {
		c.File("./web/static/manifest.json")
	})
	s.router.GET("/sw.js", func(c *gin.Context) {
		c.File("./web/static/sw.js")
	})

	// API routes
	api := s.router.Group("/api")
	{
		api.GET("/accounts", handler.GetAccounts)
		api.POST("/accounts", handler.CreateAccount)
		api.GET("/accounts/:id", handler.GetAccount)
		api.PUT("/accounts/:id", handler.UpdateAccount)
		api.DELETE("/accounts/:id", handler.DeleteAccount)
		api.GET("/stats", handler.GetStats)
		api.POST("/login", handler.Login)
		
		// Export/Import routes (with stricter rate limiting)
		exportGroup := api.Group("/export")
		exportGroup.Use(s.rateLimiter.RateLimit(10, 2)) // More restrictive for export/import
		{
			exportGroup.GET("/csv", handler.ExportCSV)
			exportGroup.GET("/json", handler.ExportJSON)
			exportGroup.POST("/csv", handler.ImportCSV)
		}
	}

	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
}

func (s *Server) serveIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Educational Game Database",
	})
}

func (s *Server) serveAdmin(c *gin.Context) {
	c.HTML(http.StatusOK, "admin.html", gin.H{
		"title": "Admin Dashboard",
	})
}

func (s *Server) Start() error {
	log.Printf("Starting server on port %s", s.port)
	log.Printf("Access the web interface at: http://localhost:%s", s.port)
	log.Printf("Admin dashboard at: http://localhost:%s/admin", s.port)
	return s.router.Run(fmt.Sprintf(":%s", s.port))
}

func (s *Server) Stop() error {
	return s.db.Close()
}
