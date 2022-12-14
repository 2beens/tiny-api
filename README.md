# tiny-api
Tiny API for various experiments (k8s, argocd, ...)

Image is available in Docker Hub [stubin87/tiny-api](https://hub.docker.com/repository/docker/stubin87/tiny-api) repo.

Kubernetes 🚢 and ArgoCD 🐙 config files are located in the [2beens/tiny-api-k8s](https://github.com/2beens/tiny-api-k8s) repo.

Using [actions/build-and-push-docker-images](https://github.com/marketplace/actions/build-and-push-docker-images) to build and push new docker images to [my dockerhub](https://hub.docker.com/repository/docker/stubin87/tiny-api).
  - only on new releases!
  
Tiny API Web client in Vue can be found here: [2beens/tiny-api-web-client](https://github.com/2beens/tiny-api-web-client).

#### Tiny Stock Exchange (internal/tiny_stock_exchange_handler.go)
- communicates with [2beens/tiny-service](https://github.com/2beens/tiny-service) via gRPC
- uses protobuf from [2beens/tiny-stock-exchange-proto](https://github.com/2beens/tiny-stock-exchange-proto)
