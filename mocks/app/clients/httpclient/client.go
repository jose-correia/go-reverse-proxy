// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package db_mock

import (
	"context"
	"go-reverse-proxy/app/clients/httpclient"
	"net/http"
	"sync"
)

// Ensure, that HttpClientMock does implement httpclient.HttpClient.
// If this is not the case, regenerate this file with moq.
var _ httpclient.HttpClient = &HttpClientMock{}

// HttpClientMock is a mock implementation of httpclient.HttpClient.
//
// 	func TestSomethingThatUsesHttpClient(t *testing.T) {
//
// 		// make and configure a mocked httpclient.HttpClient
// 		mockedHttpClient := &HttpClientMock{
// 			GetHttpClientFunc: func() *http.Client {
// 				panic("mock out the GetHttpClient method")
// 			},
// 			RequestFunc: func(ctx context.Context, method string, address string, header http.Header, parameters string, payload []byte) ([]byte, int, error) {
// 				panic("mock out the Request method")
// 			},
// 		}
//
// 		// use mockedHttpClient in code that requires httpclient.HttpClient
// 		// and then make assertions.
//
// 	}
type HttpClientMock struct {
	// GetHttpClientFunc mocks the GetHttpClient method.
	GetHttpClientFunc func() *http.Client

	// RequestFunc mocks the Request method.
	RequestFunc func(ctx context.Context, method string, address string, header http.Header, parameters string, payload []byte) ([]byte, int, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetHttpClient holds details about calls to the GetHttpClient method.
		GetHttpClient []struct {
		}
		// Request holds details about calls to the Request method.
		Request []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Method is the method argument value.
			Method string
			// Address is the address argument value.
			Address string
			// Header is the header argument value.
			Header http.Header
			// Parameters is the parameters argument value.
			Parameters string
			// Payload is the payload argument value.
			Payload []byte
		}
	}
	lockGetHttpClient sync.RWMutex
	lockRequest       sync.RWMutex
}

// GetHttpClient calls GetHttpClientFunc.
func (mock *HttpClientMock) GetHttpClient() *http.Client {
	if mock.GetHttpClientFunc == nil {
		panic("HttpClientMock.GetHttpClientFunc: method is nil but HttpClient.GetHttpClient was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetHttpClient.Lock()
	mock.calls.GetHttpClient = append(mock.calls.GetHttpClient, callInfo)
	mock.lockGetHttpClient.Unlock()
	return mock.GetHttpClientFunc()
}

// GetHttpClientCalls gets all the calls that were made to GetHttpClient.
// Check the length with:
//     len(mockedHttpClient.GetHttpClientCalls())
func (mock *HttpClientMock) GetHttpClientCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetHttpClient.RLock()
	calls = mock.calls.GetHttpClient
	mock.lockGetHttpClient.RUnlock()
	return calls
}

// Request calls RequestFunc.
func (mock *HttpClientMock) Request(ctx context.Context, method string, address string, header http.Header, parameters string, payload []byte) ([]byte, int, error) {
	if mock.RequestFunc == nil {
		panic("HttpClientMock.RequestFunc: method is nil but HttpClient.Request was just called")
	}
	callInfo := struct {
		Ctx        context.Context
		Method     string
		Address    string
		Header     http.Header
		Parameters string
		Payload    []byte
	}{
		Ctx:        ctx,
		Method:     method,
		Address:    address,
		Header:     header,
		Parameters: parameters,
		Payload:    payload,
	}
	mock.lockRequest.Lock()
	mock.calls.Request = append(mock.calls.Request, callInfo)
	mock.lockRequest.Unlock()
	return mock.RequestFunc(ctx, method, address, header, parameters, payload)
}

// RequestCalls gets all the calls that were made to Request.
// Check the length with:
//     len(mockedHttpClient.RequestCalls())
func (mock *HttpClientMock) RequestCalls() []struct {
	Ctx        context.Context
	Method     string
	Address    string
	Header     http.Header
	Parameters string
	Payload    []byte
} {
	var calls []struct {
		Ctx        context.Context
		Method     string
		Address    string
		Header     http.Header
		Parameters string
		Payload    []byte
	}
	mock.lockRequest.RLock()
	calls = mock.calls.Request
	mock.lockRequest.RUnlock()
	return calls
}