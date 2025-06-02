package main

import (
	"catalog/internal/handlers"
	"catalog/internal/models"
	"catalog/internal/repositories"
	"catalog/internal/services"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	const DSN = "host=localhost user=postgres password=postgres dbname=catalog port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	err = db.AutoMigrate(&models.Product{}, &models.Brand{}, &models.SkinType{})
	if err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	skinTypeRepo := repositories.NewSkinTypeRepository(db)
	skinTypeRepo.SeedSkinTypes()

	productRepo := repositories.NewProductRepository(db)
	imageRepo := repositories.NewProductImageRepository()
	productService := services.NewProductService(productRepo, imageRepo)
	productHandler := handlers.NewProductHandler(productService)

	router := gin.Default()
	productHandler.RegisterRoutes(router)

	router.Static("/docs", "./docs")
	url := ginSwagger.URL("http://localhost:8080/docs/swagger.yaml")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	log.Println("Server running on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
