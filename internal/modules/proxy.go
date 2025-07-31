package modules

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"golang.org/x/net/proxy"
	"github.com/JuaanReis/vorin/internal/print"
)

func CreateClientProxy(proxyStr string, timeout int) *http.Client {
	var transport *http.Transport

	if proxyStr == "" {
		transport = &http.Transport{}
	} else if strings.HasPrefix(proxyStr, "socks5://") {
		addr := strings.TrimPrefix(proxyStr, "socks5://")
		dialer, err := proxy.SOCKS5("tcp", addr, nil, proxy.Direct)
		if err != nil {
			fmt.Printf("[ERROR] SOCKS5 dialer failed: %v\n", err)
			os.Exit(1)
		}

		transport = &http.Transport{
			Dial: dialer.Dial,
		}
	} else {
		proxyURL, err := url.Parse(proxyStr)
		if err != nil {
			fmt.Printf("[ERROR] Invalid proxy URL: %v\n", err)
			os.Exit(1)
		}
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	testClient := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	testReq, _ := http.NewRequest("GET", "http://check.torproject.org", nil)
	resp, err := testClient.Do(testReq)
	print.FatalIfErr(err)
	resp.Body.Close()

	return &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}
