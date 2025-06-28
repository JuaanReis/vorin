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
	"math/rand"
	"compress/gzip"
)

type Resultado struct {
	Status int
	URL    string
	Title  string
	Size   int
	Lines  int
}

func Parser(endereco string, threads int, wordlist string, minDelay int, maxDelay int, timeout int, customHeaders map[string]string) []Resultado {
	var resultados []Resultado
	var mu sync.Mutex
	var wg sync.WaitGroup
	var reader io.Reader

	sem := make(chan struct{}, threads)


	file, err := pkg.ReadLines(wordlist)
	if err != nil {
		fmt.Printf("[ERROR]: %v\n", err)
		os.Exit(1)
	}

	client := &http.Client{
    Timeout: time.Duration(timeout) * time.Second,
    CheckRedirect: func(req *http.Request, via []*http.Request) error {
        return http.ErrUseLastResponse
    },
	}

	fakePath := "__vorin_this_should_not_exist_473827382__"
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
	structureOnly := cleanStructure(string(bodyAle))
	fakeStructureSize := len(structureOnly)
	titleAle := strings.TrimSpace(strings.ToLower(getTitle(string(bodyAle))))

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
			defer bar.Add(1)

			delay := rand.Intn(maxDelay - minDelay + 1) + minDelay
			time.Sleep(time.Duration(delay) * time.Second)

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

			for key, val := range customHeaders {
				req.Header.Set(key, val)
			}

			resp, err := client.Do(req)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			if resp.Header.Get("Content-Encoding") == "gzip" {
				gzipReader, err := gzip.NewReader(resp.Body)
			if err != nil {
				fmt.Printf("[ERROR]: gzip decode error: %v\n", err)
				return
			}
			defer gzipReader.Close()

			reader = gzipReader
			} else {
				reader = resp.Body
			}

			body, err := io.ReadAll(reader)
			if err != nil {
				fmt.Printf("[ERROR] %s\n", err)
				return
			}

			structureOnly := cleanStructure(string(body))
			structureSize := len(structureOnly)
			htmlSize := len(body)
			title := strings.TrimSpace(strings.ToLower(getTitle(string(body))))
			content := string(body)
			lines := len(strings.Split(content, "\n"))
			elapsed := time.Since(start)

			statusLabel, color := StatusColor(resp.StatusCode)
			isSameContent := title == titleAle || structureSize == fakeStructureSize

			if (resp.StatusCode >= 200 && resp.StatusCode < 399) &&
				!content404(content) &&
				!isSameContent  {
					bar.Clear()
					fmt.Printf("%s[%-3d]%s  /%-20s (Size: %-6dB, Lines: %-3d) %-6s %s\n",
					color, resp.StatusCode, Reset,
					p,
					htmlSize,
					lines,
					elapsed.Truncate(time.Millisecond),
					statusLabel,
				)
				mu.Lock()
				resultados = append(resultados, Resultado{
					Status: resp.StatusCode,
					URL:    finalURL,
					Title:  title,
					Size:   htmlSize,
					Lines:  lines,
				})
				mu.Unlock()
			}
		}(path)
	}

	wg.Wait()
	bar.Clear()

	end := time.Since(ini)
	fmt.Printf("\n%s[✓]%s Scan completed in %s\n", Green, Reset, end.Truncate(time.Millisecond))
	return resultados
}
