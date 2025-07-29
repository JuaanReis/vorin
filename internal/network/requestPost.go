package network

import (
	"bytes"
	"io"
	"net/http"
	"time"
	"regexp"
)

var jsonRegex = regexp.MustCompile(`^\s*(\{.*\}|\[.*\])\s*$`)

func PostRequestWithRetry(req *http.Request, client *http.Client, payload []byte, retries int) ([]byte, *http.Response, error) {
	var (
		body []byte
		resp *http.Response
		err  error
	)

	isJSON := jsonRegex.Match(payload)
	if isJSON {
		req.Header.Set("Content-Type", "application/json")
	}

	for i := 0; i <= retries; i++ {
		req.Body = io.NopCloser(bytes.NewReader(payload))

		resp, err = client.Do(req)
    	if err == nil && resp != nil && resp.StatusCode < 500 {
        	body, err = io.ReadAll(resp.Body)
        	resp.Body.Close()
        		if err != nil {
            	return nil, resp, err
        }
        return body, resp, nil
    }

		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}

		if i < retries {
			time.Sleep(500 * time.Millisecond)
		}
	}

	return nil, resp, err
}
