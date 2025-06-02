package main

import (
	"backend/internal/app"
	"backend/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("[WARN] Error loading .env file")
	}
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("[DB] Connection error: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("[DB] Migrations error: %v", err)
	}
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			_, err := db.Exec(`
			UPDATE users
			SET is_online = FALSE
			WHERE is_online = TRUE AND last_seen < NOW() - INTERVAL '15 seconds'
		`)
			if err != nil {
				//log.Printf("Auto-offline update failed: %v", err)
			} else {
				//log.Println("Auto-offline check completed")
			}
		}
	}()
	router := gin.Default()
	app.SetupMiddleware(router)
	app.SetupRoutes(router, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("[APP] Server listening on " + port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("[APP] Server failed: %v", err)
	}
}
