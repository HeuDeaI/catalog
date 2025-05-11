package main

import (
	"catalog/internal/handlers"
	"catalog/internal/models"
	"catalog/internal/repositories"
	"catalog/internal/routes"
	"catalog/internal/services"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	const DSN = "host=localhost user=postgres password=postgres dbname=catalog port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	err = db.AutoMigrate(&models.Product{}, &models.Brand{})
	if err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	router := gin.Default()
	routes.RegisterRoutes(router, productHandler)

	log.Println("Server running on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
