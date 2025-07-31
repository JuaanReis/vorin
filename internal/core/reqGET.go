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
	"github.com/JuaanReis/vorin/internal/collector"
	"github.com/JuaanReis/vorin/internal/model"
	"github.com/JuaanReis/vorin/internal/modules"
	"github.com/JuaanReis/vorin/internal/network"
	"github.com/JuaanReis/vorin/internal/print"
	"github.com/JuaanReis/vorin/pkg"
)

func ParserGET(cfg model.ParserConfigGet) ([]model.Resultado, time.Duration) {
	var resultados []model.Resultado
	var mu sync.Mutex
	var wg sync.WaitGroup
	var resultadoUnico sync.Map
	var reader io.Reader
	var rateLimiter <-chan time.Time
	var compiledRegexTitle, compiledRegexBody *regexp.Regexp
	threadLimiter := make(chan struct{}, cfg.Threads)
	spinnerDone := make(chan bool)
	doneBar := make(chan struct{})
	var fakeBodies [][]byte
	var fakeSizes []int
	var fakeTitles []string

	if cfg.RateLimit > 0 {
		ticker := time.NewTicker(time.Second / time.Duration(cfg.RateLimit))
		rateLimiter = ticker.C
		defer ticker.Stop()
	}

	file, err := pkg.ReadLines(cfg.Wordlist)
	print.FatalIfErr(err)

	if cfg.Shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(file), func(i, j int) { file[i], file[j] = file[j], file[i] })
	}

	client := &http.Client{Timeout: time.Duration(cfg.Timeout) * time.Second}
	if !cfg.Redirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	modules.CreateClientProxy(cfg.Proxy, cfg.Timeout)

	if cfg.Calibrate {
		enderecoBase := strings.TrimRight(strings.Replace(cfg.Endereco, "FUZZ", "", -1), "/")
		for i := 0; i < 3; i++ {
			path := enderecoBase + "/__vorin_" + pkg.RandomString(15) + "__"
			resp, err := http.Get(path)
			print.FatalIfErr(err)
			body, err := io.ReadAll(resp.Body)
			print.FatalIfErr(err)
			resp.Body.Close()
			size, title := collector.DataTargetFake(body)
			fakeBodies = append(fakeBodies, body)
			fakeSizes = append(fakeSizes, size)
			fakeTitles = append(fakeTitles, title)
		}
	}

	if cfg.RegexTitle != "" {
		compiledRegexTitle, err = regexp.Compile("(?i)" + cfg.RegexTitle)
		print.FatalIfErr(err)
	}

	if cfg.RegexBody != "" {
		compiledRegexBody, err = regexp.Compile("(?i)" + cfg.RegexBody)
		print.FatalIfErr(err)
	}

	startTime := time.Now()
	if cfg.Silence {
		go print.Spinner("[Vorin] Running", spinnerDone)
	}

	var current, errors, reqPerSec int32
	if !cfg.Silence {
		go print.UpdateProgressBar(len(file), &current, &errors, &reqPerSec, startTime, doneBar)
	}

	for _, p := range file {
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
				delay := cfg.MinDelay + rand.Float64()*(cfg.MaxDelay-cfg.MinDelay)
				time.Sleep(time.Duration((delay + rand.Float64()*0.1) * float64(time.Second)))
			}

			finalURL := strings.Replace(cfg.Endereco, "FUZZ", p, -1)
			req, err := http.NewRequest("GET", finalURL, nil)
			if err != nil {
				atomic.AddInt32(&errors, 1)
				return
			}

			network.MountHeaders(req, p, cfg.Stealth, cfg.Bypass, cfg.CustomHeaders, cfg.Cookies)
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

			if len(cfg.FilterCode) > 0 && cfg.FilterCode[resp.StatusCode] {
				return
			}

			bodyStr := string(body)
			if analyzer.Content404(bodyStr) {
				return
			}

			lines, size, title, htmlSize, content, text := collector.DataTaget(body)
			if lines == 0 {
				return
			}

			if cfg.FilterBodyContent != "" && !strings.Contains(strings.ToLower(content), strings.ToLower(cfg.FilterBodyContent)) {
				return
			}
			if cfg.FilterTitleContent != "" && !strings.Contains(strings.ToLower(title), strings.ToLower(cfg.FilterTitleContent)) {
				return
			}
			if compiledRegexTitle != nil && !compiledRegexTitle.MatchString(title) {
				return
			}
			if compiledRegexBody != nil && !compiledRegexBody.MatchString(content) {
				return
			}

			if cfg.Calibrate {
				for i := range fakeBodies {
					if collector.IsSameContent(title, fakeTitles[i], cfg, htmlSize, lines, size, fakeSizes[i], content) {
						return
					}
				}
			}

			key := strings.ToLower(p)
			_, loaded := resultadoUnico.LoadOrStore(key, true)
			if loaded {
				return
			}

			tm := time.Since(startTime).Truncate(time.Millisecond)
			label, color := print.StatusColor(resp.StatusCode)

			if cfg.Live {
				fmt.Print("\r\033[K")
				if cfg.StatusOnly {
					fmt.Printf("%s[%3d]%s %-26s\n", color, resp.StatusCode, print.Reset, p)
				} else if cfg.Verbose {
					fmt.Printf("%s[%3d]%s  %4dw  %5dB  %4dL  %6s  %s\n", color, resp.StatusCode, print.Reset, text, htmlSize, lines, tm, label)
					fmt.Printf(" ├─ URL     : %s\n ├─ FUZZ    : %s\n", finalURL, p)
					if title != "" {
						fmt.Printf(" └─ Title   : %s\n\n", title)
					}
				} else {
					fmt.Printf("%s[%3d]%s  %-26s Words: %-6d Size: %-6dB Lines: %-5d %-6s %-11s\n", color, resp.StatusCode, print.Reset, p, text, htmlSize, lines, tm, label)
				}
			}

			mu.Lock()
			resultados = append(resultados, model.Resultado{
				Label: label, Status: resp.StatusCode, URL: p, Title: title,
				Text: text, Size: htmlSize, Lines: lines, Time: tm,
				Color: color, Endereco: finalURL,
			})
			mu.Unlock()
		}(p)
	}

	wg.Wait()
	fmt.Print("\r\033[K\033[1A\r\033[K")
	if !cfg.Silence {
		<-doneBar
	}
	if cfg.Silence {
		spinnerDone <- true
	}
	fmt.Println()
	return resultados, time.Since(startTime)
}
