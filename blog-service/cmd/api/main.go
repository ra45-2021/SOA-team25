package main

import (
	"log"
	"strings"
	"context"
    "fmt"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"blog-service/internal/config"
	"blog-service/internal/db"
	"blog-service/internal/handler"
	"blog-service/internal/middleware"
	"blog-service/internal/s3"
	"blog-service/internal/service"
	"blog-service/internal/store"
)

func main() {
	r := gin.Default()
	r.Use(middleware.CORS())

	secret := []byte(config.MustEnv("JWT_SECRET"))
	dsn := config.MustEnv("MYSQL_DSN")

	sqlDB := db.MustMySQL(dsn)
	db.MustInitSchema(sqlDB)

	blogStore := store.NewMySQLBlogStore(sqlDB)

	authBase := config.MustEnv("AUTH_BASE_URL")
	authClient := service.NewAuthClient(authBase)

	s3Endpoint := config.MustEnv("S3_ENDPOINT")
	s3Access := config.MustEnv("S3_ACCESS_KEY")
	s3Secret := config.MustEnv("S3_SECRET_KEY")
	bucket := config.MustEnv("S3_BUCKET")
	publicBase := config.MustEnv("S3_PUBLIC_BASE_URL")

	endpoint := strings.TrimPrefix(strings.TrimPrefix(s3Endpoint, "http://"), "https://")
	useSSL := strings.HasPrefix(s3Endpoint, "https://")

	mc, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3Access, s3Secret, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatal(err)
	}

	uploader := s3.NewUploader(mc, bucket, publicBase)

	ctx := context.Background()
	exists, err := mc.BucketExists(ctx, bucket)
	if err != nil {
		log.Printf("Error checking bucket: %v", err)
	} else if !exists {
		err = mc.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("CRITICAL: Could not create bucket %s: %v", bucket, err)
		}
		log.Printf("Successfully created bucket: %s", bucket)
	}

	policy := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]}]}`, bucket)
	err = mc.SetBucketPolicy(ctx, bucket, policy)
	if err != nil {
		log.Printf("Warning: Could not set public policy: %v", err)
	}

	blogService := service.NewBlogService(blogStore, authClient, uploader)
	blogHandler := handler.NewBlogHandler(blogService)

	r.GET("/health", func(c *gin.Context) { c.String(200, "ok") })

	r.GET("/blogs", blogHandler.GetAll)
	r.GET("/blogs/:id", blogHandler.GetByID)
	r.POST("/blogs/images", middleware.Auth(secret), blogHandler.UploadImages)
	r.POST("/blogs", middleware.Auth(secret), blogHandler.Create)

	port := config.MustEnvDefault("PORT", "8080")
	log.Fatal(r.Run(":" + port))
}
