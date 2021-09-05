<h1 align="center">
  <br>
  <img width="500" src="docs/logo.png" alt="ArminC AutoExec">
</h1>

<h4 align="center">A Golang scalable reverse proxy</h4>

<p align="center">
  <a href="#about">About</a> •
  <a href="#architecture-design">Architecture Design</a> •
  <a href="#slis">SLIs</a>
  <a href="#dependencies">Dependencies</a> •
  <a href="#configuration">Configuration</a> •
  <a href="#test">Test</a> •
  <a href="#deployment">Deployment</a> •
</p>

---

## About

The `go-reverse-proxy` service is a Golang reverse proxy implementation, completely configured to be deployed in a Kubernetes environment and scale as needed.

## Architecture Design
The system architecture design can be viewed in the [Architecture](docs/architecture.md) doc.

## SLIs
### Health
- `Uptime`
- `Request volume`
- `Success rate`: percentage of non 4xx-5xx status code responses;

### Resource usage
- `Number of pods`
- `Network I/0 Usage`
- `% CPU usage per pod`
- `% Memory usage per pod`
  
### Performance

- `Request Processing Time`: time elapsed since the client request is read by the proxy, until it is forwarded to the downstream service;
- `Response Processing Time`: time elapsed since the downstream service responds to the proxy, until it is forwarded to the client;
- `Latency`: time elapsed since the proxy receives the client request, until it responds back to him.

## Dependencies
- [go-kit](https://github.com/go-kit/kit)
- [gorilla/mux](https://github.com/gorilla/mux)
- [go-retryablehttp](https://github.com/hashicorp/go-retryablehttp)
- [go-yaml](https://github.com/go-yaml/yaml)
- [httpcache](https://github.com/bxcodec/httpcache)

## Configuration

The template configuration file can be found in [Proxy Config](proxy-configs/proxyConfig.yaml)

```yaml
proxy:

  listen:
    address: "127.0.0.1"
    port: 8080

  services:

    - name: my-service
      domain: my-service.my-company.com
      hosts:
        - address: "10.0.0.1"
          port: 9090
        - address: "10.0.0.2"
          port: 9090
```

## Test
```sh
make test
```

## Deployment

Steps to deploy the service can be found on the [Deployment](docs/deployment.md) doc.


## Usage
