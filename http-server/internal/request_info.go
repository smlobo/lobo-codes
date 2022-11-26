package internal

import (
	"net/http"
	"strings"
)

func requestInfo(request *http.Request, info *RequestInfo) {
	info.UserAgent = request.Header.Get("User-Agent")
	info.RemoteAddress = strings.Trim(strings.Split(request.RemoteAddr, ":")[0], "[]")
	geoInfo(info)
}
