package mocks

import (
	"net/http"
	"sync"

	"github.com/ONSdigital/dp-authorisation/auth"
)

var (
	muAuthHandlerMock sync.Mutex
)

type PermissionCheckCalls struct {
	Calls int
}

type AuthHandlerMock struct {
	Required *PermissionCheckCalls
}

func NewAuthHandlerMock() *AuthHandlerMock {
	muAuthHandlerMock.Lock()
	defer muAuthHandlerMock.Unlock()

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
	muAuthHandlerMock.Lock()
	defer muAuthHandlerMock.Unlock()

	return func(w http.ResponseWriter, r *http.Request) {
		c.Calls++
		h.ServeHTTP(w, r)
	}
}
