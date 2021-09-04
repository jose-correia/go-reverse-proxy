// +build unit

package httpclient_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-reverse-proxy/app/clients/httpclient"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

type failedRoundTripFunc func(req *http.Request) *http.Response

func (f failedRoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("failed to perform HTTP request")
}

func newHTTPClient(fn roundTripFunc, failRequests bool) *http.Client {
	if failRequests {
		return &http.Client{
			Transport: failedRoundTripFunc(fn),
		}
	}

	return &http.Client{
		Transport: roundTripFunc(fn),
	}
}

func TestRequest(t *testing.T) {
	expectedMethod := "GET"
	expectedURL := `http://127.0.0.1:8080`
	expectedHeader := http.Header{}
	expectedHeader.Add("Content-Type", "application/json")
	responseBody := `{
  "objects": [
    {
      "message": "Hello World!"
    },
  ]
}`

	// replace the *http.Client w/ one with overriden Transport
	mockHTTPClient := newHTTPClient(
		func(req *http.Request) *http.Response {
			buf, _ := ioutil.ReadAll(req.Body)
			var prettyJSON bytes.Buffer
			_ = json.Indent(&prettyJSON, buf, "", "  ")
			assert.Equal(t, req.Method, expectedMethod)
			assert.Equal(t, expectedURL, req.URL.String())
			assert.Equal(t, expectedHeader, req.Header)
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(responseBody)),
				Header:     make(http.Header),
			}
		},
		false,
	)

	httpClient := httpclient.New(log.NewNopLogger(), 5*time.Second, mockHTTPClient)

	var reqPayload []byte
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	resp, statusCode, err := httpClient.Request(
		context.TODO(),
		"GET",
		"127.0.0.1:8080",
		header,
		"",
		reqPayload)

	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Nil(t, err)
}

func TestRequestWithParameters(t *testing.T) {
	expectedMethod := "GET"
	expectedURL := `http://127.0.0.1:8080?par1=test&par2=test`
	expectedHeader := http.Header{}
	expectedHeader.Add("Content-Type", "application/json")
	responseBody := ""

	// replace the *http.Client w/ one with overriden Transport
	mockHTTPClient := newHTTPClient(
		func(req *http.Request) *http.Response {
			buf, _ := ioutil.ReadAll(req.Body)
			var prettyJSON bytes.Buffer
			_ = json.Indent(&prettyJSON, buf, "", "  ")
			assert.Equal(t, req.Method, expectedMethod)
			assert.Equal(t, expectedURL, req.URL.String())
			assert.Equal(t, expectedHeader, req.Header)
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(responseBody)),
				Header:     make(http.Header),
			}
		},
		false,
	)

	httpClient := httpclient.New(log.NewNopLogger(), 5*time.Second, mockHTTPClient)

	var reqPayload []byte
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	resp, statusCode, err := httpClient.Request(
		context.TODO(),
		"GET",
		"127.0.0.1:8080",
		header,
		"par1=test&par2=test",
		reqPayload)

	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Nil(t, err)
}

func TestRequestAddressIsDown(t *testing.T) {
	// replace the *http.Client w/ one with overriden Transport
	mockHTTPClient := newHTTPClient(
		nil,
		true,
	)

	httpClient := httpclient.New(log.NewNopLogger(), 5*time.Second, mockHTTPClient)

	var reqPayload []byte
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	resp, statusCode, err := httpClient.Request(
		context.TODO(),
		"GET",
		"127.0.0.1:8080",
		header,
		"",
		reqPayload)

	assert.Nil(t, resp)
	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.NotNil(t, err)
}
