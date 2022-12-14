package util

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func FileName(name string) string {
	basePath := viper.GetString("storage.content_root")
	return fmt.Sprintf("%s%c%s", basePath, os.PathSeparator, name)
}
