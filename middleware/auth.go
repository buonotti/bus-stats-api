package middleware

import (
	"net/http"

	"github.com/buonotti/bus-stats-api/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var IdentityKey = "identity"

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) <= len(BEARER_SCHEMA)+1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing auth token"})
			return
		}
		tokenString := authHeader[len(BEARER_SCHEMA):]
		token, err := util.JWTAuthService().ValidateToken(tokenString)
		if token.Valid && err == nil {
			claims := token.Claims.(jwt.MapClaims)
			util.ApiLogger.Infof("Authenticated user %s", claims["uid"].(string))
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
