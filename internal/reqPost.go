package internal


import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"github.com/JuaanReis/vorin/pkg"
	"github.com/schollz/progressbar/v3"
	"regexp"
)

func ParserPost(endereco string, threads int, userlist string, passlist string, payloadTemplate string,  minDelay float64, maxDelay float64, timeout int, customHeaders map[string]string, randomAgent bool, shuffle bool, live bool, statusOnly bool, regexBody string, regexTitle string) ([]Resultado, time.Duration) {
    var resultados []Resultado
    var mu sync.Mutex
    var wg sync.WaitGroup
	var compiledRegexTitle, compiledRegexBody *regexp.Regexp
    var err error
    if regexTitle != "" {
        compiledRegexTitle, err = regexp.Compile("(?i)" + regexTitle)
        if err != nil {
            fmt.Printf("[ERROR]: Invalid regexTitle: %v\n", err)
            os.Exit(1)
        }
    }
    if regexBody != "" {
        compiledRegexBody, err = regexp.Compile("(?i)" + regexBody)
        if err != nil {
            fmt.Printf("[ERROR]: Invalid regexBody: %v\n", err)
            os.Exit(1)
        }
    }

    fakeUser := "__vorin_fake_user_473827382__"
    fakePass := "__vorin_fake_pass_473827382__"
    fakePayload := strings.ReplaceAll(payloadTemplate, "USERFUZZ", fakeUser)
    fakePayload = strings.ReplaceAll(fakePayload, "PASSFUZZ", fakePass)

    reqFake, err := http.NewRequest("POST", endereco, strings.NewReader(fakePayload))
    if err != nil {
        fmt.Printf("[ERROR]: %v\n", err)
        os.Exit(1)
    }
    reqFake.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    for k, v := range customHeaders {
        reqFake.Header.Set(k, v)
    }
    if randomAgent {
        reqFake.Header.Set("User-Agent", RandomUserAgent())
    }

    users, err := pkg.ReadLines(userlist)
    if err != nil {
        fmt.Printf("[ERROR]: %v\n", err)
        os.Exit(1)
    }
    passes, err := pkg.ReadLines(passlist)
    if err != nil {
        fmt.Printf("[ERROR]: %v\n", err)
        os.Exit(1)
    }

    if shuffle { 
        rand.Seed(time.Now().UnixNano())
        rand.Shuffle(len(users), func(i, j int) { users[i], users[j] = users[j], users[i] })
        rand.Shuffle(len(passes), func(i, j int) { passes[i], passes[j] = passes[j], passes[i] })
    }

    client := &http.Client{
        Timeout: time.Duration(timeout) * time.Second,
    }

    respFake, err := client.Do(reqFake)
    if err != nil || respFake == nil {
        fmt.Printf("[ERROR]: %v\n", err)
        os.Exit(1)
    }
    defer respFake.Body.Close()

    bodyFake, err := io.ReadAll(respFake.Body)
    if err != nil {
        fmt.Printf("[ERROR]: %v\n", err)
        os.Exit(1)
    }
    fakeContent := string(bodyFake)
    fakeTitle := getTitle(fakeContent)
    fakeSize := len(fakeContent)
    fakeStatus := respFake.StatusCode

    totalComb := len(users) * len(passes)
    sem := make(chan struct{}, threads)
    bar := progressbar.NewOptions(totalComb,
        progressbar.OptionSetDescription("Testing POST payloads..."),
        progressbar.OptionShowCount(),
        progressbar.OptionSetWidth(30),
        progressbar.OptionClearOnFinish(),
    )

    ini := time.Now()

    for _, user := range users {
        for _, pass := range passes {
            wg.Add(1)
            sem <- struct{}{}
            go func(user, pass string) {
                defer wg.Done()
                defer func() { <-sem }()
                bar.Add(1)

                if minDelay > 0 || maxDelay > 0 {
                    delayRange := maxDelay - minDelay
                    delay := minDelay + rand.Float64()*delayRange
                    time.Sleep(time.Duration(delay * float64(time.Second)))
                }

                // Replace FUZZ for user and password in the payload
                payload := strings.ReplaceAll(payloadTemplate, "USERFUZZ", user)
                payload = strings.ReplaceAll(payload, "PASSFUZZ", pass)

                req, err := http.NewRequest("POST", endereco, strings.NewReader(payload))
                if err != nil {
                    return
                }
                req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
                for k, v := range customHeaders {
                    req.Header.Set(k, v)
                }
                if randomAgent {
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

                if live {
                    bar.Clear()
                    if statusOnly {
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
                    URL:      endereco,
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
    bar.Clear()
    end := time.Since(ini)
    return resultados, end
}