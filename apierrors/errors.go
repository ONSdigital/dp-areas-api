package apierrors

import (
	"errors"
)

// A list of error messages for Topic API
var (
	ErrContentNotFound                = errors.New("content not found")
	ErrContentUnrecognisedParameter   = errors.New("content query not recognised")
	ErrEmptyRequestBody               = errors.New("request body empty")
	ErrInternalServer                 = errors.New("internal error")
	ErrInvalidReleaseDate             = errors.New("invalid topic release date, must have the following format: 2022-05-22T09:21:45Z")
	ErrNotFound                       = errors.New("not found")
	ErrTopicInvalidState              = errors.New("topic state is not a valid state name")
	ErrTopicNotFound                  = errors.New("topic not found")
	ErrTopicStateTransitionNotAllowed = errors.New("topic state transition not allowed")
	ErrTopicUploadEmpty               = errors.New("topic upload section is not populated")
	ErrUnableToParseJSON              = errors.New("failed to parse json body")
	ErrUnableToReadMessage            = errors.New("failed to read message body")
)
