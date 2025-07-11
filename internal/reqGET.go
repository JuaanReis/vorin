package internal

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
	"github.com/JuaanReis/vorin/pkg"
)

func ParserGET(cfg ParserConfigGet) ([]Resultado, time.Duration) {
	var spinnerDone = make(chan bool)
	var resultados []Resultado
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

	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/124.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 Version/15.1 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x86_64) Gecko/20100101 Firefox/113.0",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 15_5 like Mac OS X) AppleWebKit/605.1.15 Mobile/15E148",
		"Mozilla/5.0 (Linux; Android 11; SM-G991B) AppleWebKit/537.36 Chrome/91.0.4472.120 Mobile Safari/537.36",
		"Googlebot/2.1 (+http://www.google.com/bot.html)",
	}

	if cfg.RateLimit > 0 {
		ticker := time.NewTicker(time.Second / time.Duration(cfg.RateLimit))
		rateLimiter = ticker.C 
		defer ticker.Stop()
	}

	rawPaths, err := pkg.ReadLines(cfg.Wordlist)
	FatalIfErr(err)

	var file []string
	for _, path := range rawPaths {
		file = append(file, path)
		if strings.Contains(path, ".") {
			continue
		}

		for _, ext := range cfg.Extension {
			if ext != "" {
				file = append(file, path + ext)
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
		if cfg.Bypass {
			finalPaths = append(finalPaths, ApplyBypassTechniques(path)...)
		} else {
			finalPaths = append(finalPaths, path)
		}
	}


	client := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
	}

	if !cfg.Redirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	CreateClientProxy(cfg.Proxy, cfg.Timeout)

	var fakePath string
	if cfg.Compare != "" {
		fakePath = cfg.Compare
	} else {
		fakePath = "__vorin_this_should_not_exist_40028922__"
	}

	enderecoBase := strings.Replace(cfg.Endereco, "Fuzz", "", -1)
	pathAle := strings.TrimRight(enderecoBase, "/") + "/" + fakePath
	respAle, err := http.Get(pathAle)
	FatalIfErr(err)

	bodyAle, err := io.ReadAll(respAle.Body)
	FatalIfErr(err)

	defer respAle.Body.Close()
	fakeStructureSize, titleAle := DataTargetFake(bodyAle)

	startTime := time.Now()

	if cfg.Silence {
		go Spinner("[Vorin] Running", spinnerDone)
	}

	if cfg.RegexTitle != "" {
		var err error
		compiledRegexTitle, err = regexp.Compile("(?i)" +cfg.RegexTitle)
		FatalIfErr(err)
	}

	if cfg.RegexBody != "" {
		var err error
		compiledRegexBody, err = regexp.Compile("(?i)" + cfg.RegexBody)
		FatalIfErr(err)
	}

	var (
		current    int32
		errors     int32
		reqPerSec  int32
	)

	totalPaths := len(finalPaths)

	if !cfg.Silence {
		go UpdateProgressBar(totalPaths, &current, &errors, &reqPerSec, startTime, doneBar)
	}

	for _, p := range finalPaths {
		wg.Add(1)
		threadLimiter <- struct{}{}
		go func(p string) {
			defer wg.Done()
			defer func() { <- threadLimiter }()
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

			finalURL := strings.Replace(cfg.Endereco, "Fuzz", p, -1)
			start := time.Now()
			req, err := http.NewRequest("GET", finalURL, nil)
			if err != nil {
				return
			}
			MountHeaders(req, p, cfg.Stealth, cfg.Bypass, cfg.CustomHeaders)
			if cfg.RandomIp && (!cfg.Bypass || !cfg.Stealth) {
				ip := RandomIP()
				req.Header.Set("X-Client-IP", ip)
				req.Header.Set("X-Forwarded-For", ip)
				req.Header.Set("CF-Connecting-IP", ip)
			}
			if cfg.RandomAgent && !cfg.Stealth && !cfg.Bypass {
				random := rand.Intn(len(userAgents))
				req.Header.Set("User-Agent", userAgents[random])
			}
			body, resp, err := GetRequestWithRetry(req, client, reader, cfg.Retries)
			if err != nil || resp == nil {
				return
			}
			lines, structureSize, title, htmlSize, content, text := DataTaget(body)
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
			statusLabel, color := StatusColor(resp.StatusCode)
			if IsSameContent(title, titleAle, cfg, htmlSize, lines, structureSize, fakeStructureSize, content) {
				return
			}
			key := strings.ToLower(p)
			if cfg.Code[resp.StatusCode] && !content404(content) && lines > 0 && mustMatch && matchRegexTitle && matchRegexBody {
				if !resultadoUnico[key] {
					resultadoUnico[key] = true
					if cfg.Live {
						fmt.Print("\r\033[K")
						if cfg.StatusOnly {
							fmt.Printf("%s[%3d]%s %-26s\n",
								color, resp.StatusCode, Reset, p,
							)
						} else {
							fmt.Printf("%s[%3d]%s  %-26s Words: %-6d Size: %-6dB Lines: %-5d %-6s %-11s\n",
								color, resp.StatusCode, Reset,
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
					resultados = append(resultados, Resultado{
						Label:  statusLabel,
						Status: resp.StatusCode,
						URL:    p,
						Title:  title,
						Text: text,
						Size:   htmlSize,
						Lines:  lines,
						Time:   tm,
						Color:  color,
					})
					mu.Unlock()
				}
			}
		}(p)   
	}

	fmt.Print("\r\033[K")             
	fmt.Print("\033[1A\r\033[K")    
	wg.Wait()
	close(progressChan)
	<-doneBar
	fmt.Println()
	if cfg.Silence {
		spinnerDone <- true
	}
	end := time.Since(startTime)
	return resultados, end
}
