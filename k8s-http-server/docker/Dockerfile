FROM golang:alpine AS build

WORKDIR /app
COPY . .

RUN go build -o k8s-http-server cmd/k8s-http-server.go

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/k8s-http-server .
# Copy each subdomain, domain, and common directories
COPY IP2LOCATION-LITE-DB11.IPV6.BIN .
COPY ./common common
COPY ./domain domain
COPY ./amelia amelia
COPY ./ryan ryan
COPY ./sheldon sheldon
COPY ./test-vue/dist test-vue/dist

CMD ["/app/k8s-http-server"]

