package middleware

import (
	"net/http"

	"github.com/buonotti/bus-stats-api/errors"
	"github.com/buonotti/bus-stats-api/jwt"
	"github.com/buonotti/bus-stats-api/logging"
	"github.com/buonotti/bus-stats-api/services"
	"github.com/gin-gonic/gin"
	goJWT "github.com/golang-jwt/jwt/v4"
)

const BEARER_SCHEMA = "Bearer "

var IdentityKey = "identity"

// Auth is a middleware that checks if the request has a valid JWT token and then authorizes the request
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token string from the header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) <= len(BEARER_SCHEMA)+1 {
			err := errors.TokenError.New("missing auth token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, services.ErrorResponse{Message: err.Error()})
			return
		}

		// Cut off the "Bearer " part of the header
		tokenString := authHeader[len(BEARER_SCHEMA):]

		// Parse and validate the token
		token, err := jwt.Service().ValidateToken(tokenString)
		if token.Valid && err == nil {
			claims := token.Claims.(goJWT.MapClaims)
			uid := claims["uid"].(string)
			logging.ApiLogger.Infof("authorizing user %s", claims["uid"].(string))
			if !authorizeRoute(c, uid) {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func authorizeRoute(c *gin.Context, uid string) bool {
	routeParam := c.Param("id")

	// If no id is provided in the route, then the user is authorized because the route is not user-specific
	if routeParam == "" {
		return true
	}

	// If the id in the route is the same as the id in the token, then the user is authorized otherwise he tries to access another user's data
	if routeParam != uid {
		return false
	}
	return true
}
