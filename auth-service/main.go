package main
//mongosh -> use auth_db -> show collections -> db.users.find().pretty()

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const roleClaimKey = "http://schemas.microsoft.com/ws/2008/06/identity/claims/role"

type UserDoc struct {
	MongoID      primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	ID           int64              `bson:"id" json:"id"`
	Username     string             `bson:"username" json:"username"`
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"passwordHash" json:"-"`
	Role         string             `bson:"role" json:"role"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
}

type RegisterReq struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required"`
}


type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	ID          int64  `json:"id"`
	AccessToken string `json:"accessToken"`
}

type CounterDoc struct {
	ID    string `bson:"_id"`
	Value int64  `bson:"value"`
}

func main() {
	r := gin.Default()
	r.Use(cors())

	jwtSecret := []byte(mustEnv("JWT_SECRET"))

	client := mustMongo(mustEnv("MONGO_URI"))
	dbName := mustEnvDefault("MONGO_DB", "auth_db")
	db := client.Database(dbName)

	users := db.Collection("users")
	counters := db.Collection("counters")

	ensureIndexes(users)

	r.GET("/health", func(c *gin.Context) { c.String(200, "ok") })


	r.POST("/users", func(c *gin.Context) {
		role := strings.ToUpper(strings.TrimSpace(c.Query("role")))
		if role == "" {
			role = "TOURIST"
		}
		if role != "GUIDE" && role != "TOURIST" {
			c.JSON(400, gin.H{"error": "role must be GUIDE or TOURIST"})
			return
		}

		var req RegisterReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		role = strings.ToUpper(strings.TrimSpace(req.Role))

		if role != "GUIDE" && role != "TOURIST" {
			c.JSON(400, gin.H{"error": "role must be GUIDE or TOURIST"})
			return
		}


		pwHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(500, gin.H{"error": "hash failed"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		newID, err := nextID(ctx, counters)
		if err != nil {
			c.JSON(500, gin.H{"error": "id generation failed"})
			return
		}

		doc := UserDoc{
			ID:           newID,
			Email:        strings.TrimSpace(req.Email),
			Username:     strings.TrimSpace(req.Username),
			PasswordHash: string(pwHash),
			Role:         role,
			CreatedAt:    time.Now(),
		}

		_, err = users.InsertOne(ctx, doc)
		if err != nil {
			c.JSON(400, gin.H{"error": "username or email already exists"})
			return
		}

		token, err := issueToken(jwtSecret, doc)
		if err != nil {
			c.JSON(500, gin.H{"error": "token failed"})
			return
		}

		c.JSON(201, AuthResponse{ID: doc.ID, AccessToken: token})
	})

	r.POST("/users/login", func(c *gin.Context) {
		var req LoginReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		var u UserDoc
		if err := users.FindOne(ctx, bson.M{"username": req.Username}).Decode(&u); err != nil {
			c.JSON(401, gin.H{"error": "invalid credentials"})
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)) != nil {
			c.JSON(401, gin.H{"error": "invalid credentials"})
			return
		}

		token, err := issueToken(jwtSecret, u)
		if err != nil {
			c.JSON(500, gin.H{"error": "token failed"})
			return
		}

		c.JSON(200, AuthResponse{ID: u.ID, AccessToken: token})
	})

	r.GET("/users/me", authMiddleware(jwtSecret), func(c *gin.Context) {
		c.JSON(200, gin.H{
			"id":       c.MustGet("id").(int64),
			"username": c.MustGet("username").(string),
			"role":     c.MustGet("role").(string),
		})
	})

	port := mustEnvDefault("PORT", "8080")
	log.Fatal(r.Run(":" + port))
}

func issueToken(secret []byte, u UserDoc) (string, error) {
	claims := jwt.MapClaims{
		"id":       u.ID,    
		"username": u.Username, 
		roleClaimKey: u.Role,   
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(secret)
}

func authMiddleware(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing token"})
			return
		}

		raw := strings.TrimPrefix(h, "Bearer ")
		tok, err := jwt.Parse(raw, func(token *jwt.Token) (interface{}, error) { return secret, nil })
		if err != nil || !tok.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		claims := tok.Claims.(jwt.MapClaims)

		var id int64
		switch v := claims["id"].(type) {
		case float64:
			id = int64(v)
		case int64:
			id = v
		case string:
			parsed, _ := strconv.ParseInt(v, 10, 64)
			id = parsed
		default:
			id = 0
		}

		username, _ := claims["username"].(string)
		role, _ := claims[roleClaimKey].(string)

		c.Set("id", id)
		c.Set("username", username)
		c.Set("role", role)
		c.Next()
	}
}

func nextID(ctx context.Context, counters *mongo.Collection) (int64, error) {
	res := counters.FindOneAndUpdate(
		ctx,
		bson.M{"_id": "users"},
		bson.M{"$inc": bson.M{"value": 1}},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	)

	var c CounterDoc
	if err := res.Decode(&c); err != nil {
		return 0, err
	}
	return c.Value, nil
}

func ensureIndexes(users *mongo.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, _ = users.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "username", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "id", Value: 1}}, Options: options.Index().SetUnique(true)},
	})
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

func mustMongo(uri string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		panic(err)
	}
	return client
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
