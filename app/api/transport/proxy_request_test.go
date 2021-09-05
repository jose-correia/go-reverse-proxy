package transport_test

import (
	"context"
	"fmt"
	"go-reverse-proxy/app/api/transport"
	"go-reverse-proxy/app/common/log"
	"go-reverse-proxy/app/values"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	proxyMock "go-reverse-proxy/mocks/app/handlers/proxy"

	"github.com/stretchr/testify/assert"
)

func TestProxyRequest(t *testing.T) {
	url := "http://127.0.0.1:5000/api/?parameter-key=test"
	responseBody := `{
  "message": "Hello World!", 
}`
	forwardRequestProviderMock := &proxyMock.HandlerMock{
		ForwardFunc: func(ctx context.Context, request *values.Request) ([]byte, int, error) {
			return []byte(responseBody), http.StatusOK, nil
		},
	}

	handler := transport.NewForwardRequest(
		log.NewNopLogger(),
		forwardRequestProviderMock,
	)

	req := httptest.NewRequest("PATCH", url, nil)
	req.Host = "service.com"

	w := httptest.NewRecorder()
	w.Header().Add("Content-Type", "application/json")
	handler.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	assert.Equal(t, "service.com", forwardRequestProviderMock.ForwardCalls()[0].Request.HostHeader)
	assert.Equal(t, "parameter-key=test", forwardRequestProviderMock.ForwardCalls()[0].Request.Parameters)
	assert.Len(t, forwardRequestProviderMock.ForwardCalls(), 1)

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, responseBody, string(body))
}

func TestProxyRequestError(t *testing.T) {
	url := "http://127.0.0.1:5000/api/"

	forwardRequestProviderMock := &proxyMock.HandlerMock{
		ForwardFunc: func(ctx context.Context, request *values.Request) ([]byte, int, error) {
			return []byte{}, http.StatusInternalServerError, fmt.Errorf("error")
		},
	}

	handler := transport.NewForwardRequest(
		log.NewNopLogger(),
		forwardRequestProviderMock,
	)

	req := httptest.NewRequest("GET", url, nil)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	assert.Len(t, forwardRequestProviderMock.ForwardCalls(), 1)
}
