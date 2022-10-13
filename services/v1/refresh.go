package v1

import "github.com/buonotti/bus-stats-api/util"

type RefreshRequest struct {
	Token string `json:"token"`
	Id string `json:"id"`
}

type RefreshResponse struct {
	Result string `json:"result"`
	Token  string `json:"token"`
}

func RefreshToken(data RefreshRequest) (RefreshResponse,error) {
	err := util.TokenValidString(data.Token)
	if err != nil {
		return RefreshResponse{}, err
	}
	newToken, err := util.GenerateToken(data.Id)
	return RefreshResponse{
		Result: "OK",
		Token: newToken,
	}, nil
}
