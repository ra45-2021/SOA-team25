package main

//mysql -u root -p -> root_pass -> USE blog_db -> SHOW TABLES; -> SELECT * FROM blogs; (ili koji god prompt)

import (
	"log"

	"github.com/gin-gonic/gin"

	"blog-service/internal/config"
	"blog-service/internal/db"
	"blog-service/internal/handler"
	"blog-service/internal/middleware"
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
	blogService := service.NewBlogService(blogStore)
	blogHandler := handler.NewBlogHandler(blogService)

	r.GET("/health", func(c *gin.Context) { c.String(200, "ok") })

	r.GET("/blogs", blogHandler.GetAll)
	r.GET("/blogs/:id", blogHandler.GetByID)
	r.POST("/blogs", middleware.Auth(secret), blogHandler.Create)

	port := config.MustEnvDefault("PORT", "8080")
	log.Fatal(r.Run(":" + port))
}
