package configuration

import (
	"context"
	"log"
)

type Handler interface {
	ReadFromYaml(ctx context.Context, path string) (*Configuration, error)
}

type DefaultHandler struct {
	logger log.Logger
}

func New(logger log.Logger) Handler {
	var svc Handler
	svc = &DefaultHandler{
		logger: logger,
	}

	return svc
}

func (h *DefaultHandler) ReadFromYaml(
	ctx context.Context,
	path string,
) (*Configuration, error) {

	return nil, nil
}
