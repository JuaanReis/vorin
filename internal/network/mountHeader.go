package network

import (
	"net/http"
)

func MountHeaders(req *http.Request, path string, stealth, bypass bool, custom map[string]string) {
	headers := map[string]string{}

	if stealth {
		for k, v := range GetRandomHeaders() {
			headers[k] = v
		}
	}

	for k, v := range custom {
		headers[k] = v
	}

	for k, v := range headers {
		req.Header.Set(k, v)
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Pragma", "no-cache")
	}
}
