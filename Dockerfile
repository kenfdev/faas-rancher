FROM golang:1.7.5

RUN mkdir -p /go/src/github.com/kenfdev/faas-rancher/

WORKDIR /go/src/github.com/kenfdev/faas-rancher

COPY vendor     vendor
COPY handlers	handlers
COPY types      types
COPY rancher     rancher
COPY server.go  .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o faas-rancher .

EXPOSE 8080
ENV http_proxy      ""
ENV https_proxy     ""

CMD ["./faas-rancher"]
