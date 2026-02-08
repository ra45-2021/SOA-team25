package main
//mongosh -> use auth_db -> show collections -> db.users.find().pretty()

import (
	"log"

	"github.com/gin-gonic/gin"

	"auth-service/internal/config"
	"auth-service/internal/db"
	"auth-service/internal/handler"
	"auth-service/internal/middleware"
	"auth-service/internal/service"
	"auth-service/internal/store"
)

func main() {
	r := gin.Default()
	r.Use(middleware.CORS())

	jwtSecret := []byte(config.MustEnv("JWT_SECRET"))

	client := db.MustMongo(config.MustEnv("MONGO_URI"))
	dbName := config.MustEnvDefault("MONGO_DB", "auth_db")
	mdb := client.Database(dbName)

	userStore := store.NewMongoUserStore(mdb)
	authSvc := service.NewAuthService(userStore, jwtSecret)
	authHandler := handler.NewAuthHandler(authSvc)

	r.GET("/health", func(c *gin.Context) { c.String(200, "ok") })

	r.POST("/users", authHandler.Register)
	r.POST("/users/login", authHandler.Login)

	r.GET("/users/me", middleware.Auth(jwtSecret), authHandler.Me)

	port := config.MustEnvDefault("PORT", "8080")
	log.Fatal(r.Run(":" + port))
}
