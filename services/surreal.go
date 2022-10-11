package services

import "github.com/go-resty/resty/v2"

func FormatResponseString(response *resty.Response) string {
	str := response.String()
	slice := str[1 : len(str)-1]
	return slice
}
