package values

import "fmt"

// Type YamlConfig the structure where the proxy
// configuration .yaml will be parsed into
type YamlConfig struct {
	Proxy ProxyYamlConfig
}

// ToConfiguration parses a yamlConfiguration into the internal
// Configurations representation of the reverse proxy service
func (y *YamlConfig) ToConfiguration() (*Configuration, error) {
	services := make(map[string]*Service)

	for _, service := range y.Proxy.Services {
		var hosts []*Host
		for _, host := range service.Hosts {
			hosts = append(hosts, &Host{
				Address: host.Address,
				Port:    host.Port,
			})
		}

		services[service.Domain] = &Service{
			Name:   service.Name,
			Domain: service.Domain,
			Hosts:  hosts,
		}
	}

	hostAddress := y.Proxy.Listen.Address
	hostPort := y.Proxy.Listen.Port

	if len(services) < 1 || hostAddress == "" || hostPort == 0 {
		return nil, fmt.Errorf("the .yaml configuration is invalid")
	}

	return &Configuration{
		Host: &Host{
			Address: y.Proxy.Listen.Address,
			Port:    y.Proxy.Listen.Port,
		},
		Services: services,
	}, nil
}

type ProxyYamlConfig struct {
	Listen   HostYamlConfig
	Services []ServiceYamlConfig `yaml:",flow"`
}

type ServiceYamlConfig struct {
	Name   string
	Domain string
	Hosts  []HostYamlConfig `yaml:",flow"`
}
type HostYamlConfig struct {
	Address string
	Port    int32
}
