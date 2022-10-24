package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/buonotti/bus-stats-api/models"
	"github.com/buonotti/bus-stats-api/util"
)

type RefreshRequest struct {
	Token string `json:"token" binding:"required,jwt"`
	Id    string `json:"id" binding:"required"`
}

type RefreshResponse struct {
	Result string `json:"result"`
	Token  string `json:"token"`
}

func RefreshToken(data RefreshRequest) (RefreshResponse, error, int) {
	err := util.TokenValidString(data.Token)
	if err != nil {
		return RefreshResponse{}, fmt.Errorf("token is invalid"), http.StatusBadRequest
	}

	checkForIdQuery := util.Query("SELECT * FROM user WHERE id = ?", data.Id)
	selectResponse, err := util.RestClient.R().SetBody(checkForIdQuery).Post(util.DatabaseUrl())
	if err != nil {
		return RefreshResponse{}, fmt.Errorf("could not verify user id"), http.StatusInternalServerError
	}

	var userSelectresult models.UserSelectResult
	responseString := util.FormatResponseString(selectResponse)
	err = json.Unmarshal([]byte(responseString), &userSelectresult)
	if err != nil {
		return RefreshResponse{}, fmt.Errorf("unexpected database response for userSelectResult"), http.StatusBadRequest
	}

	if len(userSelectresult.Result) == 0 {
		return RefreshResponse{}, fmt.Errorf("user with the provided id is not registered"), http.StatusUnauthorized
	}

	newToken, err := util.GenerateToken(data.Id)
	return RefreshResponse{
		Result: "OK",
		Token:  newToken,
	}, nil, http.StatusOK
}
