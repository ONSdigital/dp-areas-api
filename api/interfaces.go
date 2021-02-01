package api

import (
	"net/http"

	"github.com/ONSdigital/dp-authorisation/auth"
)

//go:generate moq -out mock/auth.go -pkg mock . AuthHandler

// AuthHandler interface for adding auth to endpoints
type AuthHandler interface {
	Require(required auth.Permissions, handler http.HandlerFunc) http.HandlerFunc
}
