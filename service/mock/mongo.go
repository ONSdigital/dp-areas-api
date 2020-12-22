// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/ONSdigital/dp-areas-api/service"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"sync"
)

// Ensure, that MongoServerMock does implement service.MongoServer.
// If this is not the case, regenerate this file with moq.
var _ service.MongoServer = &MongoServerMock{}

// MongoServerMock is a mock implementation of service.MongoServer.
//
//     func TestSomethingThatUsesMongoServer(t *testing.T) {
//
//         // make and configure a mocked service.MongoServer
//         mockedMongoServer := &MongoServerMock{
//             CheckerFunc: func(in1 context.Context, in2 *healthcheck.CheckState) error {
// 	               panic("mock out the Checker method")
//             },
//             CloseFunc: func(in1 context.Context) error {
// 	               panic("mock out the Close method")
//             },
//         }
//
//         // use mockedMongoServer in code that requires service.MongoServer
//         // and then make assertions.
//
//     }
type MongoServerMock struct {
	// CheckerFunc mocks the Checker method.
	CheckerFunc func(in1 context.Context, in2 *healthcheck.CheckState) error

	// CloseFunc mocks the Close method.
	CloseFunc func(in1 context.Context) error

	// calls tracks calls to the methods.
	calls struct {
		// Checker holds details about calls to the Checker method.
		Checker []struct {
			// In1 is the in1 argument value.
			In1 context.Context
			// In2 is the in2 argument value.
			In2 *healthcheck.CheckState
		}
		// Close holds details about calls to the Close method.
		Close []struct {
			// In1 is the in1 argument value.
			In1 context.Context
		}
	}
	lockChecker sync.RWMutex
	lockClose   sync.RWMutex
}

// Checker calls CheckerFunc.
func (mock *MongoServerMock) Checker(in1 context.Context, in2 *healthcheck.CheckState) error {
	if mock.CheckerFunc == nil {
		panic("MongoServerMock.CheckerFunc: method is nil but MongoServer.Checker was just called")
	}
	callInfo := struct {
		In1 context.Context
		In2 *healthcheck.CheckState
	}{
		In1: in1,
		In2: in2,
	}
	mock.lockChecker.Lock()
	mock.calls.Checker = append(mock.calls.Checker, callInfo)
	mock.lockChecker.Unlock()
	return mock.CheckerFunc(in1, in2)
}

// CheckerCalls gets all the calls that were made to Checker.
// Check the length with:
//     len(mockedMongoServer.CheckerCalls())
func (mock *MongoServerMock) CheckerCalls() []struct {
	In1 context.Context
	In2 *healthcheck.CheckState
} {
	var calls []struct {
		In1 context.Context
		In2 *healthcheck.CheckState
	}
	mock.lockChecker.RLock()
	calls = mock.calls.Checker
	mock.lockChecker.RUnlock()
	return calls
}

// Close calls CloseFunc.
func (mock *MongoServerMock) Close(in1 context.Context) error {
	if mock.CloseFunc == nil {
		panic("MongoServerMock.CloseFunc: method is nil but MongoServer.Close was just called")
	}
	callInfo := struct {
		In1 context.Context
	}{
		In1: in1,
	}
	mock.lockClose.Lock()
	mock.calls.Close = append(mock.calls.Close, callInfo)
	mock.lockClose.Unlock()
	return mock.CloseFunc(in1)
}

// CloseCalls gets all the calls that were made to Close.
// Check the length with:
//     len(mockedMongoServer.CloseCalls())
func (mock *MongoServerMock) CloseCalls() []struct {
	In1 context.Context
} {
	var calls []struct {
		In1 context.Context
	}
	mock.lockClose.RLock()
	calls = mock.calls.Close
	mock.lockClose.RUnlock()
	return calls
}
