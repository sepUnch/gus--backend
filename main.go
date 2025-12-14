package main

import (
	"log"

	"github.com/Zain0205/gdgoc-subbmission-be-go/config"
	"github.com/Zain0205/gdgoc-subbmission-be-go/database"
	"github.com/Zain0205/gdgoc-subbmission-be-go/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	database.ConnectDatabase()
	database.MigrateDatabase()
	database.SeedAchievementTypes(database.DB)

	r := gin.Default()

	// === CORS HARUS DI ATAS ROUTES ===
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Izinkan frontend lokal
		if origin == "http://localhost:3000" || origin == "http://localhost:8000" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// ROUTES DITEMPATKAN SETELAH CORS
	routes.RegisterRoutes(r, database.DB)

	r.Static("/uploads", "./public/uploads")
	
	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
