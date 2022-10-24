package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/buonotti/bus-stats-api/models"
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

func RegisterUser(data RegisterRequest) (RegisterResponse, error, int) {
	selectResponse, err := util.RestClient.R().
		SetBody(util.Query("SELECT * FROM user WHERE email = ?;", data.Email)).
		Post(util.DatabaseUrl())
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("could not verify username uniqueness"), http.StatusInternalServerError
	}

	var selectUserResponse models.UserSelectResult
	responseString := util.FormatResponseString(selectResponse)
	err = json.Unmarshal([]byte(responseString), &selectUserResponse)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("unexpected database response for userSelectResult"), http.StatusBadRequest
	}

	if len(selectUserResponse.Result) >= 1 {
		return RegisterResponse{}, fmt.Errorf("username already present"), http.StatusBadRequest
	}

	insertResponse, err := util.RestClient.R().
		SetBody(util.Query("CREATE user SET email = ?, password = ?, image = {name: '', type: ''} RETURN id;", data.Email, data.Password)).
		Post(util.DatabaseUrl())
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("could not create user"), http.StatusInternalServerError
	}

	var insertUserResponse models.UserInsertResult
	responseString = util.FormatResponseString(insertResponse)
	err = json.Unmarshal([]byte(responseString), &insertUserResponse)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("unexpected database response for userInsertResult"), http.StatusInternalServerError
	}

	userId := util.SplitDatabaseId(insertUserResponse.Result[0].Id)

	if err != nil {
		return RegisterResponse{}, fmt.Errorf("could not generate token"), http.StatusInternalServerError
	}

	return RegisterResponse{
		Result: "OK",
		Id:     userId,
	}, nil, http.StatusOK
}
