package main

import (
	"github.com/gin-gonic/gin"
	"github.com/klemis/packs-calculator/internal/handlers"
	"github.com/klemis/packs-calculator/internal/repositories"
	"github.com/klemis/packs-calculator/internal/services"
	"log"
)

func main() {
	log.Println("Starting API server...")
	// Initialize packs calculator service, database and repository.
	service, cleanup, err := initializePacksCalculatorService()
	if err != nil {
		log.Fatalf("failed to initialize services: %v", err)
	}
	defer cleanup()

	// Initialize handlers.
	handler := handlers.NewHandler(service)

	router := gin.Default()
	registerRoutes(router, handler)

	log.Println("API server listening on port 8080...")
	if err = router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

// registerRoutes sets up the API routes for the application.
func registerRoutes(router *gin.Engine, handlers *handlers.Handler) {
	// Serve static files (index.html, styles and js) from the "static" dir.
	router.Static("/static", "static")

	// API endpoints
	v1 := router.Group("/api/v1")
	{
		v1.POST("/packs", handlers.AddPackSize)
		v1.DELETE("/packs", handlers.DeletePackSize)
		v1.GET("/calculate", handlers.CalculatePacks)
	}
}

// initializePacksCalculatorService sets up the database and returns a new PacksCalculatorService instance.
func initializePacksCalculatorService() (services.PacksCalculator, func(), error) {
	// Initialize the database.
	db, cleanup, err := repositories.InitAndCloseDB()
	if err != nil {
		return nil, nil, err
	}

	// Initialize pack size repository and service.
	packSizeRepo := repositories.NewSQLPackSizeRepository(db)
	service := services.NewPacksCalculatorService(packSizeRepo)

	return service, cleanup, nil
}
