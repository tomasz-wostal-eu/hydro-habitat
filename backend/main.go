package main

import (
	"hydro-habitat/backend/api"
	"hydro-habitat/backend/config"
	"hydro-habitat/backend/store"
	"log"

	_ "hydro-habitat/backend/docs" // Import generated docs

	"github.com/gin-gonic/gin"
)

// @title Hydro Habitat API
// @version 1.0
// @description This is the API for the Hydro Habitat application.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	gin.SetMode(cfg.GinMode)

	db, err := store.NewPostgresStore(cfg)
	if err != nil {
		log.Fatalf("could not connect to the database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	router := api.NewRouter(db)

	log.Printf("Server starting on port %s", cfg.APIPort)
	if err := router.Run(":" + cfg.APIPort); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
