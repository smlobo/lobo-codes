FROM golang:alpine

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN go build -o k8s-http-server cmd/k8s-http-server.go

CMD ["/app/k8s-http-server"]
