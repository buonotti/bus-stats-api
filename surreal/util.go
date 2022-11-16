package surreal

import (
	"strings"

	"github.com/go-resty/resty/v2"
)

func SplitDatabaseId(id string) string {
	return strings.Split(id, ":")[1]
}

func ScaffoldDB() error {
	_, err := Query(`
	DEFINE TABLE user SCHEMAFULL;
	DEFINE FIELD email ON user TYPE string;
	DEFINE FIELD password ON user TYPE string;
	DEFINE FIELD image ON user TYPE object;
	DEFINE FIELD image.name ON user TYPE string;
	DEFINE FIELD image.type ON user TYPE string;
	DEFINE TABLE stop SCHEMAFULL;
	DEFINE FIELD name ON stop TYPE string;
	DEFINE FIELD location ON stop TYPE array;
	DEFINE TABLE line SCHEMAFULL;
	DEFINE FIELD name ON line TYPE string;
	`)
	return err
}

func PingDB() error {
	_, err := Query("INFO FOR DB;")
	return err
}

func FormatResponse(response *resty.Response) string {
	str := response.String()
	slice := str[1 : len(str)-1]
	return slice
}
