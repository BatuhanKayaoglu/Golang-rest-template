package middleware

import (
	"golang-rest-api-template/pkg/auth"
	"golang-rest-api-template/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BearerSchema = "Bearer "
		header := c.GetHeader("Authorization")
		if header == "" {
			response.Unauthorized(c, "Missing Authorization Header")
			return
		}

		if !strings.HasPrefix(header, BearerSchema) {
			response.Unauthorized(c, "Invalid Authorization Header")
			return
		}

		tokenStr := header[len(BearerSchema):]
		claims := &auth.CustomClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return auth.JwtKey, nil
		})

		if err != nil {
			response.Unauthorized(c, "Invalid token")
			return
		}

		if !token.Valid {
			response.Unauthorized(c, "Invalid token")
			return
		}

		c.Set("username", claims.Username)
		c.Set("userId", claims.UserId)
		c.Next()
	}
}
