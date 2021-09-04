package transport

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"

	encoder "go-reverse-proxy/app/common/encoder"

	"github.com/go-kit/kit/log"
)

type forwardRequestHTTPProvider interface {
	Forward(
		ctx context.Context,
		method string,
		header http.Header,
		hostHeader string,
		parameters string,
		payload []byte,
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

	response, statusCode, err := c.provider.Forward(
		req.Context(),
		req.Method,
		req.Header,
		req.Host,
		req.URL.RawQuery,
		payload,
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
		encoder.Encode(req.Context(), &encoder.Error{Code: 500, Message: err.Error()}, w)
		return
	}
	c.logger.Log("transport", "proxyRequest/HTTP")
}
