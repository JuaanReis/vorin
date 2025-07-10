package internal

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
	"github.com/JuaanReis/vorin/pkg"
	"github.com/schollz/progressbar/v3"
)

type Resultado struct {
	Status int
	Resposta string
	URL    string
	Title  string
	Size   int
	Lines  int
	Time   time.Duration
	Label  string
	Color  string
	User string
	Pass string
}

var spinnerDone = make(chan bool)

func ParserGET(endereco string, threads int, wordlist string, minDelay float64, maxDelay float64, timeout int, customHeaders map[string]string, code map[int]bool, stealth bool, proxy string, silence bool, live bool, bypass bool, extension []string, rateLimit int, filterSize int, filterLine int, filterTitle string, randomAgent bool, shuffle bool, filterTitleContent string, filterBodyContent string, filterBody string, regexBody string, regexTitle string, redirect bool, statusOnly bool, retries int, compare string, randomIp bool) ([]Resultado, time.Duration) {
	var resultados []Resultado
	var mu sync.Mutex
	var wg sync.WaitGroup
	var reader io.Reader
	var rateLimiter <-chan time.Time
	var resultadoUnico = make(map[string]bool)
	progressChan := make(chan struct{})
	var compiledRegexTitle *regexp.Regexp
	var compiledRegexBody *regexp.Regexp

	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/124.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 Version/15.1 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x86_64) Gecko/20100101 Firefox/113.0",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 15_5 like Mac OS X) AppleWebKit/605.1.15 Mobile/15E148",
		"Mozilla/5.0 (Linux; Android 11; SM-G991B) AppleWebKit/537.36 Chrome/91.0.4472.120 Mobile Safari/537.36",
		"Googlebot/2.1 (+http://www.google.com/bot.html)",
	}

	if rateLimit > 0 {
		rateLimiter = time.Tick(time.Second / time.Duration(rateLimit))
	}

	sem := make(chan struct{}, threads)

	rawPaths, err := pkg.ReadLines(wordlist)
	if err != nil {
		fmt.Printf("[ERROR]: %v\n", err)
		os.Exit(1)
	}

	var file []string
	for _, path := range rawPaths {
		file = append(file, path)
		if strings.Contains(path, ".") {
			continue
		}

		for _, ext := range extension {
			if ext != "" {
				file = append(file, path+ext)
			}
		}
	}

	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(file), func(i, j int) {
			file[i], file[j] = file[j], file[i]
		})
	}

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	if !redirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	CreateClientProxy(proxy, timeout)

	var fakePath string
	if compare != "" {
		fakePath = compare
	} else {
		fakePath = "__vorin_this_should_not_exist_473827382__"
	}
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

	totalPaths := 0
	for _, path := range file {
		if bypass {
			totalPaths += len(ApplyBypassTechniques(path))
		} else {
			totalPaths++
		}
	}

	bar := progressbar.NewOptions((totalPaths),
		progressbar.OptionSetDescription("Testing paths..."),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(30),
		progressbar.OptionClearOnFinish(),
	)

	ini := time.Now()

	if silence {
		go Spinner("[Vorin] Running", spinnerDone)
	}

	go func() {
		for range progressChan {
			bar.Add(1)
		}
	}()

	if regexTitle != "" {
		var err error
		compiledRegexTitle, err = regexp.Compile("(?i)" + regexTitle)
		if err != nil {
			fmt.Printf("[ERROR regex-title]: %v\n", err)
			os.Exit(1)
		}
	}

	if regexBody != "" {
		var err error
		compiledRegexBody, err = regexp.Compile("(?i)" + regexBody)
		if err != nil {
			fmt.Printf("[ERROR regex-body]: %v\n", err)
			os.Exit(1)
		}
	}

	for _, path := range file {
		paths := []string{path}
		if bypass {
			paths = ApplyBypassTechniques(path)
		}
		for _, p := range paths {
			wg.Add(1)
			sem <- struct{}{}
			go func(p string) {
				defer wg.Done()
				defer func() { <-sem }()
				if !silence {
					progressChan <- struct{}{}
				}

				if rateLimiter != nil {
					<-rateLimiter
				}

				if stealth || (minDelay > 0 && maxDelay > 0) {
					delayRange := maxDelay - minDelay
					delay := minDelay + rand.Float64()*delayRange
					jitter := rand.Float64() * 0.1
					time.Sleep(time.Duration((delay + jitter) * float64(time.Second)))
				}

				finalURL := strings.Replace(endereco, "Fuzz", p, -1)
				start := time.Now()

				req, err := http.NewRequest("GET", finalURL, nil)
				if err != nil {
					return
				}

				MountHeaders(req, p, stealth, bypass, customHeaders)

				if randomIp && (!bypass || !stealth) {
					ip := RandomIP()
					req.Header.Set("X-Client-IP", ip)
					req.Header.Set("X-Forwarded-For", ip)
					req.Header.Set("CF-Connecting-IP", ip)
				}

				if randomAgent && !stealth && !bypass {
					random := rand.Intn(len(userAgents))
					req.Header.Set("User-Agent", userAgents[random])
				}

				body, resp, err := GetRequestWithRetry(req, client, reader, retries)
				if err != nil || resp == nil {
					return
				}

				lines, structureSize, title, htmlSize, content := DataTaget(body)

				elapsed := time.Since(start)

				mustMatch := true
				if filterBodyContent != "" && !strings.Contains(strings.ToLower(content), strings.ToLower(filterBodyContent)) {
					mustMatch = false
				}

				if filterTitleContent != "" && !strings.Contains(strings.ToLower(title), strings.ToLower(filterTitleContent)) {
					mustMatch = false
				}

				titleMatch := title == titleAle || title == filterTitle
				sizeTooSmall := htmlSize <= filterSize
				linesTooFew := lines <= filterLine || lines == 0
				structureMatch := structureSize == fakeStructureSize
				erroBody := false
				if filterBody != "" {
					erroBody = strings.Contains(strings.ToLower(content), strings.ToLower(filterBody))
				}
				isSameContent := titleMatch || structureMatch || sizeTooSmall || linesTooFew || erroBody

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

				key := strings.ToLower(p)

				if code[resp.StatusCode] && !content404(content) && !isSameContent && lines > 0 && mustMatch && matchRegexTitle && matchRegexBody {
					if !resultadoUnico[key] {
						resultadoUnico[key] = true
						if live {
							bar.Clear()
							if statusOnly {
								fmt.Printf("%s[%3d]%s %-26s\n",
									color, resp.StatusCode, Reset, p,
								)
							} else {
								fmt.Printf("%s[%3d]%s  %-26s Size: %-6dB Lines: %-5d %-6s %-11s\n",
									color, resp.StatusCode, Reset,
									p,
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
	}

	wg.Wait()
	close(progressChan)
	if silence {
		spinnerDone <- true
	}

	bar.Clear()
	end := time.Since(ini)
	return resultados, end
}
