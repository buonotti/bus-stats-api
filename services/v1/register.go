package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/buonotti/bus-stats-api/models"
	"github.com/buonotti/bus-stats-api/services"
	"github.com/buonotti/bus-stats-api/util"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"` // TODO add sha256
}

type RegisterResponse struct {
	Result string `json:"result"`
	Id     string `json:"id"`
}

func RegisterUser(data RegisterRequest) (RegisterResponse, int, error) {
	selectResponse, err := util.RestClient.R().
		SetBody(util.Query("SELECT * FROM user WHERE email = ?;", data.Email)).
		Post(util.DatabaseUrl())
	if err != nil {
		util.ApiLogger.Error(err)
		return RegisterResponse{}, http.StatusBadRequest, services.CredentialError
	}

	var selectUserResponse models.UserSelectResult
	responseString := util.FormatResponseString(selectResponse)
	err = json.Unmarshal([]byte(responseString), &selectUserResponse)
	if err != nil {
		util.ApiLogger.Debug(selectResponse.Request.Header)
		util.ApiLogger.Error(err)
		return RegisterResponse{}, http.StatusBadRequest, services.FormatError
	}

	if len(selectUserResponse.Result) >= 1 {
		return RegisterResponse{}, http.StatusBadRequest, fmt.Errorf("email already in use")
	}

	insertResponse, err := util.RestClient.R().
		SetBody(util.Query("CREATE user SET email = ?, password = ?, image = {name: '', type: ''} RETURN id;", data.Email, data.Password)).
		Post(util.DatabaseUrl())
	if err != nil {
		util.ApiLogger.Error(err)
		return RegisterResponse{}, http.StatusBadRequest, services.FormatError
	}

	var insertUserResponse models.UserInsertResult
	responseString = util.FormatResponseString(insertResponse)
	err = json.Unmarshal([]byte(responseString), &insertUserResponse)
	if err != nil {
		util.ApiLogger.Error(err)
		return RegisterResponse{}, http.StatusBadRequest, services.FormatError
	}

	userId := util.SplitDatabaseId(insertUserResponse.Result[0].Id)

	return RegisterResponse{
		Result: "OK",
		Id:     userId,
	}, http.StatusOK, nil
}
