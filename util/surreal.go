package util

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/buonotti/bus-stats-api/models"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

func FormatResponseString(response *resty.Response) string {
	str := response.String()
	slice := str[1 : len(str)-1]
	return slice
}

func Query(query string, args ...interface{}) string {
	for i := 0; strings.Contains(query, "?"); i++ {
		arg := args[i]
		if _, isUid := arg.(models.UserId); isUid {
			query = strings.Replace(query, "?", string(arg.(models.UserId)), 1)
		} else if _, isString := arg.(string); isString {
			query = strings.Replace(query, "?", fmt.Sprintf("'%s'", arg.(string)), 1)
		} else if _, isInt := arg.(int); isInt {
			query = strings.Replace(query, "?", fmt.Sprintf("%d", arg.(int)), 1)
		} else if _, isFloat := arg.(float64); isFloat {
			query = strings.Replace(query, "?", fmt.Sprintf("%f", arg.(float64)), 1)
		} else if _, isStringer := arg.(fmt.Stringer); isStringer {
			query = strings.Replace(query, "?", fmt.Sprintf("'%s'", arg.(fmt.Stringer).String()), 1)
		}
	}

	return query
}

func DatabaseUrl() string {
	protocol := viper.GetString(ConfigValue("database.{env}.protocol"))
	host := viper.GetString(ConfigValue("database.{env}.host"))
	port := viper.GetInt(ConfigValue("database.{env}.port"))

	return protocol + "://" + host + ":" + strconv.Itoa(port) + "/sql"
}

func SplitDatabaseId(id string) string {
	return strings.Split(id, ":")[1]
}
