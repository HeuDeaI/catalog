package main

import (
	"catalog/internal/handlers"
	"catalog/internal/models"
	"catalog/internal/repositories"
	"catalog/internal/services"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dsn = "host=localhost user=postgres password=postgres dbname=catalog port=5432 sslmode=disable"

func main() {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
	api := router.Group("/products")
	{
		api.POST("/", productHandler.CreateProduct)
		api.GET("/", productHandler.GetAllProducts)
		api.GET("/:id", productHandler.GetProductByID)
		api.PUT("/:id", productHandler.UpdateProduct)
		api.DELETE("/:id", productHandler.DeleteProduct)
	}

	log.Println("Server running on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
