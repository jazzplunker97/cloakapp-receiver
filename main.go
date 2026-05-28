package main

import (
	"log"
	"net/http"
	"os"

	"cloakapp-receiver/db"
	"cloakapp-receiver/handlers"
	"github.com/gin-gonic/gin"

	_ "cloakapp-receiver/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title CloakApp Telemetry Receiver API
// @version 1.0
// @description High-performance telemetry receiver with InfluxDB storage.
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key

func main() {
	// Initialize Database
	influx, err := db.NewInfluxDB()
	if err != nil {
		log.Printf("Warning: Could not connect to InfluxDB: %v", err)
	} else {
		defer influx.Close()
	}

	// Initialize KeyStore
	ks, err := db.NewKeyStore("keys.json")
	if err != nil {
		log.Fatalf("Failed to initialize keystore: %v", err)
	}

	r := gin.Default()

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-API-Key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	handler := handlers.NewTelemetryHandler(influx)
	adminHandler := handlers.NewAdminHandler(influx, ks)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/telemetry", handler.Receive)

	// Auth Middleware for protected API
	authMiddleware := func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		if key == "" || !ks.Validate(key) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing API Key"})
			c.Abort()
			return
		}
		c.Next()
	}

	admin := r.Group("/admin")
	{
		admin.GET("/", adminHandler.Dashboard)
		admin.GET("/api/data", adminHandler.GetData)
		admin.GET("/api/keys", adminHandler.ListKeys)
		admin.POST("/api/keys", adminHandler.CreateKey)
		admin.DELETE("/api/keys/:key", adminHandler.DeleteKey)
	}

	// External Data API (Protected)
	api := r.Group("/api")
	api.Use(authMiddleware)
	{
		api.GET("/telemetry", adminHandler.GetData)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
