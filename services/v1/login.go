package v1

import (
	"encoding/json"
	"net/http"

	"github.com/buonotti/bus-stats-api/models"
	"github.com/buonotti/bus-stats-api/services"
	"github.com/buonotti/bus-stats-api/util"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"` // TODO sha265
}

type LoginResponse struct {
	Token string `json:"token"`
}

func LoginUser(data LoginRequest) (models.BaseUser, int, error) {
	selectResponse, err := util.RestClient.R().
		SetBody(util.Query("SELECT * FROM user WHERE email = ?", data.Email)).
		Post(util.DatabaseUrl())
	if err != nil {
		return models.BaseUser{}, http.StatusBadRequest, services.FormatError
	}

	var selectUserResponse models.UserSelectResult
	responseString := util.FormatResponseString(selectResponse)
	err = json.Unmarshal([]byte(responseString), &selectUserResponse)
	if err != nil {
		return models.BaseUser{}, http.StatusBadRequest, services.FormatError
	}

	if len(selectUserResponse.Result) <= 0 {
		return models.BaseUser{}, http.StatusUnauthorized, services.CredentialError
	}

	if data.Password == selectUserResponse.Result[0].Password {
		return models.BaseUser{
			Id:    util.SplitDatabaseId(selectUserResponse.Result[0].Id),
			Email: data.Email,
		}, http.StatusOK, nil
	}

	return models.BaseUser{}, http.StatusUnauthorized, services.CredentialError
}
