package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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

	// === CORS ===
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
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

	routes.RegisterRoutes(r, database.DB)

	r.Static("/uploads", "./public/uploads")

	// --- DEBUGGING START (Boleh dihapus nanti kalau sudah fix) ---
	dir, _ := os.Getwd()
	fmt.Println("--------------------------------------------------")
	fmt.Println("📂 POSISI TERMINAL (Current Working Directory):")
	fmt.Println(dir)
	fmt.Println("--------------------------------------------------")

	uploadPath := "./public/uploads"
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		fmt.Println("❌ Folder Upload TIDAK DITEMUKAN di path:", uploadPath)
		absPath, _ := filepath.Abs(uploadPath)
		fmt.Println("   (Mencoba mencari di:", absPath, ")")
	} else {
		fmt.Println("✅ Folder Upload DITEMUKAN!")

		// Cek sampel file
		avatarPath := "./public/uploads/avatars"
		if _, err := os.Stat(avatarPath); !os.IsNotExist(err) {
			files, _ := os.ReadDir(avatarPath)
			if len(files) > 0 {
				fmt.Println("   - Contoh file:", files[0].Name())
			} else {
				fmt.Println("   - Folder avatars kosong")
			}
		}
	}
	fmt.Println("--------------------------------------------------")
	// --- DEBUGGING END ---

	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
