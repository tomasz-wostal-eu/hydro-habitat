package api

import (
	"hydro-habitat/backend/store"
	// Import docs package with explicit name to use it for initialization
	docs "hydro-habitat/backend/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(db *store.DB) *gin.Engine {
	router := gin.Default()
	
	// Programmatically set Swagger info
	docs.SwaggerInfo.Title = "Hydro Habitat API"
	docs.SwaggerInfo.Description = "This is the API for the Hydro Habitat application."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Prosta polityka CORS na potrzeby deweloperskie
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check endpoint
	router.GET("/health", HealthCheck)

	// Initialize Swagger docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiV1 := router.Group("/api/v1")
	tankStore := store.NewTankStore(db)
	tankHandler := NewTankHandler(tankStore)

	tanks := apiV1.Group("/tanks")
	{
		tanks.POST("", tankHandler.CreateTank)
		tanks.GET("", tankHandler.GetAllTanks)
		tanks.GET("/:id", tankHandler.GetTankByID)
		tanks.PUT("/:id", tankHandler.UpdateTank)
		tanks.DELETE("/:id", tankHandler.DeleteTank)
	}

	return router
}
