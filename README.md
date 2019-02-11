# kube-event-logger

Capture Kubernetes events and log them out to stdout so they can be picked up by a central log collector.

## Usage

### Out-of-cluster

```
go build main.go -o kube-event-logger
KUBECONFIG=~/.kube/config ./kube-event-logger
```

### In-cluster

Apply the provided kubernetes manifest:
```
kubectl apply -f ./kubernetes.yaml
```

## Config

Take a look at [.kube-event-logger.yaml](.kube-event-logger.yaml) for the possible configuration. If no configuration is found, or configuration fails to load the application will only report on pods across all namespaces by default.
