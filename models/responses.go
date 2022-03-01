package models

import (
	"context"
	errs "github.com/ONSdigital/dp-areas-api/apierrors"
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

func NewDBReadError(ctx context.Context, err error) *ErrorResponse {
	if err.Error() == errs.ErrNoRows.Error() {
		responseErr := NewError(ctx, err, InvalidAreaCodeError, err.Error())
		return NewErrorResponse(http.StatusNotFound, nil, responseErr)
	}
	responseErr := NewError(ctx, err, AreaDataIdGetError, err.Error())
	return NewErrorResponse(http.StatusInternalServerError, nil, responseErr)

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
