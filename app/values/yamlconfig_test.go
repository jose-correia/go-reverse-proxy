package values_test

import (
	"go-reverse-proxy/app/values"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToConfiguration(t *testing.T) {

	yamlConfig := &values.YamlConfig{
		Proxy: values.ProxyYamlConfig{
			Listen: values.HostYamlConfig{
				Address: "127.0.0.1",
				Port:    5000,
			},
			Services: []values.ServiceYamlConfig{
				{
					Name:   "service",
					Domain: "service.com",
					Hosts: []values.HostYamlConfig{
						{
							Address: "127.0.0.2",
							Port:    5001,
						},
						{
							Address: "127.0.0.3",
							Port:    5002,
						},
					},
				},
			},
		},
	}

	configuration, err := yamlConfig.ToConfiguration()

	expectedConfig := &values.Configuration{
		Host: &values.Host{
			Address: "127.0.0.1",
			Port:    5000,
		},
		Services: map[string]*values.Service{
			"service.com": {
				Name:   "service",
				Domain: "service.com",
				Hosts: []*values.Host{
					{
						Address: "127.0.0.2",
						Port:    5001,
					},
					{
						Address: "127.0.0.3",
						Port:    5002,
					},
				},
				NextHostIndex: 0,
			},
		},
		MaxForwardRetries: 0,
	}

	assert.Equal(t, expectedConfig, configuration)
	assert.Nil(t, err)
}

func TestToConfigurationNoServices(t *testing.T) {

	yamlConfig := &values.YamlConfig{
		Proxy: values.ProxyYamlConfig{
			Listen: values.HostYamlConfig{
				Address: "127.0.0.1",
				Port:    5000,
			},
			Services: []values.ServiceYamlConfig{},
		},
	}

	configuration, err := yamlConfig.ToConfiguration()

	assert.NotNil(t, err)
	assert.Nil(t, configuration)
}

func TestToConfigurationNoHostAddress(t *testing.T) {

	yamlConfig := &values.YamlConfig{
		Proxy: values.ProxyYamlConfig{
			Listen: values.HostYamlConfig{
				Address: "",
				Port:    5000,
			},
			Services: []values.ServiceYamlConfig{
				{
					Name:   "service",
					Domain: "service.com",
					Hosts: []values.HostYamlConfig{
						{
							Address: "127.0.0.2",
							Port:    5001,
						},
						{
							Address: "127.0.0.3",
							Port:    5002,
						},
					},
				},
			},
		},
	}

	configuration, err := yamlConfig.ToConfiguration()

	assert.NotNil(t, err)
	assert.Nil(t, configuration)
}
