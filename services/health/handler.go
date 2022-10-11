package health

import (
	"encoding/json"
	"fmt"

	"github.com/buonotti/bus-stats-api/services"
)

func Result() DbInfoResponse {
	response, err := services.RestClient.R().SetBody("INFO FOR DB;").Post(services.ApiUrl("sql"))
	if err != nil {
		fmt.Println(err)
	}

	var dbResponse DbInfoResponse
	jsonData := services.FormatResponseString(response)
	err = json.Unmarshal([]byte(jsonData), &dbResponse)
	if err != nil {
		fmt.Println(err)
	}

	return dbResponse
}
