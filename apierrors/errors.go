package apierrors

import (
	"errors"
)

//A lit of error messages for Permissions API
var (
	//ErrAreaNotFound is an error when an Area cannot be found in mongoDB
	ErrAreaNotFound      = errors.New("Area not found")
	ErrVersionNotFound   = errors.New("Version not found")
	ErrInternalServer    = errors.New("internal error")
)
