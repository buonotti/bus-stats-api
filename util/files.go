package util

import (
	"fmt"

	"github.com/spf13/viper"
)

func BuildFileName(name string, ext string) string {
	basePath := viper.GetString("storage")
	return fmt.Sprintf("%s/%s.%s", basePath, name, ext)
}
