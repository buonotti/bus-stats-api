package middleware

import (
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/buonotti/bus-stats-api/models"
	serviceV1 "github.com/buonotti/bus-stats-api/services/v1"
	"github.com/buonotti/bus-stats-api/util"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var IdentityKey = "identity"

func Auth() *jwt.GinJWTMiddleware {
	middleware, _ := jwt.New(&jwt.GinJWTMiddleware{
		Realm:      "bus-stats",
		Key:        []byte(viper.GetString("api.key")),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.BaseUser); ok {
				return jwt.MapClaims{IdentityKey: v.Id}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.BaseUser{
				Id: claims[IdentityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			util.ApiLogger.Debugf("logging in %s", c.RemoteIP())
			var loginRequest serviceV1.LoginRequest
			if err := c.ShouldBind(&loginRequest); err != nil {
				return "", err
			}
			user, err, _ := serviceV1.LoginUser(loginRequest)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			return &user, nil
		},
		LoginResponse: func(c *gin.Context, code int, message string, time time.Time) {
			token := message
			c.JSON(code, serviceV1.LoginResponse{
				Token: token,
			})
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			util.ApiLogger.Debugf("authorizing %s", c.RemoteIP())
			if v, ok := data.(*models.BaseUser); ok {
				if v != nil {
					urlId := c.Param("id")
					if urlId != "" {
						if urlId != v.Id {
							util.ApiLogger.Warnf("user with id %s tried to access resources of %s", v.Id, urlId)
							return false
						} else {
							return true
						}
					}
					
				}
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"message": message,
			})
		},
		TokenLookup:    "header: Authorization",
		TokenHeadName:  "Bearer",
		TimeFunc:       time.Now,
		SendCookie:     true,
		SecureCookie:   false, // TODO
		CookieHTTPOnly: true,
		CookieDomain:   "localhost:8080", // TODO
		CookieName:     "token",
		CookieSameSite: http.SameSiteDefaultMode,
		CookieMaxAge:   1 * time.Hour,
	})
	return middleware
}
