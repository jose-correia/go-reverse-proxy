// Package httpclient contains a destiny-abstract HTTP client
package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
)

type HttpClient interface {
	// Request can be used to send an HTTP request to any given destination.
	Request(
		ctx context.Context,
		method string,
		address string,
		header http.Header,
		parameters string,
		payload []byte) ([]byte, int, error)
	// GetHttpCLient is a getter for the base *net.http struct so that
	// it can be wrapper in other modules
	GetHttpClient() *http.Client
}

type defaultClient struct {
	requestTimeout time.Duration
	httpClient     *http.Client
	logger         log.Logger
}

const (
	defaultRetryDelayMin = 50 * time.Millisecond
	defaultRetryMax      = 3
)

func New(
	logger log.Logger,
	timeout time.Duration,
	httpClient *http.Client,
) HttpClient {
	// Instantiate a new HTTP clients that implements retryable
	// HTTP requests
	retryableClient := retryablehttp.NewClient()
	retryableClient.RetryWaitMin = defaultRetryDelayMin
	retryableClient.RetryMax = defaultRetryMax

	// This detail is used to inject a Mockable HTTP client during tests
	if httpClient != nil {
		retryableClient.HTTPClient = httpClient
	}

	httpClient = retryableClient.StandardClient()

	var svc HttpClient
	svc = &defaultClient{
		requestTimeout: timeout,
		httpClient:     httpClient,
		logger:         logger,
	}
	return svc
}

func (c *defaultClient) Request(
	ctx context.Context,
	method string,
	address string,
	header http.Header,
	parameters string,
	payload []byte,
) ([]byte, int, error) {
	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	url := c.buildURL(address, parameters)

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(payload))
	if err != nil {
		c.logger.Log("module", "httpclient", "payload", payload, "err", err, "step", "http.NewRequest")
		return nil, http.StatusInternalServerError, err
	}

	req.Header = header

	res, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Log("module", "httpclient", "payload", payload, "err", err, "step", "http.Do")
		return nil,
			http.StatusInternalServerError,
			errors.Wrapf(err, "failed to request service url: /%s", url)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.logger.Log("module", "httpclient", "payload", payload, "err", err, "step", "ioutil.ReadAll")
		return nil,
			res.StatusCode,
			errors.Wrapf(err, "failed to decode response from service url: /%s", url)
	}

	c.logger.Log("module", "httpclient", "request", url)
	return body, res.StatusCode, nil
}

func (c *defaultClient) buildURL(address string, parameters string) string {
	// In order to have all the HTTP logic encapsulated in this service,
	// and since we know that hosts will be configured by their IP address,
	// we add the http prefix to the IP here.
	httpAdress := fmt.Sprintf("http://%s", address)

	if parameters != "" {
		return fmt.Sprintf("%s?%s", httpAdress, parameters)
	}

	return httpAdress
}

func (c *defaultClient) GetHttpClient() *http.Client {
	return c.httpClient
}
