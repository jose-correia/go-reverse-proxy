package values

import "net/http"

type Request struct {
	Method     string
	Header     http.Header
	HostHeader string
	Parameters string
	Payload    []byte
}
