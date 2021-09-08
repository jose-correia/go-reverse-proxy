package values

import "net/http"

// Request is used to represent the client request
type Request struct {
	Method     string      // HTTP method
	Endpoint   string      // Endpoint of the downstream service that is being requested
	Header     http.Header // Request headers
	HostHeader string      // Host header
	Parameters string      // URL query parameters
	Payload    []byte      // Request payload data
}
