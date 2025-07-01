package dirbrute

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"io"
	"time"
	"os"
	"github.com/JuaanReis/vorin/pkg"
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
	Time time.Duration
	Label string
	Color string
}

var spinnerDone = make(chan bool)

func Parser(endereco string, threads int, wordlist string, minDelay int, maxDelay int, timeout int, customHeaders map[string]string, code map[int]bool, stealth bool, proxy string, silence bool, live bool) ([]Resultado, time.Duration) {
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

	CreateClientProxy(proxy, timeout)

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
	stringBodyAle := string(bodyAle)
	structureOnly := cleanStructure(stringBodyAle)
	fakeStructureSize := len(structureOnly)
	titleAle := strings.TrimSpace(strings.ToLower(getTitle(stringBodyAle)))

	bar := progressbar.NewOptions(len(file),
		progressbar.OptionSetDescription("Testing paths..."),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(30),
		progressbar.OptionClearOnFinish(),
	)

	ini := time.Now()

	if silence {
		go Spinner("[Vorin] Running", spinnerDone)
	}


	for _, path := range file {
		wg.Add(1)
		sem <- struct{}{}

		go func(p string) {
			defer wg.Done()
			defer func() { <-sem }()
			if !silence {
				defer bar.Add(1)
			}
			if stealth || (minDelay > 0 && maxDelay > 0) {
				delay := rand.Intn(maxDelay - minDelay + 1) + minDelay
				time.Sleep(time.Duration(delay) * time.Second)
			}

			finalURL := strings.Replace(endereco, "Fuzz", p, -1)

			start := time.Now()

			req, err := http.NewRequest("GET", finalURL, nil)
			if err != nil {
				fmt.Printf("[ERROR]: %s", err)
				return
			}

			if !stealth {
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
			}

			if stealth {
				HeadersC := GetRandomHeaders()
				for key, value := range HeadersC {
					req.Header.Set(key, value)
				}
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

			lines, structureSize, title, htmlSize, content := DataTaget(body)

			elapsed := time.Since(start)

			isSameContent := title == titleAle || structureSize == fakeStructureSize
			tm := elapsed.Truncate(time.Millisecond)
			statusLabel, color := StatusColor(resp.StatusCode)

			if code[resp.StatusCode] &&
				!content404(content) &&
				!isSameContent  {
					if live {
						bar.Clear()
						fmt.Printf("%s[%-3d]%s  /%-30s  Size: %-6dB  Lines: %-3d  %-6s  %s\n",
						color, resp.StatusCode, Reset,
						p,
						htmlSize,
						lines,
						elapsed.Truncate(time.Millisecond),
						statusLabel,
						)
					}
				mu.Lock()
				resultados = append(resultados, Resultado{
					Label: statusLabel,
					Status: resp.StatusCode,
					URL:    p,
					Title:  title,
					Size:   htmlSize,
					Lines:  lines,
					Time: tm,
					Color: color,
				})
				mu.Unlock()
			}
		}(path)
	}

	wg.Wait()
	if silence {
		spinnerDone <- true
	}

	bar.Clear()
	end := time.Since(ini)

	return resultados, end
}
