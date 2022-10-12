package util

import "github.com/go-resty/resty/v2"

var RestClient = resty.New().
	SetBasicAuth("root", "root").
	SetHeader("Content-Type", "application/json").
	SetHeader("NS", "buonotti").
	SetHeader("DB", "bus-stats").
	SetDisableWarn(true)
