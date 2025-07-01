package dirbrute

import (
	"net/url"
	"net/http"
	"time"
	"os"
	"fmt"
)

func CreateClientProxy(proxy string, timeout int) *http.Client {
	if proxy == "" {
		return &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	}

	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		fmt.Printf("[ERROR] Invalid proxy: %v\n", err)
		os.Exit(1)
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	}

	// Testa o proxy
	testClient := &http.Client{Transport: transport, Timeout: 5 * time.Second}
	testReq, _ := http.NewRequest("GET", "https://www.google.com/", nil)
	resp, err := testClient.Do(testReq)
	if err != nil {
		fmt.Printf("[ERROR] Proxy did not respond: %v\n", err)
		os.Exit(1)
	}
	resp.Body.Close()

	return &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}
