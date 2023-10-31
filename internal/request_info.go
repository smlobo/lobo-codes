package internal

import (
	"net/http"

	"go.opentelemetry.io/otel"
)

func requestInfo(request *http.Request, tableName string) {
	// Child span for logging request info
	_, span := otel.Tracer("k8s-http-server").Start(request.Context(), "request-info-logging")
	defer span.End()

	var info RequestInfo

	info.UserAgent = request.Header.Get("User-Agent")
	info.OrigRemoteAddress = request.RemoteAddr
	info.RemoteAddress = request.Header.Get("X-Forwarded-For")

	// Once the information is extracted from the request, the remainder of processing can be
	// done concurrently
	go cassandraLogRequest(&info, tableName, request)
}
