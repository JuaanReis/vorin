package core

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"github.com/JuaanReis/vorin/internal/analyzer"
	"github.com/JuaanReis/vorin/internal/model"
	"github.com/JuaanReis/vorin/internal/modules"
	"github.com/JuaanReis/vorin/internal/network"
	"github.com/JuaanReis/vorin/internal/print"
	"github.com/JuaanReis/vorin/internal/collector"
	"github.com/JuaanReis/vorin/pkg"
)

func ParserGET(cfg model.ParserConfigGet) ([]model.Resultado, time.Duration) {
	var spinnerDone = make(chan bool)
	var resultados []model.Resultado
	var mu sync.Mutex
	var wg sync.WaitGroup
	var reader io.Reader
	var rateLimiter <-chan time.Time
	var resultadoUnico = make(map[string]bool)
	progressChan := make(chan struct{})
	var compiledRegexTitle *regexp.Regexp
	var compiledRegexBody *regexp.Regexp
	threadLimiter := make(chan struct{}, cfg.Threads)
	doneBar := make(chan struct{})

	if cfg.RateLimit > 0 {
		ticker := time.NewTicker(time.Second / time.Duration(cfg.RateLimit))
		rateLimiter = ticker.C
		defer ticker.Stop()
	}

	rawPaths, err := pkg.ReadLines(cfg.Wordlist)
	print.FatalIfErr(err)

	var file []string
	for _, path := range rawPaths {
		file = append(file, path)
		if strings.Contains(path, ".") {
			continue
		}

		for _, ext := range cfg.Extension {
			if ext != "" {
				file = append(file, path+ext)
			}
		}
	}

	if cfg.Shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(file), func(i, j int) {
			file[i], file[j] = file[j], file[i]
		})
	}

	var finalPaths []string
	for _, path := range file {
		finalPaths = append(finalPaths, path)
	}

	client := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
	}

	if !cfg.Redirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	modules.CreateClientProxy(cfg.Proxy, cfg.Timeout)

	var fakePath string
	if cfg.Compare != "" {
		fakePath = cfg.Compare
	} else {
		fakePath = "__vorin_this_should_not_exist_40028922__"
	}

	enderecoBase := strings.Replace(cfg.Endereco, "FUZZ", "", -1)
	pathAle := strings.TrimRight(enderecoBase, "/") + "/" + fakePath
	respAle, err := http.Get(pathAle)
	print.FatalIfErr(err)

	bodyAle, err := io.ReadAll(respAle.Body)
	print.FatalIfErr(err)

	defer respAle.Body.Close()
	fakeStructureSize, titleAle := collector.DataTargetFake(bodyAle)

	startTime := time.Now()

	if cfg.Silence {
		go print.Spinner("[Vorin] Running", spinnerDone)
	}

	if cfg.RegexTitle != "" {
		var err error
		compiledRegexTitle, err = regexp.Compile("(?i)" + cfg.RegexTitle)
		print.FatalIfErr(err)
	}

	if cfg.RegexBody != "" {
		var err error
		compiledRegexBody, err = regexp.Compile("(?i)" + cfg.RegexBody)
		print.FatalIfErr(err)
	}

	var (
		current   int32
		errors    int32
		reqPerSec int32
	)

	totalPaths := len(finalPaths)

	if !cfg.Silence {
		go print.UpdateProgressBar(totalPaths, &current, &errors, &reqPerSec, startTime, doneBar)
	}
		for _, p := range finalPaths {
			wg.Add(1)
			threadLimiter <- struct{}{}
			go func(p string) {
				defer wg.Done()
				defer func() { 
					if r := recover(); r != nil {
						atomic.AddInt32(&errors, 1)
					}
					<-threadLimiter
				}()
				atomic.AddInt32(&current, 1)
				atomic.AddInt32(&reqPerSec, 1)
				if rateLimiter != nil {
					<-rateLimiter
				}
			if cfg.Stealth || (cfg.MinDelay > 0 && cfg.MaxDelay > 0) {
				delayRange := cfg.MaxDelay - cfg.MinDelay
				delay := cfg.MinDelay + rand.Float64()*delayRange
				jitter := rand.Float64() * 0.1
				time.Sleep(time.Duration((delay + jitter) * float64(time.Second)))
			}
			finalURL := strings.Replace(cfg.Endereco, "FUZZ", p, -1)
			start := time.Now()
			req, err := http.NewRequest("GET", finalURL, nil)
			if err != nil {
				atomic.AddInt32(&errors, 1)
				return
			}
			network.MountHeaders(req, p, cfg.Stealth, cfg.Bypass, cfg.CustomHeaders)
			if cfg.RandomIp && (!cfg.Bypass || !cfg.Stealth) {
				ip := modules.RandomIP()
				req.Header.Set("X-Client-IP", ip)
				req.Header.Set("X-Forwarded-For", ip)
				req.Header.Set("CF-Connecting-IP", ip)
			}
			if cfg.RandomAgent && (!cfg.Bypass || !cfg.Stealth) {
				req.Header.Set("User-Agent", modules.RandomUserAgent())
			}
			body, resp, err := network.GetRequestWithRetry(req, client, reader, cfg.Retries)
			if err != nil || resp == nil {
				atomic.AddInt32(&errors, 1)
				return
			}
			lines, structureSize, title, htmlSize, content, text := collector.DataTaget(body)
			elapsed := time.Since(start)
			mustMatch := true
			if cfg.FilterBodyContent != "" && !strings.Contains(strings.ToLower(content), strings.ToLower(cfg.FilterBodyContent)) {
				mustMatch = false
			}
			if cfg.FilterTitleContent != "" && !strings.Contains(strings.ToLower(title), strings.ToLower(cfg.FilterTitleContent)) {
				mustMatch = false
			}
			matchRegexTitle := true
			matchRegexBody := true
			if compiledRegexTitle != nil && !compiledRegexTitle.MatchString(title) {
				matchRegexTitle = false
			}
			if compiledRegexBody != nil && !compiledRegexBody.MatchString(content) {
				matchRegexBody = false
			}
			tm := elapsed.Truncate(time.Millisecond)
			statusLabel, color := print.StatusColor(resp.StatusCode)
			if collector.IsSameContent(title, titleAle, cfg, htmlSize, lines, structureSize, fakeStructureSize, content) {
				return
			}
			
			key := strings.ToLower(p)
			if (len(cfg.FilterCode) == 0 || !cfg.FilterCode[resp.StatusCode]) && !analyzer.Content404(content) && lines > 0 && mustMatch && matchRegexTitle && matchRegexBody {
				if !resultadoUnico[key] {
					resultadoUnico[key] = true
					if cfg.Live {
						fmt.Print("\r\033[K")
						if cfg.StatusOnly {
							fmt.Printf("%s[%3d]%s %-26s\n",
								color, resp.StatusCode, print.Reset, p,
							)
						} else if cfg.Verbose {
							fmt.Print("\r\033[K")
							fmt.Printf("%s[%3d]%s  %4dw  %5dB  %4dL  %6s  %s\n",
								color, resp.StatusCode, print.Reset,
								text,
								htmlSize,
								lines,
								elapsed.Truncate(time.Millisecond),
								statusLabel,
							)
							fmt.Printf(" ├─ URL     : %s\n", finalURL)
							fmt.Printf(" ├─ FUZZ    : %s\n", p)
							if title != "" {
								fmt.Printf(" └─ Title   : %s\n\n", title)
							}
						} else {
							fmt.Printf("%s[%3d]%s  %-26s Words: %-6d Size: %-6dB Lines: %-5d %-6s %-11s\n",
								color, resp.StatusCode, print.Reset,
								p,
								text,
								htmlSize,
								lines,
								elapsed.Truncate(time.Millisecond),
								statusLabel,
							)
						}
					}
					mu.Lock()
					resultados = append(resultados, model.Resultado{
						Label:  statusLabel,
						Status: resp.StatusCode,
						URL:    p,
						Title:  title,
						Text:   text,
						Size:   htmlSize,
						Lines:  lines,
						Time:   tm,
						Color:  color,
						Endereco: finalURL,
					})
					mu.Unlock()
				}
			}
		}(p)
	}

	wg.Wait()
	fmt.Print("\r\033[K")
	fmt.Print("\033[1A\r\033[K")
	close(progressChan)

	if !cfg.Silence {
		<-doneBar
	}

	fmt.Println()

	if cfg.Silence {
		spinnerDone <- true
	}

	end := time.Since(startTime)
	return resultados, end
}
