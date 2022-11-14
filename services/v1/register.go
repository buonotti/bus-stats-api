package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/buonotti/bus-stats-api/logging"
	"github.com/buonotti/bus-stats-api/models"
	"github.com/buonotti/bus-stats-api/services"
	"github.com/buonotti/bus-stats-api/surreal"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"` // TODO add sha256
}

type RegisterResponse struct {
	Result string `json:"result"`
}

func RegisterUser(data RegisterRequest) (RegisterResponse, int, error) {
	selectResponse, err := surreal.Query("SELECT * FROM user WHERE email = ?;", data.Email)
	if err != nil {
		logging.DbLogger.Error(err)
		return RegisterResponse{}, http.StatusBadRequest, services.CredentialError
	}

	var selectUserResponse models.UserSelectResult
	responseString := surreal.FormatResponse(selectResponse)
	err = json.Unmarshal([]byte(responseString), &selectUserResponse)
	if err != nil {
		logging.ApiLogger.Error(err)
		return RegisterResponse{}, http.StatusBadRequest, services.FormatError
	}

	if len(selectUserResponse.Result) >= 1 {
		return RegisterResponse{}, http.StatusBadRequest, fmt.Errorf("email already in use")
	}

	insertResponse, err := surreal.Query("CREATE user SET email = ?, password = ?, image = {name: '', type: ''};", data.Email, data.Password)
	if err != nil {
		logging.DbLogger.Error(err)
		return RegisterResponse{}, http.StatusBadRequest, services.FormatError
	}

	var insertUserResponse models.UserInsertResult
	responseString = surreal.FormatResponse(insertResponse)
	err = json.Unmarshal([]byte(responseString), &insertUserResponse)
	if err != nil {
		logging.ApiLogger.Error(err)
		return RegisterResponse{}, http.StatusBadRequest, services.FormatError
	}

	return RegisterResponse{
		Result: "ok",
	}, http.StatusOK, nil
}
