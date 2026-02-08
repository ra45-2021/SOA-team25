package middleware

import (
	"strconv"
	"strings"

	"auth-service/internal/model"

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
		tok, err := jwt.Parse(raw, func(token *jwt.Token) (interface{}, error) { return secret, nil })
		if err != nil || !tok.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := tok.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

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
		if id <= 0 {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token payload"})
			return
		}

		username, _ := claims["username"].(string)
		role, _ := claims[model.RoleClaimKey].(string)

		c.Set("id", id)
		c.Set("username", username)
		c.Set("role", role)
		c.Next()
	}
}
