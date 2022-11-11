package util

import (
	"fmt"

	"github.com/spf13/viper"
)

func FileName(name string, ext string) string {
	basePath := viper.GetString("storage.content_root")
	return fmt.Sprintf("%s/%s.%s", basePath, name, ext)
}
