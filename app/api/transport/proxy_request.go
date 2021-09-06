package transport

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	encoder "go-reverse-proxy/app/common/encoder"
	"go-reverse-proxy/app/values"

	"github.com/go-kit/kit/log"
)

type forwardRequestHTTPProvider interface {
	Forward(
		ctx context.Context,
		request *values.Request,
	) ([]byte, int, error)
}

type forwardRequestHTTPHandler struct {
	logger      log.Logger
	provider    forwardRequestHTTPProvider
	routePrefix string
}

func NewForwardRequest(
	logger log.Logger,
	provider forwardRequestHTTPProvider,
	routePrefix string,
) *forwardRequestHTTPHandler {
	return &forwardRequestHTTPHandler{
		logger:      logger,
		provider:    provider,
		routePrefix: routePrefix,
	}
}

func (c *forwardRequestHTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	pathSplit := strings.Split(req.URL.Path, c.routePrefix)

	var endpoint string
	if len(pathSplit) > 1 {
		endpoint = pathSplit[1]
	}

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
		&values.Request{
			Method:     req.Method,
			Endpoint:   endpoint,
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
