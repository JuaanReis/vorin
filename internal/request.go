package internal

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"time"
)

func GetRequest(req *http.Request, client *http.Client, reader io.Reader) ([]byte, *http.Response, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode == 429 {
		fmt.Println("[INFO] Received 429. Waiting 30s...")
		time.Sleep(30 * time.Second)
	}

	var bodyReader io.Reader
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			resp.Body.Close()
			return nil, nil, err
		}
		defer gzipReader.Close()
		bodyReader = gzipReader
	} else {
		bodyReader = resp.Body
	}

	body, err := io.ReadAll(bodyReader)
	defer resp.Body.Close()
	if err != nil {
		return nil, resp, err
	}

	return body, resp, nil
}

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
