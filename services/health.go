package services


import (
	"encoding/json"
	"fmt"
	"github.com/buonotti/bus-stats-api/util"
)

type DbInfoResponse struct {
	Time   string       `json:"time"`
	Status string       `json:"status"`
	Result DbInfoResult `json:"result"`
}

type DbInfoResult struct {
	Dl any `json:"dl"`
	Dt any `json:"dt"`
	Sc any `json:"sc"`
	Tb any `json:"tb"`
}

func DbInfo() DbInfoResponse {
	response, err := util.RestClient.R().SetBody("INFO FOR DB;").Post(util.DatabaseUrl("sql"))
	if err != nil {
		fmt.Println(err)
	}

	var dbResponse DbInfoResponse
	jsonData := util.FormatResponseString(response)
	err = json.Unmarshal([]byte(jsonData), &dbResponse)
	if err != nil {
		fmt.Println(err)
	}

	return dbResponse
}