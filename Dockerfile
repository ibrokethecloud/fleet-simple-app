FROM golang:1.13 AS builder
RUN mkdir -p /src/github.com/ibrokethecloud/fleet-simple-demo
ARG VERSION
COPY . /src/github.com/ibrokethecloud/fleet-simple-demo
RUN cd /src/github.com/ibrokethecloud/fleet-simple-demo \
    && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.buildVersion=$VERSION" -o /root/fleet-simple-demo

## Using alpine
FROM alpine
COPY --from=builder /root/fleet-simple-demo /fleet-simple-demo
WORKDIR /
RUN touch /health
ENTRYPOINT ["/fleet-simple-demo"]
