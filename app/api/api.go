package api

import (
	"go-reverse-proxy/app/common/middlewares"
	"net/http"

	"github.com/go-kit/kit/log"
	httpkit "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

type API struct{ Router *mux.Router }

type EndpointRegister = func(router *mux.Router, options ...httpkit.ServerOption)

func New(
	logger log.Logger,
	externalEndpointRegister EndpointRegister,
) API {
	router := mux.NewRouter().StrictSlash(false)
	router.Use(
		middlewares.Logger(logger),
	)

	apiRouter := router.PathPrefix("/").Subrouter()
	externalEndpointRegister(apiRouter)

	return API{
		Router: router,
	}
}

func (api API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.Router.ServeHTTP(w, r)
}
