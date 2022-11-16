package surreal

import (
	"fmt"
	"strings"

	"github.com/buonotti/bus-stats-api/models"
	"github.com/go-resty/resty/v2"
)

func Query(query string, args ...any) (*resty.Response, error) {
	return restClient.R().SetBody(parseQuery(query, args...)).Post(Url())
}

const invalidChars = `~!#$%^&*()+{}|:"<>?[]\;',/`

type typePair struct {
	V any
	T string
}

func parseQuery(query string, args ...any) string {
	sanitizedArgs, err := sanitizeArgs(args...)
	if err != nil {
		return ""
	}
	for i := 0; strings.Contains(query, "?"); i++ {
		arg := sanitizedArgs[i]
		switch sanitizedArgs[i].T {
		case "uid":
			query = strings.Replace(query, "?", string(arg.V.(models.UserId)), 1)
		case "string":
			query = strings.Replace(query, "?", string(arg.V.(string)), 1)
		case "number":
			query = strings.Replace(query, "?", fmt.Sprintf("%f", arg.V.(float64)), 1)
		case "stringer":
			query = strings.Replace(query, "?", arg.V.(fmt.Stringer).String(), 1)
		default:
			return ""

		}
	}

	return query
}

func sanitizeArgs(args ...any) ([]typePair, error) {
	var sanitizedArgs []typePair
	for _, arg := range args {
		if _, isUid := arg.(models.UserId); isUid {
			cArg := arg.(models.UserId)
			if strings.ContainsAny(string(cArg), invalidChars) {
				return nil, fmt.Errorf("invalid character contained in userId")
			}
			sanitizedArgs = append(sanitizedArgs, typePair{V: cArg, T: "uid"})
		} else if _, isString := arg.(string); isString {
			cArg := arg.(string)
			if strings.ContainsAny(cArg, invalidChars) {
				return nil, fmt.Errorf("invalid character contained in string")
			}
			sanitizedArgs = append(sanitizedArgs, typePair{V: cArg, T: "string"})
		} else if _, isInt := arg.(int); isInt {
			cArg := arg.(int)
			sanitizedArgs = append(sanitizedArgs, typePair{V: cArg, T: "number"})
		} else if _, isFloat := arg.(float64); isFloat {
			cArg := arg.(float64)
			sanitizedArgs = append(sanitizedArgs, typePair{V: cArg, T: "number"})
		} else if _, isStringer := arg.(fmt.Stringer); isStringer {
			cArg := arg.(fmt.Stringer)
			if strings.ContainsAny(cArg.String(), invalidChars) {
				return nil, fmt.Errorf("invalid character contained in string")
			}
			sanitizedArgs = append(sanitizedArgs, typePair{V: cArg, T: "stringer"})
		} else {
			return nil, fmt.Errorf("invalid argument type")
		}
	}
	return sanitizedArgs, nil
}
