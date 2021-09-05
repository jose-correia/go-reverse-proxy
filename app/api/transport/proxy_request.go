package transport

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	encoder "go-reverse-proxy/app/common/encoder"
	"go-reverse-proxy/app/values"

	"github.com/go-kit/kit/log"
)

type forwardRequestHTTPProvider interface {
	Forward(
		ctx context.Context,
		request *values.Request,
	) (
		[]byte, int, error)
}

type forwardRequestHTTPHandler struct {
	logger   log.Logger
	provider forwardRequestHTTPProvider
}

func NewForwardRequest(logger log.Logger, provider forwardRequestHTTPProvider) *forwardRequestHTTPHandler {
	return &forwardRequestHTTPHandler{logger: logger, provider: provider}
}

func (c *forwardRequestHTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	payload, err := ioutil.ReadAll(req.Body)
	if err != nil {
		c.logger.Log("transport", "proxyRequest/HTTP", "error", err.Error())
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	fmt.Println(req.URL.RawQuery)
	response, statusCode, err := c.provider.Forward(
		req.Context(),
		&values.Request{
			Method:     req.Method,
			Header:     req.Header,
			HostHeader: req.Host,
			Parameters: req.URL.RawQuery,
			Payload:    payload,
		},
	)
	if err != nil {
		c.logger.Log("transport", "proxyRequest/HTTP", "error", err.Error())
		encoder.Encode(req.Context(), &encoder.Error{Code: statusCode, Message: err.Error()}, w)
		return
	}

	w.WriteHeader(statusCode)
	_, err = io.Copy(w, bytes.NewReader(response))
	if err != nil {
		c.logger.Log("transport", "proxyRequest/HTTP", "error", err.Error())
		encoder.Encode(req.Context(), &encoder.Error{Code: http.StatusInternalServerError, Message: err.Error()}, w)
		return
	}
	c.logger.Log("transport", "proxyRequest/HTTP")
}
