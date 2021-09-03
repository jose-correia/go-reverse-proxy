// Package configuration contains the logic that is used to extract
// the proxy configuration from a .yaml file.
package configuration

import (
	"context"
	"go-reverse-proxy/app/values"
	"io/ioutil"

	"github.com/go-kit/kit/log"
	"gopkg.in/yaml.v2"
)

type Handler interface {
	// FromFile reads the proxy configuration .yaml file and converts it
	// to the internal Configuration structure.
	FromFile(ctx context.Context, filepath string) (*values.Configuration, error)
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

func (h *DefaultHandler) FromFile(
	ctx context.Context,
	filepath string,
) (*values.Configuration, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		h.logger.Log("module", "configurationHandler", "error", err)
		return nil, err
	}

	configuration, err := ParseYamlData(data)
	if err != nil {
		h.logger.Log("module", "configurationHandler", "error", err)
		return nil, err
	}

	return configuration, nil
}

func ParseYamlData(data []byte) (*values.Configuration, error) {
	yamlConfig := values.YamlConfig{}

	err := yaml.Unmarshal([]byte(data), &yamlConfig)
	if err != nil {
		return nil, err
	}

	return yamlConfig.ToConfiguration()
}
