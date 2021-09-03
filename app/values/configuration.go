package values

// Type Configuration is used to represent the configuration
// of the reverse proxy service.
type Configuration struct {
	Host     *Host      // the host configuration of the reverse proxy
	Services []*Service // list of supported downstream services
}

// Type Service is used to represent a downsteam service
// that the reverse proxy can forward too
type Service struct {
	Name   string  // name of the service
	Domain string  // domain of the service
	Hosts  []*Host // host list containing the service instances
}

type Host struct {
	Address string // IPv4 address
	Port    int32  // Port that is listening
}
