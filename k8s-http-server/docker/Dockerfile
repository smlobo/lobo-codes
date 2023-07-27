FROM golang:alpine AS build

WORKDIR /app
COPY . .

RUN go build -o k8s-http-server cmd/k8s-http-server.go
RUN echo ${GOLANG_VERSION} > golang_version.txt

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/k8s-http-server .
COPY --from=build /app/golang_version.txt .

# Copy each subdomain, domain, and common directories
COPY IP2LOCATION-LITE-DB11.IPV6.BIN .
COPY ./common common
COPY ./domain domain
COPY ./amelia amelia
COPY ./ryan ryan
COPY ./sheldon sheldon
COPY ./test-vue/dist test-vue/dist
COPY ./config config
COPY ./ssl-certificates ssl-certificates

CMD ["/app/k8s-http-server"]

