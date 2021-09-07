package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	cacheControlMaxAgeArg     = "--cache-control-max-age"
	defaultCacheControlMaxAge = 60
	portArg                   = "--port"
	defaultPort               = ":8000"
	idArg                     = "--id"
)

type Server struct {
	id                 string
	useCache           bool
	cacheControlMaxAge int
}

func (s *Server) handleRequest(w http.ResponseWriter, req *http.Request) {
	log.Println("Server", s.id, "received request")

	if s.useCache {
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", s.cacheControlMaxAge))
	}

	fmt.Fprintf(w, "Hello World!\n")
}

func main() {
	var err error
	var useCache bool
	cacheControlMaxAge := defaultCacheControlMaxAge
	port := defaultPort
	var id string

	for idx, value := range os.Args {
		if idx == len(os.Args)-1 {
			break
		}

		nextArg := os.Args[idx+1]

		if value == cacheControlMaxAgeArg {
			useCache = true
			if idx < len(os.Args)-1 {
				cacheControlMaxAge, err = strconv.Atoi(nextArg)
				if err != nil {
					log.Println("Invalid cache max age argument")
					useCache = false
				}
			}
		}

		if value == portArg {
			port = nextArg
		}

		if value == idArg {
			id = nextArg
		}
	}

	server := &Server{
		id:                 id,
		useCache:           useCache,
		cacheControlMaxAge: cacheControlMaxAge,
	}

	log.Println("Server", server.id, "is listening in port", port, "| useCache:", useCache)

	http.HandleFunc("/", server.handleRequest)
	http.ListenAndServe(port, nil)
}
