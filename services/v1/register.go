package v1

import (
	"encoding/json"
	"net/http"

	"github.com/buonotti/bus-stats-api/errors"
	"github.com/buonotti/bus-stats-api/logging"
	"github.com/buonotti/bus-stats-api/models"
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
		logging.Logger.Errorf("%s\n", err.Error())
		return RegisterResponse{}, http.StatusBadRequest, err
	}

	var selectUserResponse models.UserSelectResult
	responseString := surreal.FormatResponse(selectResponse)
	err = json.Unmarshal([]byte(responseString), &selectUserResponse)
	if err != nil {
		err = errors.SurrealDeserializaError.Wrap(err, "cannot deserialize surreal response")
		logging.Logger.Errorf("%s\n", err.Error())
		return RegisterResponse{}, http.StatusBadRequest, err
	}

	if len(selectUserResponse.Result) >= 1 {
		return RegisterResponse{}, http.StatusBadRequest, errors.UserAlreadyExistsError.New("user already exists")
	}

	insertResponse, err := surreal.Query("CREATE user SET email = ?, password = ?, image = {name: '', type: ''};", data.Email, data.Password)
	if err != nil {
		logging.Logger.Errorf("%s\n", err.Error())
		return RegisterResponse{}, http.StatusBadRequest, err
	}

	var insertUserResponse models.UserInsertResult
	responseString = surreal.FormatResponse(insertResponse)
	err = json.Unmarshal([]byte(responseString), &insertUserResponse)
	if err != nil {
		err = errors.SurrealDeserializaError.Wrap(err, "cannot deserialize surreal response")
		logging.Logger.Errorf("%s\n", err.Error())
		return RegisterResponse{}, http.StatusBadRequest, err
	}

	return RegisterResponse{
		Result: "ok",
	}, http.StatusOK, nil
}
