package jwt

import (
	"fmt"
	"time"

	"github.com/buonotti/bus-stats-api/config"
	"github.com/golang-jwt/jwt/v4"
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

func Service() JwtService {
	return &jwtServices{
		secretKey: config.Get("api.key"),
		issure:    "Buonotti",
	}
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

func (service *jwtServices) RefreshToken(token string, userId string) (string, error) {
	if _, err := service.ValidateToken(token); err != nil {
		return "", err
	}

	return service.GenerateToken(userId)
}
