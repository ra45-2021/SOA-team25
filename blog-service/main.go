package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/go-sql-driver/mysql"
)

type CreateBlogReq struct {
	Title               string   `json:"title" binding:"required"`
	DescriptionMarkdown string   `json:"descriptionMarkdown" binding:"required"`
	Images              []string `json:"images"`
}

type BlogDto struct {
	ID                 int64     `json:"id"`
	Title              string    `json:"title"`
	DescriptionMarkdown string   `json:"descriptionMarkdown"`
	CreatedAt          time.Time `json:"createdAt"`
	Images             []string  `json:"images,omitempty"`
	AuthorUserID       string    `json:"authorUserId"`
}

func main() {
	r := gin.Default()
	r.Use(cors())

	secret := []byte(mustEnv("JWT_SECRET"))
	dsn := mustEnv("MYSQL_DSN")

	db := mustMySQL(dsn)
	mustInitSchema(db)

	r.GET("/health", func(c *gin.Context) { c.String(200, "ok") })

	r.GET("/blogs", func(c *gin.Context) {
		rows, err := db.Query(`
			SELECT id, title, description_markdown, created_at, images_json, author_user_id
			FROM blogs
			ORDER BY created_at DESC`)
		if err != nil {
			c.JSON(500, gin.H{"error": "db error"})
			return
		}
		defer rows.Close()

		var out []BlogDto
		for rows.Next() {
			var b BlogDto
			var imagesJSON sql.NullString
			if err := rows.Scan(&b.ID, &b.Title, &b.DescriptionMarkdown, &b.CreatedAt, &imagesJSON, &b.AuthorUserID); err != nil {
				continue
			}
			if imagesJSON.Valid && imagesJSON.String != "" {
				_ = json.Unmarshal([]byte(imagesJSON.String), &b.Images)
			}
			out = append(out, b)
		}
		c.JSON(200, out)
	})

	r.GET("/blogs/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil || id <= 0 {
			c.JSON(400, gin.H{"error": "invalid id"})
			return
		}

		var b BlogDto
		var imagesJSON sql.NullString
		err = db.QueryRow(`
			SELECT id, title, description_markdown, created_at, images_json, author_user_id
			FROM blogs
			WHERE id = ?`, id).
			Scan(&b.ID, &b.Title, &b.DescriptionMarkdown, &b.CreatedAt, &imagesJSON, &b.AuthorUserID)

		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}
		if err != nil {
			c.JSON(500, gin.H{"error": "db error"})
			return
		}
		if imagesJSON.Valid && imagesJSON.String != "" {
			_ = json.Unmarshal([]byte(imagesJSON.String), &b.Images)
		}

		c.JSON(200, b)
	})

	r.POST("/blogs", auth(secret), func(c *gin.Context) {
		var req CreateBlogReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		userID := c.MustGet("userId").(string)
		now := time.Now()

		imagesBytes, _ := json.Marshal(req.Images)
		imagesStr := string(imagesBytes)
		if len(req.Images) == 0 {
			imagesStr = ""
		}

		res, err := db.Exec(`
			INSERT INTO blogs(title, description_markdown, created_at, images_json, author_user_id)
			VALUES(?,?,?,?,?)`,
			req.Title, req.DescriptionMarkdown, now, nullableString(imagesStr), userID,
		)
		if err != nil {
			c.JSON(500, gin.H{"error": "insert failed"})
			return
		}

		newID, _ := res.LastInsertId()
		out := BlogDto{
			ID:                  newID,
			Title:               req.Title,
			DescriptionMarkdown: req.DescriptionMarkdown,
			CreatedAt:           now,
			Images:              req.Images,
			AuthorUserID:        userID,
		}
		c.JSON(201, out)
	})

	port := mustEnvDefault("PORT", "8080")
	log.Fatal(r.Run(":" + port))
}

func mustMySQL(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(10)

	for i := 0; i < 20; i++ {
		if err := db.Ping(); err == nil {
			return db
		}
		time.Sleep(1 * time.Second)
	}
	panic("mysql not ready")
}

func mustInitSchema(db *sql.DB) {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS blogs (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  description_markdown TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  images_json JSON NULL,
  author_user_id VARCHAR(24) NOT NULL
);
`)
	if err != nil {
		panic(err)
	}
}

func auth(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {

		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing token"})
			return
		}

		raw := strings.TrimPrefix(h, "Bearer ")

		token, err := jwt.Parse(raw, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		idVal, ok := claims["id"]
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token payload"})
			return
		}

		var userID int64

		switch v := idVal.(type) {
		case float64:
			userID = int64(v)
		case int64:
			userID = v
		case string:
			parsed, _ := strconv.ParseInt(v, 10, 64)
			userID = parsed
		default:
			userID = 0
		}

		if userID <= 0 {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token payload"})
			return
		}

		c.Set("userId", userID)
		c.Next()
	}
}


func nullableString(s string) interface{} {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func mustEnv(key string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		panic("missing env: " + key)
	}
	return v
}
func mustEnvDefault(key, def string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	return v
}
