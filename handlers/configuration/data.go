package configuration

type Configuration struct {
	Host     *Host     // the host configuration of the reverse proxy
	Services []Service // list of supported downstream services
}

type Service struct {
	Name   string  // name of the service
	Domain string  // domain of the service
	Hosts  []*Host // host configuration of the service
}

type Host struct {
	Address string // IPv4 address
	Port    int32  // Port that is listening
}
