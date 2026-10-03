package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"orbit-backend-golang/internal/config"
	"orbit-backend-golang/internal/database"
	"orbit-backend-golang/internal/routes"
)

func main() {

	// Load application configuration
	cfg := config.LoadConfig()

	// Connect to PostgreSQL
	db, err := database.ConnectPostgres(cfg.DatabaseURL)

	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	defer db.Close()

	log.Println("PostgreSQL connected successfully!")

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := routes.MainRouter(db)

	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatal(err)
	}

	// Health check endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "API is running",
		})
	})

	// Start server
	log.Println("Server running on port: http://localhost:" + cfg.AppPort)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal("Server failed: ", err)
	}
}
