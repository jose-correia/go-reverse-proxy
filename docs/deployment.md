# Deployment

- [Deployment](#deployment)
  - [Simple local deployment](#simple-local-deployment)
  - [Deployment using docker-compose](#deployment-using-docker-compose)
  - [Deployment to Minikube as a Helm Chart](#deployment-to-minikube-as-a-helm-chart)
    - [1. Instaling Minikube](#1-instaling-minikube)
    - [2. Building the Docker image](#2-building-the-docker-image)
    - [3. Installing Helm Chart:](#3-installing-helm-chart)


The service is listening in the port `8080`.


## Simple local deployment
```sh
make run
```

After compiling the binary file and executing the service, the output will be:
```sh
go run cmd/proxy/main.go
timestamp=2021-09-04T18:34:15.489618Z service=reverse-proxy address=127.0.0.1:8080 status=listening...
```

---

## Deployment using docker-compose
```sh
docker-compose run reverse-proxy
```

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

