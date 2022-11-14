package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func ExtractUidFromToken(token *jwt.Token) string {
	return token.Claims.(jwt.MapClaims)["uid"].(string)
}

func ExtractUidFromHeader(c *gin.Context) string {
	const BEARER_SCHEMA = "Bearer "
	authHeader := c.GetHeader("Authorization")
	tokenString := authHeader[len(BEARER_SCHEMA):]
	token, err := Service().ValidateToken(tokenString)
	if err != nil {
		return ""
	}
	claims := token.Claims.(jwt.MapClaims)
	return claims["uid"].(string)
}
