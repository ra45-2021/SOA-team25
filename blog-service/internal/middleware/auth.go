package middleware

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(secret []byte) gin.HandlerFunc {
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

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token payload"})
			return
		}

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
