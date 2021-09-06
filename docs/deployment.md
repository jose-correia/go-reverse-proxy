<h1 align="center">Deployment</h4>

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

