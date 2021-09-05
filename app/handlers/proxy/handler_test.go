package proxy_test

import (
	"context"
	"go-reverse-proxy/app/common/log"
	"go-reverse-proxy/app/handlers/loadbalancing"
	"go-reverse-proxy/app/handlers/proxy"
	"go-reverse-proxy/app/values"
	"net/http"

	http_mock "go-reverse-proxy/mocks/app/clients/httpclient"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newProxyHandler(
	configuration *values.Configuration,
) (proxy.Handler, *http_mock.HttpClientMock, loadbalancing.Handler) {
	logger := log.NewLogger()

	httpClient := &http_mock.HttpClientMock{}
	loadBalancer := loadbalancing.New(logger)

	return proxy.New(
		logger,
		*configuration,
		httpClient,
		loadBalancer,
	), httpClient, loadBalancer
}

func TestForward(t *testing.T) {
	configuration := &values.Configuration{
		Host: &values.Host{
			Address: "127.0.0.1",
			Port:    8080,
		},
		Services: map[string]*values.Service{
			"my-domain.com": {
				Name:   "my-service",
				Domain: "my-domain.com",
				Hosts: []*values.Host{
					{
						Address: "127.0.0.1",
						Port:    5000,
					},
					{
						Address: "127.0.0.1",
						Port:    5001,
					},
				},
				NextHostIndex: 0,
			},
		},
		RetryableStatusCodes: []int{},
		MaxForwardRetries:    0,
	}

	handler, httpClient, _ := newProxyHandler(
		configuration,
	)

	httpClient.RequestFunc = func(
		ctx context.Context,
		method string,
		address string,
		header http.Header,
		parameters string,
		payload []byte,
	) ([]byte, int, error) {
		return []byte{}, http.StatusOK, nil
	}

	response, status, err := handler.Forward(
		context.Background(),
		&values.Request{
			Method:     "GET",
			Header:     http.Header{},
			HostHeader: "my-domain.com",
			Parameters: "",
			Payload:    []byte{},
		},
	)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, []byte{}, response)
	assert.Equal(t, 1, len(httpClient.RequestCalls()))
	assert.Equal(t, "127.0.0.1:5000", httpClient.RequestCalls()[0].Address)
	assert.Equal(t, "GET", httpClient.RequestCalls()[0].Method)
}

func TestForwardServiceNotFound(t *testing.T) {
	configuration := &values.Configuration{
		Host: &values.Host{
			Address: "127.0.0.1",
			Port:    8080,
		},
		Services: map[string]*values.Service{},
	}

	handler, _, _ := newProxyHandler(
		configuration,
	)

	response, status, err := handler.Forward(
		context.Background(),
		&values.Request{
			Method:     "GET",
			Header:     http.Header{},
			HostHeader: "my-domain.com",
			Parameters: "",
			Payload:    []byte{},
		},
	)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, status)
	assert.Equal(t, []byte{}, response)
}

func TestForwardWithRetriesExceeded(t *testing.T) {
	configuration := &values.Configuration{
		Host: &values.Host{
			Address: "127.0.0.1",
			Port:    8080,
		},
		Services: map[string]*values.Service{
			"my-domain.com": {
				Name:   "my-service",
				Domain: "my-domain.com",
				Hosts: []*values.Host{
					{
						Address: "127.0.0.1",
						Port:    5000,
					},
					{
						Address: "127.0.0.1",
						Port:    5001,
					},
				},
				NextHostIndex: 0,
			},
		},
		RetryableStatusCodes: []int{http.StatusInternalServerError},
		MaxForwardRetries:    2,
	}

	handler, httpClient, _ := newProxyHandler(
		configuration,
	)

	httpClient.RequestFunc = func(
		ctx context.Context,
		method string,
		address string,
		header http.Header,
		parameters string,
		payload []byte,
	) ([]byte, int, error) {
		return []byte{}, http.StatusInternalServerError, nil
	}

	response, status, err := handler.Forward(
		context.Background(),
		&values.Request{
			Method:     "GET",
			Header:     http.Header{},
			HostHeader: "my-domain.com",
			Parameters: "",
			Payload:    []byte{},
		},
	)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, status)
	assert.Equal(t, []byte{}, response)
	assert.Equal(t, 3, len(httpClient.RequestCalls()))
	assert.Equal(t, "127.0.0.1:5000", httpClient.RequestCalls()[0].Address)
	assert.Equal(t, "127.0.0.1:5001", httpClient.RequestCalls()[1].Address)
	assert.Equal(t, "127.0.0.1:5000", httpClient.RequestCalls()[2].Address)
}
