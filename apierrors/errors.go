package apierrors

import (
	"errors"
)

// A list of error messages for Areas API
var (
	ErrAreaNotFound             = errors.New("Area not found")
	ErrVersionNotFound          = errors.New("Version not found")
	ErrInternalServer           = errors.New("internal error")
	ErrInvalidQueryParameter    = errors.New("invalid query parameter")
	ErrQueryParamLimitExceedMax = errors.New("limit exceeds max value")
	ErrNoRows                   = errors.New("no rows in result set")
)
