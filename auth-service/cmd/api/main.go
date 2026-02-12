package main

import (
	"auth-service/internal/config"
	"auth-service/internal/db"
	"auth-service/internal/handler"
	"auth-service/internal/middleware"
	"auth-service/internal/service"
	"auth-service/internal/store"
	"auth-service/pb"
	"log"
	"net"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	jwtSecret := []byte(config.MustEnv("JWT_SECRET"))
	mongoURI := config.MustEnv("MONGO_URI")
	dbName := config.MustEnvDefault("MONGO_DB", "auth_db")

	client := db.MustMongo(mongoURI)
	mdb := client.Database(dbName)

	userStore := store.NewMongoUserStore(mdb)
	authSvc := service.NewAuthService(userStore, jwtSecret)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, handler.NewAuthGRPCHandler(authSvc))

	go func() {
		log.Println("gRPC server listening on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve grpc: %v", err)
		}
	}()

	r := gin.Default()
	r.Use(middleware.CORS())

	authHandler := handler.NewAuthHandler(authSvc)

	r.GET("/health", func(c *gin.Context) { c.String(200, "ok") })
	r.POST("/users", authHandler.Register)
	r.POST("/users/login", authHandler.Login)
	r.GET("/users/me", middleware.Auth(jwtSecret), authHandler.Me)
	r.GET("/users/:id", authHandler.GetUserByID)

	port := config.MustEnvDefault("PORT", "8080")
	log.Fatal(r.Run(":" + port))
}