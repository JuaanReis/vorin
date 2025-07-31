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

func ParserPost(cfg model.ParserConfigPost) ([]model.Resultado, time.Duration) {
	var resultados []model.Resultado
	var mu sync.Mutex
	var wg sync.WaitGroup
	var compiledRegexTitle, compiledRegexBody *regexp.Regexp
	var err error

	doneBar := make(chan struct{})
	var rateLimiter <-chan time.Time

	if cfg.RegexTitle != "" {
		compiledRegexTitle, err = regexp.Compile("(?i)" + cfg.RegexTitle)
		print.FatalIfErr(err)
	}
	if cfg.RegexBody != "" {
		compiledRegexBody, err = regexp.Compile("(?i)" + cfg.RegexBody)
		print.FatalIfErr(err)
	}
	if cfg.RateLimit > 0 {
		ticker := time.NewTicker(time.Second / time.Duration(cfg.RateLimit))
		rateLimiter = ticker.C
		defer ticker.Stop()
	}

	users, err := pkg.ReadLines(cfg.Userlist)
	print.FatalIfErr(err)
	passes, err := pkg.ReadLines(cfg.Passlist)
	print.FatalIfErr(err)

	if cfg.Shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(users), func(i, j int) { users[i], users[j] = users[j], users[i] })
		rand.Shuffle(len(passes), func(i, j int) { passes[i], passes[j] = passes[j], passes[i] })
	}

	client := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
	}
	modules.CreateClientProxy(cfg.Proxy, cfg.Timeout)

	var fakeContent string
	var fakeTitle string
	var fakeSize int
	var fakeStatus int

	if cfg.Calibrate {
		fakeUser := "__vorin_fake_user__"
		fakePass := "__vorin_fake_pass__"
		fakePayload := strings.ReplaceAll(cfg.PayloadTemplate, "USERFUZZ", fakeUser)
		fakePayload = strings.ReplaceAll(fakePayload, "PASSFUZZ", fakePass)

		reqFake, err := http.NewRequest("POST", cfg.Endereco, strings.NewReader(fakePayload))
		print.FatalIfErr(err)
		reqFake.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		for k, v := range cfg.CustomHeaders {
			reqFake.Header.Set(k, v)
		}
		for k, v := range cfg.Cookies {
			reqFake.AddCookie(&http.Cookie{Name: k, Value: v})
		}
		if cfg.RandomAgent {
			reqFake.Header.Set("User-Agent", modules.RandomUserAgent())
		}

		respFake, err := client.Do(reqFake)
		print.FatalIfErr(err)
		defer respFake.Body.Close()
		bodyFake, err := io.ReadAll(respFake.Body)
		print.FatalIfErr(err)
		fakeContent = string(bodyFake)
		fakeTitle = analyzer.GetTitle(fakeContent)
		fakeSize = len(fakeContent)
		fakeStatus = respFake.StatusCode
	}

	totalPaths := len(users) * len(passes)
	sem := make(chan struct{}, cfg.Threads)

	var current, errors, reqPerSec int32
	ini := time.Now()

	if !cfg.Silence {
		go print.UpdateProgressBar(totalPaths, &current, &errors, &reqPerSec, ini, doneBar)
	}

	for _, user := range users {
		for _, pass := range passes {
			wg.Add(1)
			sem <- struct{}{}
			go func(user, pass string) {
				defer wg.Done()
				defer func() {
					if r := recover(); r != nil {
						atomic.AddInt32(&errors, 1)
					}
					<-sem
				}()

				if rateLimiter != nil {
					<-rateLimiter
				}
				atomic.AddInt32(&current, 1)
				atomic.AddInt32(&reqPerSec, 1)

				if cfg.MinDelay > 0 || cfg.MaxDelay > 0 {
					delay := cfg.MinDelay + rand.Float64()*(cfg.MaxDelay-cfg.MinDelay)
					time.Sleep(time.Duration(delay * float64(time.Second)))
				}

				payload := strings.ReplaceAll(cfg.PayloadTemplate, "USERFUZZ", user)
				payload = strings.ReplaceAll(payload, "PASSFUZZ", pass)
				payloadByte := []byte(payload)

				req, err := http.NewRequest("POST", cfg.Endereco, strings.NewReader(payload))
				if err != nil {
					atomic.AddInt32(&errors, 1)
					return
				}
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				for k, v := range cfg.CustomHeaders {
					req.Header.Set(k, v)
				}
				for k, v := range cfg.Cookies {
					req.AddCookie(&http.Cookie{Name: k, Value: v})
				}
				if cfg.RandomIp {
					ip := modules.RandomIP()
					req.Header.Set("X-Client-IP", ip)
					req.Header.Set("X-Forwarded-For", ip)
					req.Header.Set("CF-Connecting-IP", ip)
				}
				if cfg.RandomAgent {
					req.Header.Set("User-Agent", modules.RandomUserAgent())
				}

				start := time.Now()
				body, resp, err := network.PostRequestWithRetry(req, client, payloadByte, cfg.Retries)
				if err != nil || resp == nil {
					atomic.AddInt32(&errors, 1)
					return
				}

				elapsed := time.Since(start)
				content := string(body)
				lines, size := collector.CountLinesAndSize(content)
				title := analyzer.GetTitle(content)
				word := collector.ExtractTextFromHTML(content)
				text := collector.CountWords(word)
				statusLabel, color := print.StatusColor(resp.StatusCode)

				if cfg.Calibrate {
					if resp.StatusCode == fakeStatus &&
						strings.Contains(strings.ToLower(content), strings.ToLower(fakeContent)) &&
						strings.Contains(strings.ToLower(title), strings.ToLower(fakeTitle)) &&
						abs(size-fakeSize) <= 15 {
						return
					}
				}

				if compiledRegexTitle != nil && !compiledRegexTitle.MatchString(title) {
					return
				}
				if compiledRegexBody != nil && !compiledRegexBody.MatchString(content) {
					return
				}
				if len(cfg.FilterCode) > 0 && !cfg.FilterCode[resp.StatusCode] {
					return
				}

				if cfg.Live {
					fmt.Print("\r\033[K")
					if cfg.StatusOnly {
						fmt.Printf("%s[%3d]%s user=%s pass=%s\n", color, resp.StatusCode, print.Reset, user, pass)
					} else if cfg.Verbose {
						fmt.Printf("%s[%3d]%s  %4dw  %5dB  %4dL  %6s  %s\n", color, resp.StatusCode, print.Reset, text, size, lines, elapsed.Truncate(time.Millisecond), statusLabel)
						fmt.Printf(" ├─ FUZZ    : %s | %s\n", user, pass)
						if title != "" {
							fmt.Printf(" └─ Title   : %s\n\n", title)
						}
					} else {
						fmt.Printf("%s[%3d]%s user=%s pass=%s Size: %-6dB Lines: %-5d %-6s %-11s\n", color, resp.StatusCode, print.Reset, user, pass, size, lines, elapsed.Truncate(time.Millisecond), statusLabel)
					}
				}

				mu.Lock()
				resultados = append(resultados, model.Resultado{
					Label:    statusLabel,
					Status:   resp.StatusCode,
					Title:    title,
					Size:     size,
					Lines:    lines,
					Time:     elapsed.Truncate(time.Millisecond),
					Color:    color,
					Resposta: content,
					User:     user,
					Pass:     pass,
				})
				mu.Unlock()
			}(user, pass)
		}
	}

	wg.Wait()
	fmt.Print("\r\033[K")
	fmt.Print("\033[1A\r\033[K")

	if !cfg.Silence {
		<-doneBar
	}

	fmt.Println()
	end := time.Since(ini)
	return resultados, end
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
