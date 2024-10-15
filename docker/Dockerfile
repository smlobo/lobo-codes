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
COPY ./ssl-certificates ssl-certificates

COPY ./config config

COPY ./common common
COPY ./domain domain
COPY ./amelia amelia
COPY ./ryan ryan
COPY ./bliu bliu
COPY ./sheldon sheldon

COPY ./test-vue/dist test-vue/dist
COPY ./test-vue/static test-vue/static

COPY ./wasm/index.html wasm/index.html
COPY ./wasm/static wasm/static
COPY ./wasm/wasm-rotating-cube/dist wasm/wasm-rotating-cube/dist
COPY ./wasm/h-tree/dist wasm/h-tree/dist
COPY ./wasm/fractal-circle/dist wasm/fractal-circle/dist
COPY ./wasm/julia-set/dist wasm/julia-set/dist
COPY ./wasm/collision-system/dist wasm/collision-system/dist

CMD ["/app/k8s-http-server"]

