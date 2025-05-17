package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	MinIOURL    = "localhost:9000"
	MinIOUser   = "admin"
	MinIOPass   = "Jbjslxmjs"
	MinIOBucket = "products"
)

func createMinioClient() (*minio.Client, error) {
	minioClient, err := minio.New(MinIOURL, &minio.Options{
		Creds:  credentials.NewStaticV4(MinIOUser, MinIOPass, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	return minioClient, nil
}

type ProductImageRepository interface {
	UploadImage(localFilePath, objectName string) (string, error)
	DeleteImage(objectName string) error
	UpdateImage(oldObjectName, localFilePath, objectName string) (string, error)
}

type productImageRepository struct {
	minioClient *minio.Client
	bucket      string
	minioURL    string
}

func NewProductImageRepository() ProductImageRepository {
	client, err := createMinioClient()
	if err != nil {
		log.Fatalf("Error creating the MinIO client: %v", err)
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, MinIOBucket)
	if err != nil {
		log.Fatalf("Error checking the existence of the bucket: %v", err)
	}
	if !exists {
		if err = client.MakeBucket(ctx, MinIOBucket, minio.MakeBucketOptions{}); err != nil {
			log.Fatalf("Error creating bucket: %v", err)
		}
	}

	return &productImageRepository{
		minioClient: client,
		bucket:      MinIOBucket,
		minioURL:    MinIOURL,
	}
}

func (r *productImageRepository) UploadImage(localFilePath string, objectName string) (string, error) {
	_, err := r.minioClient.FPutObject(context.Background(), r.bucket, objectName, localFilePath, minio.PutObjectOptions{ContentType: "image/png"})
	if err != nil {
		return "", err
	}
	imageURL := fmt.Sprintf("%s/%s/%s", r.minioURL, r.bucket, objectName)
	return imageURL, nil
}

func (r *productImageRepository) DeleteImage(objectName string) error {
	return r.minioClient.RemoveObject(context.Background(), r.bucket, objectName, minio.RemoveObjectOptions{})
}

func (r *productImageRepository) UpdateImage(oldObjectName string, localFilePath string, objectName string) (string, error) {
	if oldObjectName != "" {
		if err := r.DeleteImage(oldObjectName); err != nil {
			return "", fmt.Errorf("couldn't delete old image: %v", err)
		}
	}
	return r.UploadImage(localFilePath, objectName)
}
