package utils

import (
	"strconv"

	errs "github.com/ONSdigital/dp-areas-api/apierrors"
)

// ValidatePositiveInt obtains the positive int value of query
func ValidatePositiveInt(parameter string) (val int, err error) {
	val, err = strconv.Atoi(parameter)
	if err != nil {
		return -1, errs.ErrInvalidQueryParameter
	}
	if val < 0 {
		return -1, errs.ErrInvalidQueryParameter
	}
	return val, nil
}
