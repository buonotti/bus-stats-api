package errors

import "github.com/joomcode/errorx"

var ApiErrors = errorx.NewNamespace("api")
var MalformedRequestError = ApiErrors.NewType("malformed_request")
var RouteNotFoundError = ApiErrors.NewType("route_not_found")
var UnauthorizedError = ApiErrors.NewType("unauthorized")
var TokenError = ApiErrors.NewType("token_error")