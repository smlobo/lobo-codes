
.PHONY: static
static:
	go run cmd/serve-static.go

.PHONY: template
template:
	go run cmd/serve-template.go

static2:
	go run cmd/serve-static2.go :3000 static ../ip2location-db/IP2LOCATION-LITE-DB11.IPV6.BIN
