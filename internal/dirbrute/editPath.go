 package dirbrute

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"io"
	"time"
	"vorin/pkg"
)

type Resultado struct {
	Status int
	URL    string
}

func Parser(endereco string, threads int, wordlist string) []Resultado {
	var resultados []Resultado
	var mu sync.Mutex
	var wg sync.WaitGroup

	sem := make(chan struct{}, threads)

	file, err := ReadLines(wordlist)
	if err != nil {
		fmt.Println("[ERROR]:", err)
		return nil
	}

	u, err := url.Parse(endereco)
	if err != nil {
		fmt.Println("[ERROR]:", err)
		return nil
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	for _, path := range file {
		wg.Add(1)
		sem <- struct{}{}

		go func(p string) {
			defer wg.Done()
			defer func() { <-sem }()

			temp := *u
			temp.Path = p
			finalURL := temp.String()

			start := time.Now()
			resp, err := client.Get(finalURL)
			if err != nil {
				fmt.Printf("[ERROR GET] %s => %v\n", finalURL, err)
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("[ERROR] %s", err)
			}
			htmlSize := len(body)
			elapsed := time.Since(start)

			statusLabel, color := pkg.StatusColor(resp.StatusCode)

			if resp.StatusCode >= 200 && resp.StatusCode < 400 {
				fmt.Printf("%s[%d]%s   /%-30s %s (size: %dB)%s %-6s %-15s\n", color, resp.StatusCode, pkg.Reset, path, pkg.Blue, htmlSize, pkg.Reset, elapsed.Truncate(time.Millisecond), statusLabel)
				mu.Lock()
				resultados = append(resultados, Resultado{
					Status: resp.StatusCode,
					URL:    finalURL,
				})
				mu.Unlock()
			}
		}(path)
	}

	wg.Wait()
	return resultados
}