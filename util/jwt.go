package util

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type JwtService interface {
	GenerateToken(userId string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	RefreshToken(token string, userId string) (string, error)
}

type authCustomClaims struct {
	Uid string `json:"uid"`
	jwt.StandardClaims
}

type jwtServices struct {
	secretKey string
	issure    string
}

func JWTAuthService() JwtService {
	return &jwtServices{
		secretKey: getSecretKey(),
		issure:    "Bikash",
	}
}

func getSecretKey() string {
	return viper.GetString("api.key")
}

func (service *jwtServices) GenerateToken(uid string) (string, error) {
	claims := &authCustomClaims{
		uid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			Issuer:    service.issure,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(service.secretKey))
}

func (service *jwtServices) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("invalid token %s", token.Header["alg"])
		}
		return []byte(service.secretKey), nil
	})
}

func ExtractUidFromToken(token *jwt.Token) string {
	return token.Claims.(jwt.MapClaims)["uid"].(string)
}

func ExtractUidFromHeader(c *gin.Context) string {
	const BEARER_SCHEMA = "Bearer "
	authHeader := c.GetHeader("Authorization")
	tokenString := authHeader[len(BEARER_SCHEMA):]
	token, err := JWTAuthService().ValidateToken(tokenString)
	if err != nil {
		return ""
	}
	claims := token.Claims.(jwt.MapClaims)
	return claims["uid"].(string)
}

func (service *jwtServices) RefreshToken(token string, userId string) (string, error) {
	if _, err := service.ValidateToken(token); err != nil {
		return "", err
	}

	return service.GenerateToken(userId)
}
