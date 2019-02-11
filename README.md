# kube-event-logger

Capture Kubernetes events and log them out to stdout so they can be picked up by a central log collector.

## TODO

- [ ] Support more resource types (currently only pods)
- [ ] Allow for configuration
- [x] Add out of cluster support

## Usage

### Out-of-cluster

```
go build main.go -o kube-event-logger
KUBECONFIG=~/.kube/config ./kube-event-logger
```

### In-cluster

Apply the provided kubernetes manifest:
```
kubectl apply -f ./kubernetes.yaml
```
