# Golang Reverse Proxy


# Deploying to minikube

- Set the cluster docker environment as the current docker environment
`eval $(minikube docker-env)`

- Build the image:
`make docker-build`

- Install Helm Chart:
`make helm-install`


