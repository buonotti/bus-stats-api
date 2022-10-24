package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/buonotti/bus-stats-api/models"
	"github.com/buonotti/bus-stats-api/util"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"` // TODO sha265
}

type LoginResponse struct {
	Token string `json:"token"`
}

func LoginUser(data LoginRequest) (models.BaseUser, error, int) {
	selectResponse, err := util.RestClient.R().
		SetBody(util.Query("SELECT * FROM user WHERE email = ?", data.Email)).
		Post(util.DatabaseUrl())
	if err != nil {
		return models.BaseUser{}, fmt.Errorf("could not fetch user data"), http.StatusInternalServerError
	}

	var selectUserResponse models.UserSelectResult
	responseString := util.FormatResponseString(selectResponse)
	err = json.Unmarshal([]byte(responseString), &selectUserResponse)
	if err != nil {
		return models.BaseUser{}, fmt.Errorf("unexpected database response for userSelectResult"), http.StatusBadRequest
	}

	if len(selectUserResponse.Result) <= 0 {
		return models.BaseUser{}, fmt.Errorf("user does not exist"), http.StatusBadRequest
	}

	if data.Password == selectUserResponse.Result[0].Password {
		return models.BaseUser{
				Id:    util.SplitDatabaseId(selectUserResponse.Result[0].Id),
				Email: data.Email,
			},
			nil, http.StatusOK
	}

	return models.BaseUser{}, fmt.Errorf("credentials are incorrect"), http.StatusUnauthorized
}
