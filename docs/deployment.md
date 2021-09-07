<h1 align="center">Deployment</h4>

- [Configuration](#configuration)

- [Simple local deployment](#simple-local-deployment)

- [Deployment using docker-compose](#deployment-using-docker-compose)

- [Deployment to Minikube as a Helm Chart](#deployment-to-minikube-as-a-helm-chart)
  - [1. Instaling Minikube](#1-instaling-minikube)
  - [2. Building the Docker image](#2-building-the-docker-image)
  - [3. Installing Helm Chart:](#3-installing-helm-chart)
  
  

## Configuration

```yaml
CONFIGURATION_FILENAME: "proxyConfig.yaml"
MAX_HTTP_RETRIES: 2
MAX_FORWARD_RETRIES: 2
HTTP_CACHE_TTL_SECONDS: 60
METRICS_ADDR: ":8090"
```



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




## Simple local deployment
```sh
make run
```

---

## Deployment using docker-compose
```sh
docker-compose run reverse-proxy
```

**Note:** 
- Since Docker for MacOS is actually running inside a Linux VM, you won't be able to access the container even though the service is using the `network_mode: host`. With this being said, the [Simple local deployment](#simple-local-deployment) should be a better approach in this situation.

---

## Deployment to Minikube as a Helm Chart

### 1. Instaling Minikube
```console
brew install minikube
```
```sh
minikube start
```

### 2. Building the Docker image
```sh
eval $(minikube docker-env)
```

```sh
make docker-build
```

### 3. Installing Helm Chart:
```sh
make helm-install
```

