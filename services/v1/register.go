package v1

import (
	"encoding/json"
	"fmt"
	"github.com/buonotti/bus-stats-api/models"
	"github.com/buonotti/bus-stats-api/util"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	Result string `json:"result"`
	Token  string `json:"token"`
	Id string `json:"id"`
}

func RegisterUser(data RegisterRequest) (RegisterResponse, error) {
	checkForUsernameQuery := util.FormatQuery("SELECT * FROM user WHERE email = ?;", data.Email)
	selectResponse, err := util.RestClient.R().SetBody(checkForUsernameQuery).Post(util.DatabaseUrl("sql"))
	if err != nil {
		panic(err)
	}

	var selectUserResponse models.UserSelectResult
	responseString := util.FormatResponseString(selectResponse)
	err = json.Unmarshal([]byte(responseString), &selectUserResponse)
	if err != nil {
		panic(err)
	}

	if len(selectUserResponse.Result) >= 1 {
		return RegisterResponse{}, fmt.Errorf("username already present")
	}

	insertUserQuery := util.FormatQuery("CREATE user SET email = ?, password = ? RETURN id;", data.Email, data.Password)
	insertResponse, err := util.RestClient.R().SetBody(insertUserQuery).Post(util.DatabaseUrl("sql"))
	if err != nil {
		panic(err)
	}

	var insertUserResponse models.UserInsertResult
	responseString = util.FormatResponseString(insertResponse)
	err = json.Unmarshal([]byte(responseString), &insertUserResponse)
	if err != nil {
		panic(err)
	}

	userId := util.SplitDatabaseId(insertUserResponse.Result[0].Id)

	token, err := util.GenerateToken(userId)
	if err != nil {
		panic(err)
	}
	return RegisterResponse{
		Result: "OK",
		Token: token,
		Id: userId,
	}, nil
}
