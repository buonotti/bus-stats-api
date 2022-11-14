package v1

import (
	"net/http"

	"github.com/buonotti/bus-stats-api/jwt"
	"github.com/buonotti/bus-stats-api/logging"
)

type RefreshRequest struct {
	Token string `json:"token" binding:"required,jwt"`
	Id    string `json:"id" binding:"required"`
}

type RefreshResponse struct {
	Token string `json:"token"`
}

func RefreshUserToken(data RefreshRequest) (RefreshResponse, int, error) {
	token, err := jwt.Service().RefreshToken(data.Token, data.Id)
	if err != nil {
		logging.ApiLogger.Error(err)
		return RefreshResponse{}, http.StatusUnauthorized, err
	}

	return RefreshResponse{
		Token: token,
	}, http.StatusOK, nil
}
