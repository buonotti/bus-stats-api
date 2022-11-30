package errors

import (
	"os"

	"github.com/buonotti/bus-stats-api/logging"
)

var DoDebug bool = false

func CheckError(err error) {
	if err == nil {
		return
	}
	if DoDebug {
		logging.Logger.Errorf("%+v\n", err)
	} else {
		logging.Logger.Errorf("%s\n", err.Error())
	}
	os.Exit(1)
}

func CheckApiError[TData any](data TData, status int, err error) (TData, int, error) {
	if err == nil {
		return data, status, nil
	}
	logging.Logger.Errorf("%s\n", err.Error())
	return data, status, err
}
