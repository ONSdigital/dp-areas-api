// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/ONSdigital/dp-areas-api/api"
	"github.com/ONSdigital/dp-areas-api/models"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"sync"
)

// Ensure, that AreaStoreMock does implement api.AreaStore.
// If this is not the case, regenerate this file with moq.
var _ api.AreaStore = &AreaStoreMock{}

// AreaStoreMock is a mock implementation of api.AreaStore.
//
// 	func TestSomethingThatUsesAreaStore(t *testing.T) {
//
// 		// make and configure a mocked api.AreaStore
// 		mockedAreaStore := &AreaStoreMock{
// 			CheckerFunc: func(contextMoqParam context.Context, checkState *healthcheck.CheckState) error {
// 				panic("mock out the Checker method")
// 			},
// 			CloseFunc: func(ctx context.Context) error {
// 				panic("mock out the Close method")
// 			},
// 			GetAreaFunc: func(ctx context.Context, id string) (*models.Area, error) {
// 				panic("mock out the GetArea method")
// 			},
// 			GetAreasFunc: func(ctx context.Context, offset int, limit int) (*models.AreasResults, error) {
// 				panic("mock out the GetAreas method")
// 			},
// 			GetVersionFunc: func(ctx context.Context, id string, versionID int) (*models.Area, error) {
// 				panic("mock out the GetVersion method")
// 			},
// 		}
//
// 		// use mockedAreaStore in code that requires api.AreaStore
// 		// and then make assertions.
//
// 	}
type AreaStoreMock struct {
	// CheckerFunc mocks the Checker method.
	CheckerFunc func(contextMoqParam context.Context, checkState *healthcheck.CheckState) error

	// CloseFunc mocks the Close method.
	CloseFunc func(ctx context.Context) error

	// GetAreaFunc mocks the GetArea method.
	GetAreaFunc func(ctx context.Context, id string) (*models.Area, error)

	// GetAreasFunc mocks the GetAreas method.
	GetAreasFunc func(ctx context.Context, offset int, limit int) (*models.AreasResults, error)

	// GetVersionFunc mocks the GetVersion method.
	GetVersionFunc func(ctx context.Context, id string, versionID int) (*models.Area, error)

	// calls tracks calls to the methods.
	calls struct {
		// Checker holds details about calls to the Checker method.
		Checker []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// CheckState is the checkState argument value.
			CheckState *healthcheck.CheckState
		}
		// Close holds details about calls to the Close method.
		Close []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
		// GetArea holds details about calls to the GetArea method.
		GetArea []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID string
		}
		// GetAreas holds details about calls to the GetAreas method.
		GetAreas []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Offset is the offset argument value.
			Offset int
			// Limit is the limit argument value.
			Limit int
		}
		// GetVersion holds details about calls to the GetVersion method.
		GetVersion []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID string
			// VersionID is the versionID argument value.
			VersionID int
		}
	}
	lockChecker    sync.RWMutex
	lockClose      sync.RWMutex
	lockGetArea    sync.RWMutex
	lockGetAreas   sync.RWMutex
	lockGetVersion sync.RWMutex
}

// Checker calls CheckerFunc.
func (mock *AreaStoreMock) Checker(contextMoqParam context.Context, checkState *healthcheck.CheckState) error {
	if mock.CheckerFunc == nil {
		panic("AreaStoreMock.CheckerFunc: method is nil but AreaStore.Checker was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		CheckState      *healthcheck.CheckState
	}{
		ContextMoqParam: contextMoqParam,
		CheckState:      checkState,
	}
	mock.lockChecker.Lock()
	mock.calls.Checker = append(mock.calls.Checker, callInfo)
	mock.lockChecker.Unlock()
	return mock.CheckerFunc(contextMoqParam, checkState)
}

// CheckerCalls gets all the calls that were made to Checker.
// Check the length with:
//     len(mockedAreaStore.CheckerCalls())
func (mock *AreaStoreMock) CheckerCalls() []struct {
	ContextMoqParam context.Context
	CheckState      *healthcheck.CheckState
} {
	var calls []struct {
		ContextMoqParam context.Context
		CheckState      *healthcheck.CheckState
	}
	mock.lockChecker.RLock()
	calls = mock.calls.Checker
	mock.lockChecker.RUnlock()
	return calls
}

// Close calls CloseFunc.
func (mock *AreaStoreMock) Close(ctx context.Context) error {
	if mock.CloseFunc == nil {
		panic("AreaStoreMock.CloseFunc: method is nil but AreaStore.Close was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockClose.Lock()
	mock.calls.Close = append(mock.calls.Close, callInfo)
	mock.lockClose.Unlock()
	return mock.CloseFunc(ctx)
}

// CloseCalls gets all the calls that were made to Close.
// Check the length with:
//     len(mockedAreaStore.CloseCalls())
func (mock *AreaStoreMock) CloseCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockClose.RLock()
	calls = mock.calls.Close
	mock.lockClose.RUnlock()
	return calls
}

// GetArea calls GetAreaFunc.
func (mock *AreaStoreMock) GetArea(ctx context.Context, id string) (*models.Area, error) {
	if mock.GetAreaFunc == nil {
		panic("AreaStoreMock.GetAreaFunc: method is nil but AreaStore.GetArea was just called")
	}
	callInfo := struct {
		Ctx context.Context
		ID  string
	}{
		Ctx: ctx,
		ID:  id,
	}
	mock.lockGetArea.Lock()
	mock.calls.GetArea = append(mock.calls.GetArea, callInfo)
	mock.lockGetArea.Unlock()
	return mock.GetAreaFunc(ctx, id)
}

// GetAreaCalls gets all the calls that were made to GetArea.
// Check the length with:
//     len(mockedAreaStore.GetAreaCalls())
func (mock *AreaStoreMock) GetAreaCalls() []struct {
	Ctx context.Context
	ID  string
} {
	var calls []struct {
		Ctx context.Context
		ID  string
	}
	mock.lockGetArea.RLock()
	calls = mock.calls.GetArea
	mock.lockGetArea.RUnlock()
	return calls
}

// GetAreas calls GetAreasFunc.
func (mock *AreaStoreMock) GetAreas(ctx context.Context, offset int, limit int) (*models.AreasResults, error) {
	if mock.GetAreasFunc == nil {
		panic("AreaStoreMock.GetAreasFunc: method is nil but AreaStore.GetAreas was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		Offset int
		Limit  int
	}{
		Ctx:    ctx,
		Offset: offset,
		Limit:  limit,
	}
	mock.lockGetAreas.Lock()
	mock.calls.GetAreas = append(mock.calls.GetAreas, callInfo)
	mock.lockGetAreas.Unlock()
	return mock.GetAreasFunc(ctx, offset, limit)
}

// GetAreasCalls gets all the calls that were made to GetAreas.
// Check the length with:
//     len(mockedAreaStore.GetAreasCalls())
func (mock *AreaStoreMock) GetAreasCalls() []struct {
	Ctx    context.Context
	Offset int
	Limit  int
} {
	var calls []struct {
		Ctx    context.Context
		Offset int
		Limit  int
	}
	mock.lockGetAreas.RLock()
	calls = mock.calls.GetAreas
	mock.lockGetAreas.RUnlock()
	return calls
}

// GetVersion calls GetVersionFunc.
func (mock *AreaStoreMock) GetVersion(ctx context.Context, id string, versionID int) (*models.Area, error) {
	if mock.GetVersionFunc == nil {
		panic("AreaStoreMock.GetVersionFunc: method is nil but AreaStore.GetVersion was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		ID        string
		VersionID int
	}{
		Ctx:       ctx,
		ID:        id,
		VersionID: versionID,
	}
	mock.lockGetVersion.Lock()
	mock.calls.GetVersion = append(mock.calls.GetVersion, callInfo)
	mock.lockGetVersion.Unlock()
	return mock.GetVersionFunc(ctx, id, versionID)
}

// GetVersionCalls gets all the calls that were made to GetVersion.
// Check the length with:
//     len(mockedAreaStore.GetVersionCalls())
func (mock *AreaStoreMock) GetVersionCalls() []struct {
	Ctx       context.Context
	ID        string
	VersionID int
} {
	var calls []struct {
		Ctx       context.Context
		ID        string
		VersionID int
	}
	mock.lockGetVersion.RLock()
	calls = mock.calls.GetVersion
	mock.lockGetVersion.RUnlock()
	return calls
}
