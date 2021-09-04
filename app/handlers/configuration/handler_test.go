// +build unit

package configuration_test

import (
	"go-reverse-proxy/app/values"
	"testing"

	config "go-reverse-proxy/app/handlers/configuration"

	"github.com/stretchr/testify/assert"
)

var (
	mockFiledata = []byte{
		112, 114, 111, 120, 121, 58, 10, 10, 32, 32, 108, 105, 115, 116, 101, 110, 58, 10, 32, 32, 32, 32, 97, 100, 100, 114, 101, 115, 115, 58, 32, 34, 49, 50, 55, 46, 48, 46, 48, 46, 49, 34, 10, 32, 32, 32, 32, 112, 111, 114, 116, 58, 32, 56, 48, 56, 48, 10, 10, 32, 32, 115, 101, 114, 118, 105, 99, 101, 115, 58, 10, 10, 32, 32, 32, 32, 45, 32, 110, 97, 109, 101, 58, 32, 109, 121, 45, 115, 101, 114, 118, 105, 99, 101, 10, 32, 32, 32, 32, 32, 32, 100, 111, 109, 97, 105, 110, 58, 32, 109, 121, 45, 115, 101, 114, 118, 105, 99, 101, 46, 109, 121, 45, 99, 111, 109, 112, 97, 110, 121, 46, 99, 111, 109, 10, 32, 32, 32, 32, 32, 32, 104, 111, 115, 116, 115, 58, 10, 32, 32, 32, 32, 32, 32, 32, 32, 45, 32, 97, 100, 100, 114, 101, 115, 115, 58, 32, 34, 49, 48, 46, 48, 46, 48, 46, 49, 34, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 112, 111, 114, 116, 58, 32, 57, 48, 57, 48, 10, 32, 32, 32, 32, 32, 32, 32, 32, 45, 32, 97, 100, 100, 114, 101, 115, 115, 58, 32, 34, 49, 48, 46, 48, 46, 48, 46, 50, 34, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 112, 111, 114, 116, 58, 32, 57, 48, 57, 48, 10,
	}
)

func TestParseYamlData(t *testing.T) {

	expectedConfiguration := &values.Configuration{
		Host: &values.Host{
			Address: "127.0.0.1",
			Port:    8080,
		},
		Services: map[string]*values.Service{
			"my-service.my-company.com": {
				Name:   "my-service",
				Domain: "my-service.my-company.com",
				Hosts: []*values.Host{
					{
						Address: "10.0.0.1",
						Port:    9090,
					},
					{
						Address: "10.0.0.2",
						Port:    9090,
					},
				},
			},
		},
	}

	config, _ := config.ParseYamlData(mockFiledata)

	assert.Equal(t, expectedConfiguration, config)
}

func TestInvalidYaml(t *testing.T) {
	config, err := config.ParseYamlData([]byte{})

	assert.Nil(t, config)
	assert.NotNil(t, err)
}
