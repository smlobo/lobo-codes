FROM golang:alpine

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN go build -o cassandra-client cmd/cassandra-client.go

CMD ["/app/cassandra-client"]
