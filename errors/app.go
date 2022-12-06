package errors

import "github.com/joomcode/errorx"

var AppErrors = errorx.NewNamespace("app")
var UserNotFoundError = AppErrors.NewType("user_not_found")
var UserAlreadyExistsError = AppErrors.NewType("user_already_exists")
var UserNoProfileError = AppErrors.NewType("user_no_profile")