FROM golang:alpine AS builder

WORKDIR $GOPATH/src/github.com/marcusnoble/kube-event-logger

RUN apk update && \
    apk add git build-base

Add . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --installsuffix cgo --ldflags="-s" -o /kubelogger


FROM scratch
COPY --from=builder /kubelogger /bin/kubelogger
ENTRYPOINT ["/bin/kubelogger"]
