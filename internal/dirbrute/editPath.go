package dirbrute

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"github.com/JuaanReis/vorin/pkg"
	"github.com/schollz/progressbar/v3"
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

var reqCounter int32

func Parser(endereco string, threads int, wordlist string, minDelay float64, maxDelay float64, timeout int, customHeaders map[string]string, code map[int]bool, stealth bool, proxy string, silence bool, live bool, bypass bool, extension []string) ([]Resultado, time.Duration) {
	var resultados []Resultado
	var mu sync.Mutex
	var wg sync.WaitGroup
	var reader io.Reader

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
					defer bar.Add(1)
				}

				if stealth {
					atomic.AddInt32(&reqCounter, 1)
    			if atomic.LoadInt32(&reqCounter)%100 == 0 {
        		time.Sleep(10 * time.Second)
    			}
				}
				if stealth || (minDelay > 0 && maxDelay > 0) {
					delayRange := maxDelay - minDelay
					delay := minDelay + rand.Float64()*delayRange
					time.Sleep(time.Duration(delay * float64(time.Second)))
				}


				finalURL := strings.Replace(endereco, "Fuzz", p, -1)
				start := time.Now()

				req, err := http.NewRequest("GET", finalURL, nil)
				if err != nil {
					return
				}

				MountHeaders(req, p, stealth, bypass, customHeaders)

				body, resp, err := GetRequest(req, client, reader)
				if err != nil {
					return
				}

				lines, structureSize, title, htmlSize, content := DataTaget(body)

				elapsed := time.Since(start)
				isSameContent := title == titleAle || structureSize == fakeStructureSize
				tm := elapsed.Truncate(time.Millisecond)
				statusLabel, color := StatusColor(resp.StatusCode)

				if code[resp.StatusCode] && !content404(content) && !isSameContent && lines > 0{
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
			}(p)
		}
	}

	wg.Wait()
	if silence {
		spinnerDone <- true
	}

	bar.Clear()
	end := time.Since(ini)
	return resultados, end
}
