package network

import (
	"io"
	"net/http"
	"time"
)

func GetRequestWithRetry(req *http.Request, client *http.Client, reader io.Reader, retries int) ([]byte, *http.Response, error) {
	var (
		body []byte
		resp *http.Response
		err  error
	)

	for i := 0; i <= retries; i++ {
		body, resp, err = GetRequest(req, client, reader)

		if err == nil && resp != nil && resp.StatusCode < 500 {
			return body, resp, nil
		}

		if i < retries {
			time.Sleep(500 * time.Millisecond)
		}
	}

	return body, resp, err
}
