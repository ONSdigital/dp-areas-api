package mocks

import (
	"net/http"

	"github.com/ONSdigital/dp-authorisation/auth"
)

type PermissionCheckCalls struct {
	Calls int
}

type AuthHandlerMock struct {
	Required *PermissionCheckCalls
}

func NewAuthHandlerMock() *AuthHandlerMock {
	return &AuthHandlerMock{
		Required: &PermissionCheckCalls{
			Calls: 0,
		},
	}
}

func (a AuthHandlerMock) Require(required auth.Permissions, handler http.HandlerFunc) http.HandlerFunc {
	return a.Required.checkPermissions(handler)
}

func (c *PermissionCheckCalls) checkPermissions(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.Calls += 1
		h.ServeHTTP(w, r)
	}
}
