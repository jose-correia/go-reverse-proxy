# Reverse Proxy


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

## Deployment

Steps to deploy the service can be found on the [Deployment](docs/deployment.md) doc.

