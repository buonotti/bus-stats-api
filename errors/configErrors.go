package errors

import "github.com/joomcode/errorx"

var ConfigErrors = errorx.NewNamespace("config")
var ConfigFileNotFoundError = ConfigErrors.NewType("config_file_not_found")
var ConfigFileParseError = ConfigErrors.NewType("config_file_parse_error")
var ConfigValueNotFoundError = ConfigErrors.NewType("config_value_not_found")
var CannotWriteConfigFileError = ConfigErrors.NewType("cannot_write_config_file")
