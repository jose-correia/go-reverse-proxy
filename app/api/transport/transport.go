package transport

import (
	"go-reverse-proxy/app/api"
	"go-reverse-proxy/app/handlers/proxy"

	"github.com/go-kit/kit/log"
	httpkit "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func BuildEndpointRegister(
	logger log.Logger,
	svc proxy.Handler,
) api.EndpointRegister {
	return func(router *mux.Router, options ...httpkit.ServerOption) {
		authSubRouter := router.PathPrefix("/api").Subrouter()
		authSubRouter.
			Path("/").
			Handler(NewForwardRequest(logger, svc))
	}
}
