package apierrors

import (
	"errors"
)

// A list of error messages for Topic API
var (
	ErrTopicNotFound                  = errors.New("topic not found")
	ErrContentNotFound                = errors.New("content not found")
	ErrNotFound                       = errors.New("not found")
	ErrInternalServer                 = errors.New("internal error")
	ErrUnableToReadMessage            = errors.New("failed to read message body")
	ErrUnableToParseJSON              = errors.New("failed to parse json body")
	ErrTopicInvalidState              = errors.New("topic state is not a valid state name")
	ErrTopicStateTransitionNotAllowed = errors.New("topic state transition not allowed")
	ErrTopicUploadEmpty               = errors.New("topic upload section is not populated")
)
