package internal


import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
	"github.com/JuaanReis/vorin/pkg"
	"regexp"
	"sync/atomic"
)

func ParserPost(cfg ParserConfigPost) ([]Resultado, time.Duration) {
    var resultados []Resultado
    var mu sync.Mutex
    var wg sync.WaitGroup
	var compiledRegexTitle, compiledRegexBody *regexp.Regexp
    var err error
	doneBar := make(chan struct{})
    if cfg.RegexTitle != "" {
        compiledRegexTitle, err = regexp.Compile("(?i)" + cfg.RegexTitle)
        FatalIfErr(err)
    }
    if cfg.RegexBody != "" {
        compiledRegexBody, err = regexp.Compile("(?i)" + cfg.RegexBody)
        FatalIfErr(err)
    }
    
    fakeUser := "__vorin_fake_user_473827382__"
    fakePass := "__vorin_fake_pass_473827382__"
    fakePayload := strings.ReplaceAll(cfg.PayloadTemplate, "USERFUZZ", fakeUser)
    fakePayload = strings.ReplaceAll(fakePayload, "PASSFUZZ", fakePass)

    reqFake, err := http.NewRequest("POST", cfg.Endereco, strings.NewReader(fakePayload))
    FatalIfErr(err)
    reqFake.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    for k, v := range cfg.CustomHeaders {
        reqFake.Header.Set(k, v)
    }
    if cfg.RandomAgent {
        reqFake.Header.Set("User-Agent", RandomUserAgent())
    }

    users, err := pkg.ReadLines(cfg.Userlist)
    FatalIfErr(err)
    passes, err := pkg.ReadLines(cfg.Passlist)
    FatalIfErr(err)

    if cfg.Shuffle { 
        rand.Seed(time.Now().UnixNano())
        rand.Shuffle(len(users), func(i, j int) { users[i], users[j] = users[j], users[i] })
        rand.Shuffle(len(passes), func(i, j int) { passes[i], passes[j] = passes[j], passes[i] })
    }

    client := &http.Client{
        Timeout: time.Duration(cfg.Timeout) * time.Second,
    }

    respFake, err := client.Do(reqFake)
    FatalIfErr(err)
    defer respFake.Body.Close()

    bodyFake, err := io.ReadAll(respFake.Body)
    FatalIfErr(err)

    fakeContent := string(bodyFake)
    fakeTitle := getTitle(fakeContent)
    fakeSize := len(fakeContent)
    fakeStatus := respFake.StatusCode

    totalPaths := len(users) * len(passes)

    sem := make(chan struct{}, cfg.Threads)

    var (
		current    int32
		errors     int32
		reqPerSec  int32
	)

    ini := time.Now()

	if !cfg.Silence {
		go UpdateProgressBar(totalPaths, &current, &errors, &reqPerSec, ini, doneBar)
	}

    for _, user := range users {
        for _, pass := range passes {
            wg.Add(1)
            sem <- struct{}{}
            go func(user, pass string) {
                defer wg.Done()
                defer func() { <-sem }()
			    atomic.AddInt32(&current, 1)
			    atomic.AddInt32(&reqPerSec, 1)

                if cfg.MinDelay > 0 || cfg.MaxDelay > 0 {
                    delayRange := cfg.MaxDelay - cfg.MinDelay
                    delay := cfg.MinDelay + rand.Float64()*delayRange
                    time.Sleep(time.Duration(delay * float64(time.Second)))
                }

                payload := strings.ReplaceAll(cfg.PayloadTemplate, "USERFUZZ", user)
                payload = strings.ReplaceAll(payload, "PASSFUZZ", pass)

                req, err := http.NewRequest("POST", cfg.Endereco, strings.NewReader(payload))
                if err != nil {
                    return
                }
                req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
                for k, v := range cfg.CustomHeaders {
                    req.Header.Set(k, v)
                }
                if cfg.RandomAgent {
                    req.Header.Set("User-Agent", RandomUserAgent())
                }

                start := time.Now()
                resp, err := client.Do(req)
                if err != nil || resp == nil {
                    return
                }
                defer resp.Body.Close()

                body, err := io.ReadAll(resp.Body)
                if err != nil {
                    return
                }
                elapsed := time.Since(start)
                statusLabel, color := StatusColor(resp.StatusCode)
                content := string(body)
                lines, size := CountLinesAndSize(content)
                title := getTitle(content)

                if content == fakeContent || title == fakeTitle || size == fakeSize || resp.StatusCode == fakeStatus {
                    return
                }

                if compiledRegexTitle != nil && !compiledRegexTitle.MatchString(title) {
                    return
                }
                if compiledRegexBody != nil && !compiledRegexBody.MatchString(content) {
                    return
                }

                if cfg.Live {
                    fmt.Print("\r\033[K")
                    if cfg.StatusOnly {
                        fmt.Printf("%s[%3d]%s user=%s pass=%s\n", color, resp.StatusCode, Reset, user, pass)
                    } else {
                        fmt.Printf("%s[%3d]%s user=%s pass=%s Size: %-6dB Lines: %-5d %-6s %-11s\n",
                            color, resp.StatusCode, Reset,
                            user, pass,
                            size,
                            lines,
                            elapsed.Truncate(time.Millisecond),
                            statusLabel,
                        )
                    }
                }

                mu.Lock()
                resultados = append(resultados, Resultado{
                    Label:    statusLabel,
                    Status:   resp.StatusCode,
                    Title:    title,
                    Size:     size,
                    Lines:    lines,
                    Time:     elapsed.Truncate(time.Millisecond),
                    Color:    color,
                    Resposta: content,
					User: user,
					Pass: pass,
                })
                mu.Unlock()
            }(user, pass)
        }
    }

    wg.Wait()
	fmt.Print("\r\033[K")             
	fmt.Print("\033[1A\r\033[K")   
	<-doneBar
	fmt.Println()
    end := time.Since(ini)
    return resultados, end
}