package services

import (
	"errors"
)

type ErrorResponse struct {
	Message string `json:"message" example:"an error has occured"`
}

var CredentialError = errors.New("invalid credentials")
var FileError = errors.New("bad file upload")
var FormatError = errors.New("unexpected data format")
