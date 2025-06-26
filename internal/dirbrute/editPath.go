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
	"github.com/schollz/progressbar/v3"
)

type Resultado struct {
	Status int
	URL    string
	Title  string
	Size   int
	Lines  int
}

func Parser(endereco string, threads int, wordlist string, delay int) []Resultado {
	if delay == 0 {
		delay = 5
	}
	var resultados []Resultado
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

	fakePath := "01"
	enderecoBase := strings.Replace(endereco, "Fuzz", "", -1)
	pathAle := strings.TrimRight(enderecoBase, "/") + "/" + fakePath
	respAle, err := http.Get(pathAle)
	if err != nil {
		fmt.Printf("[ERROR]: %s\n", err)
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

	bar := progressbar.NewOptions(len(file),
		progressbar.OptionSetDescription("Testing paths..."),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(30),
		progressbar.OptionClearOnFinish(),
	)

	ini := time.Now()

	for _, path := range file {
		wg.Add(1)
		sem <- struct{}{}

		go func(p string) {
			defer wg.Done()
			defer func() { <-sem }()

			finalURL := strings.Replace(endereco, "Fuzz", p, -1)

			start := time.Now()

			req, err := http.NewRequest("GET", finalURL, nil)
			if err != nil {
				fmt.Printf("[ERROR]: %s", err)
				os.Exit(1)
			}

			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
			req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
			req.Header.Set("Accept-Language", "en-US,en;q=0.5")
			req.Header.Set("Accept-Encoding", "gzip, deflate")
			req.Header.Set("Connection", "keep-alive")
			req.Header.Set("Upgrade-Insecure-Requests", "1")
			req.Header.Set("Cache-Control", "max-age=0")

			resp, err := client.Do(req)
			if err != nil {
				bar.Add(1)
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)

			if err != nil {
				fmt.Printf("[ERROR] %s\n", err)
				bar.Add(1)
				return
			}

			htmlSize := len(body)
			title := getTitle(string(body))
			content := string(body)
			lines := len(strings.Split(content, "\n"))
			elapsed := time.Since(start)

			statusLabel, color := StatusColor(resp.StatusCode)

			if resp.StatusCode >= 200 && resp.StatusCode < 400 && !content404(string(body)) && htmlSize != 0 && (title != titleAle || htmlSize != htmlAle)  {

				bar.Clear()

				fmt.Printf("%s[%-3d]%s  /%-20s (Size: %-6dB, Lines: %-3d) %-6s %s\n",
					color, resp.StatusCode, Reset,
					p,
					htmlSize,
					lines,
					elapsed.Truncate(time.Millisecond),
					statusLabel,
				)
				resultados = append(resultados, Resultado {
        	Status: resp.StatusCode,
        	URL:    finalURL,
        	Title:  title,
        	Size:   htmlSize,
        	Lines:  lines,
    		})
			}

			bar.Add(1)
		}(path)
	}

	wg.Wait()
	bar.Clear()

	end := time.Since(ini)
	fmt.Printf("\n%s[✓]%s Scan completed in %s\n", Green, Reset, end.Truncate(time.Millisecond))
	return resultados
}
