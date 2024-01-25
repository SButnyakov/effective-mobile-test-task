package services

import "errors"

var (
	ErrExternalError = errors.New("an error occurred while trying to use external resource")
)
