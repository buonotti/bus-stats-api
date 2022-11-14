package v1

import (
	"encoding/json"
	"net/http"

	"github.com/buonotti/bus-stats-api/jwt"
	"github.com/buonotti/bus-stats-api/logging"
	"github.com/buonotti/bus-stats-api/models"
	"github.com/buonotti/bus-stats-api/services"
	"github.com/buonotti/bus-stats-api/surreal"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"` // TODO sha265
}

type LoginResponse struct {
	Token string `json:"token"`
	Uid   string `json:"uid"`
}

func LoginUser(data LoginRequest) (LoginResponse, int, error) {
	selectResponse, err := surreal.Query("SELECT * FROM user WHERE email = ?", data.Email)
	if err != nil {
		logging.DbLogger.Error(err)
		return LoginResponse{}, http.StatusBadRequest, services.FormatError
	}

	var selectUserResponse models.UserSelectResult
	responseString := surreal.FormatResponse(selectResponse)
	err = json.Unmarshal([]byte(responseString), &selectUserResponse)
	if err != nil {
		logging.ApiLogger.Error(err)
		return LoginResponse{}, http.StatusBadRequest, services.FormatError
	}

	if len(selectUserResponse.Result) <= 0 {
		return LoginResponse{}, http.StatusUnauthorized, services.CredentialError
	}

	if data.Password != selectUserResponse.Result[0].Password {
		return LoginResponse{}, http.StatusUnauthorized, services.CredentialError
	}

	userId := surreal.SplitDatabaseId(selectUserResponse.Result[0].Id)
	token, err := jwt.Service().GenerateToken(userId)
	if err != nil {
		logging.ApiLogger.Error(err)
		return LoginResponse{}, http.StatusUnauthorized, services.CredentialError
	}

	return LoginResponse{
		Uid:   userId,
		Token: token,
	}, http.StatusOK, nil
}
