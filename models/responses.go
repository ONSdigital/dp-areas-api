package models

import (
	"context"
	"net/http"
)

type ErrorResponse struct {
	Errors  []error           `json:"errors"`
	Status  int               `json:"-"`
	Headers map[string]string `json:"-"`
}

func NewErrorResponse(statusCode int, headers map[string]string, errors ...error) *ErrorResponse {
	return &ErrorResponse{
		Errors:  errors,
		Status:  statusCode,
		Headers: headers,
	}
}

func NewBodyReadError(ctx context.Context, err error) *ErrorResponse {
	return NewErrorResponse(http.StatusInternalServerError,
		nil,
		NewError(ctx, err, BodyReadError, BodyReadFailedDescription),
	)
}

func NewBodyUnmarshalError(ctx context.Context, err error) *ErrorResponse {
	return NewErrorResponse(http.StatusInternalServerError,
		nil,
		NewError(ctx, err, JSONUnmarshalError, ErrorUnmarshalFailedDescription),
	)
}

type SuccessResponse struct {
	Body    []byte            `json:"-"`
	Status  int               `json:"-"`
	Headers map[string]string `json:"-"`
}

func NewSuccessResponse(jsonBody []byte, statusCode int, headers map[string]string) *SuccessResponse {
	return &SuccessResponse{
		Body:    jsonBody,
		Status:  statusCode,
		Headers: headers,
	}
}
