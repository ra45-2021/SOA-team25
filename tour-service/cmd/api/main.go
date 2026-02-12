package main

import (
	"context"
	"log"
	"os"
	"tour-service/internal/client"
	"tour-service/internal/handler"
	"tour-service/internal/model"
	"tour-service/internal/service"
	"tour-service/internal/store"
	_"tour-service/pb"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("MYSQL_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to mysql: %v", err)
	}
	db.AutoMigrate(&model.Tour{}, &model.Checkpoint{})

	minioClient, err := minio.New("minio:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("failed to connect to minio: %v", err)
	}

	bucket := "checkpoint-images"
	if exists, _ := minioClient.BucketExists(context.Background(), bucket); !exists {
		minioClient.MakeBucket(context.Background(), bucket, minio.MakeBucketOptions{})
	}

	authGRPCAddr := os.Getenv("AUTH_GRPC_ADDR")
	authClient, err := client.NewAuthClient(authGRPCAddr)
	if err != nil {
		log.Fatalf("failed to connect to auth grpc: %v", err)
	}

	tourStore := store.NewTourStore(db)
	tourSvc := service.NewTourService(tourStore, authClient)
	tourHandler := handler.NewTourHandler(tourSvc, minioClient, bucket)

	r := gin.Default()
	
	tours := r.Group("/tours")
	{
		tours.POST("", tourHandler.CreateTour)
		tours.GET("/my", tourHandler.GetMyTours)
		tours.POST("/:id/checkpoints", tourHandler.AddCheckpoint)
	}

	port := ":8080"
	log.Printf("Tour service running on %s", port)
	log.Fatal(r.Run(port))
}