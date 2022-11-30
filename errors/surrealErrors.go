package errors

import "github.com/joomcode/errorx"

var SurrealErrors = errorx.NewNamespace("surreal")
var SurrealNotReachableError = SurrealErrors.NewType("not_reachable")
var SurrealDeserializaError = SurrealErrors.NewType("deserialization_error")
var SurrealExecError = SurrealErrors.NewType("exec_error")
var SurrealQueryError = SurrealErrors.NewType("query_error")
