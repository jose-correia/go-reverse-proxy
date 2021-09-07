<h1 align="center">Deployment</h4>
<p align="center">
  <a href="#initial-configuration">Initial Configuration</a> •
  <a href="#local-deployment">Local Deployment</a> •
  <a href="#deploying-to-kubernetes-as-a-helm-chart">Deploying to Kubernetes as a Helm Chart</a>
</p>



## Initial Configuration

1. Create a .env file with the environment variables configuration (the values below represent the default values if no file is created):

```yaml
CONFIGURATION_FILENAME: "proxyConfig.yaml"
MAX_HTTP_RETRIES: 2
MAX_FORWARD_RETRIES: 2
HTTP_CACHE_TTL_SECONDS: 60
METRICS_ADDR: ":8090"
```

2. Add your own service routes to the proxy configuration file which can be found in ```proxy-configs/```:

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



## Local Deployment

There are two ways for deploying the system locally:

### Simple execution

```sh
make run
```

### Using docker-compose

```sh
docker-compose run reverse-proxy
```





## Deploying to Kubernetes as a Helm Chart

The service can be easily plugged in a Kubernetes cluster by intalling the respective Helm Chart. The deployment is configured to consume the following resources:

```yaml
resource:
	limits:
		cpu: 600m # maximum CPU that the pod is allowed to request
		memory: 512Mi # maximum memory hat the pod is allowed to request
	requests:
		cpu: 100m # CPU initially allocated to the pod
		memory: 128Mi # Memory initially allocated to the pod
```

1. **For a local Kubernetes environment:** install Minikube

```console
brew install minikube
```
```sh
minikube start
```

2. Build the Docker image

```sh
eval $(minikube docker-env)
```

```sh
make docker-build
```

3. Install Helm Chart in the Kubernetes cluster

```sh
make helm-install
```

Finally, you can see your deployment  and the respective pod by executing:

```shell
kubectl get deployment | grep reverse-proxy
```

```shell
kubectl get pod | grep reverse-proxy
```

