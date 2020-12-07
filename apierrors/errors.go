package apierrors

import (
	"errors"
)

//!!! fix this list for this service ...

// A list of error messages for Topic API
var (
	ErrTopicNotFound                    = errors.New("topic not found")
	ErrVariantNotFound                  = errors.New("image download variant not found")
	ErrVariantAlreadyExists             = errors.New("image download variant already exists")
	ErrInternalServer                   = errors.New("internal error")
	ErrUnableToReadMessage              = errors.New("failed to read message body")
	ErrImageIDMismatch                  = errors.New("image id provided in body does not match 'id' path parameter")
	ErrUnableToParseJSON                = errors.New("failed to parse json body")
	ErrImageFilenameTooLong             = errors.New("image filename is too long")
	ErrImageNoCollectionID              = errors.New("image does not have a collectionID")
	ErrImageAlreadyPublished            = errors.New("image is already published")
	ErrTopicAlreadyCompleted            = errors.New("topic is already completed")
	ErrTopicInvalidState                = errors.New("topic state is not a valid state name")
	ErrImageBadInitialState             = errors.New("image state is not a valid initial state")
	ErrTopicStateTransitionNotAllowed   = errors.New("topic state transition not allowed")
	ErrTopicUploadEmpty                 = errors.New("topic upload section is not populated")
	ErrImageUploadPathEmpty             = errors.New("image upload path is not populated")
	ErrImageNotImporting                = errors.New("image is not in importing state")
	ErrImageNotPublished                = errors.New("image is not in published state")
	ErrVariantIDMismatch                = errors.New("variant id provided in body does not match 'variant' path parameter")
	ErrVariantStateTransitionNotAllowed = errors.New("image download variant state transition not allowed")
	ErrImageDownloadTypeMismatch        = errors.New("image download variant type does not match existing type")
	ErrImageDownloadInvalidState        = errors.New("image download state is not a valid state name")
	ErrImageDownloadBadInitialState     = errors.New("image download state is not a valid initial state")
)