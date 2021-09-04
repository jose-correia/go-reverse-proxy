package encoder

import (
	"context"
	"encoding/json"
	"net/http"

	httpkit "github.com/go-kit/kit/transport/http"
)

type message struct {
	Message string `json:"error"`
}
type Error struct {
	Message string
	Code    int
}

func (e *Error) GetCode() int {
	return e.Code
}

func (e *Error) Error() string {
	return e.Message
}

type HTTPError interface{ GetCode() int }

func Encode(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(asHTTPCode(err))
	_ = json.NewEncoder(w).Encode(message{Message: err.Error()})
}

func asHTTPCode(err error) int {
	switch err.(type) {
	case HTTPError:
		return err.(HTTPError).GetCode()
	default:
		return http.StatusInternalServerError
	}
}

func ErrorOption() httpkit.ServerOption {
	return httpkit.ServerErrorEncoder(Encode)
}
