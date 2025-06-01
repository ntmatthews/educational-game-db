package server

import (
	"fmt"
	"log"
	"net/http"

	"educational-game-db/internal/database"
	"educational-game-db/internal/handlers"
	"github.com/gin-gonic/gin"
)

type Server struct {
	db     *database.Database
	router *gin.Engine
	port   string
}

func NewServer(db *database.Database, port string) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	
	server := &Server{
		db:     db,
		router: router,
		port:   port,
	}

	server.setupRoutes()
	return server
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
