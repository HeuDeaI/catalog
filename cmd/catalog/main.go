package main

import (
	"catalog/internal/handlers"
	"catalog/internal/models"
	"catalog/internal/repositories"
	"catalog/internal/services"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	const (
		DSN         = "host=localhost user=postgres password=postgres dbname=catalog port=5432 sslmode=disable"
		MinIOURL    = "localhost:9000"
		MinIOUser   = "admin"
		MinIOPass   = "password"
		MinIOBucket = "products"
	)

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	err = db.AutoMigrate(&models.Product{}, &models.Brand{})
	if err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	minioClient, err := minio.New(MinIOURL, &minio.Options{
		Creds:  credentials.NewStaticV4(MinIOUser, MinIOPass, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("MinIO connection error: %v", err)
	}

	skinTypeRepo := repositories.NewSkinTypeRepository(db)
	skinTypeRepo.SeedSkinTypes()

	productRepo := repositories.NewProductRepository(db, minioClient, MinIOBucket, MinIOURL)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	router := gin.Default()
	productHandler.RegisterRoutes(router)

	log.Println("Server running on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
