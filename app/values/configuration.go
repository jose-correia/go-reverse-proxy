package values

import "fmt"

// Type Configuration is used to represent the configuration
// of the reverse proxy service.
type Configuration struct {
	Host     *Host               // the host configuration of the reverse proxy
	Services map[string]*Service // map of supported downstream services

	// list of status codes that should result in a redirect of the request
	// to another instance
	RetryableStatusCodes []int
	// maximum number of retries to different instances
	MaxForwardRetries int
}

// GetServiceByDomain returns a pointer to a Service given a domain string
func (c *Configuration) GetServiceByDomain(domain string) *Service {
	service, _ := c.Services[domain]
	return service
}

// Type Service is used to represent a downsteam service
// that the reverse proxy can forward too
type Service struct {
	Name   string  // name of the service
	Domain string  // domain of the service
	Hosts  []*Host // host list containing the service instances

	// used to store the index of the next host to use
	// this value is incrementally updated after each request
	// in order to apply a Round-Robin load balancing algorithm
	NextHostIndex int32
}

// GetNextHost retrieves a pointer to the next Host to be used
func (s *Service) GetNextHost() *Host {
	return s.Hosts[s.NextHostIndex]
}

type Host struct {
	Address string // IPv4 address
	Port    int32  // Port that is listening
}

// ToURL creates the URL representation composed of a Host address and port
func (h *Host) ToURL() string {
	return fmt.Sprintf("%s:%d", h.Address, h.Port)
}
