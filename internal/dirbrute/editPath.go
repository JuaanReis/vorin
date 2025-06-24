 package dirbrute

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"io"
	"time"
	"os"
	"vorin/pkg"
)

type Resultado struct {
	Status int
	URL    string
	Title string
	Size int
	Lines int
}

func Parser(endereco string, threads int, wordlist string, delay int) []Resultado {
	if delay == 0 {
		delay = 5
	}
	var resultados []Resultado
	var mu sync.Mutex
	var wg sync.WaitGroup

	sem := make(chan struct{}, threads)

	file, err := pkg.ReadLines(wordlist)
	if err != nil {
		fmt.Printf("[ERROR]: %v\n", err)
		os.Exit(1)
	}

	client := http.Client{
		Timeout: time.Duration(delay) * time.Second,
	}

	fakePath := "00"
	enderecoBase := strings.Replace(endereco, "Fuzz", "", -1)
	pathAle := strings.TrimRight(enderecoBase, "/") + "/" + fakePath
	respAle, error := http.Get(pathAle)
	if error != nil {
		fmt.Printf("[ERROR]: %s\n", error)
		os.Exit(1)
	}

	bodyAle, err := io.ReadAll(respAle.Body)
	if err != nil {
		fmt.Printf("[ERROR]: %v\n", err)
		os.Exit(1)
	}
	defer respAle.Body.Close()
	htmlAle := len(bodyAle)
	titleAle := getTitle(string(bodyAle))

	ini := time.Now()

	for _, path := range file {
		wg.Add(1)
		sem <- struct{}{}

		go func(p string) {
			defer wg.Done()
			defer func() { <-sem }()

			finalURL := strings.Replace(endereco, "Fuzz", path, -1)

			start := time.Now()
			resp, err := client.Get(finalURL)
			if err != nil {
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("[ERROR] %s\n", err)
			}
			htmlSize := len(body)
			title := getTitle(string(body))
			content := string(body)
			lines := len(strings.Split(content, "\n"))
			elapsed := time.Since(start)

			statusLabel, color := StatusColor(resp.StatusCode)

			if resp.StatusCode >= 200 && resp.StatusCode < 400 && !content404(string(body)) && (title != titleAle || htmlSize != htmlAle)  {
				fmt.Printf("%s[%-3d]%s  /%-20s (Size: %-6dB, Lines: %-3d) %-6s %s\n",
					color, resp.StatusCode, Reset,
					path,
					htmlSize,
					lines,
					elapsed.Truncate(time.Millisecond),
					statusLabel,
				)

				mu.Lock()
				resultados = append(resultados, Resultado{
					Status: resp.StatusCode,
					URL:    finalURL,
					Title: title,
					Size: htmlSize,
					Lines: lines,
				})
				mu.Unlock()
			}
		}(path)
	}
	wg.Wait()
	end := time.Since(ini)
	fmt.Printf("\n%s[✓]%s Scan completed in %s\n", Green, Reset, end.Truncate(time.Millisecond))
	return resultados
}