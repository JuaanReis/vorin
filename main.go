package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/JuaanReis/vorin/internal/core"
	"github.com/JuaanReis/vorin/internal/flags"
	"github.com/JuaanReis/vorin/internal/model"
	"github.com/JuaanReis/vorin/internal/modules"
	"github.com/JuaanReis/vorin/internal/network"
	"github.com/JuaanReis/vorin/internal/output"
	"github.com/JuaanReis/vorin/internal/print"
)

func main() {
	banner, err := os.ReadFile("./internal/banner.txt")
	print.FatalIfErr(err)
	bannerString := string(banner)
	cfg := flags.ParseFlags()

	if cfg.Help {
		flags.PrintHelp()
		os.Exit(0)
	}

	flags.ValidateFlags(cfg)

	if cfg.StatusCodeFlags != "" {
		cfg.FilterCodeFlags = ""
	}

	chosenMethod := strings.ToUpper(cfg.Method)

	cfg.StatusCodeFlags = strings.ReplaceAll(cfg.StatusCodeFlags, " ", "")

	minDelay := float64(0)
	maxDelay := float64(0)

	minDelay, maxDelay, err = modules.ParseDelay(cfg.Delay)
	if err != nil {
		fmt.Printf("[ERROR]: %v\n", err)
		os.Exit(1)
	}

	customHeader := print.ParseHeaderFlags(cfg.HeaderFlags)
	customCookie := print.ParseCookiesFlags(cfg.Cookies)

	if cfg.Stealth {
		if cfg.Rate == 0 {
			cfg.Rate = 15
		}
		if cfg.Threads == 35 {
			cfg.Threads = 30
		}
		if cfg.Timeout == 5 {
			cfg.Timeout = 7
		}
		if minDelay == 0.1 && maxDelay == 0.2 {
			minDelay = 0.2
			maxDelay = 0.2
		}
		customHeader = network.GetRandomHeaders()
	}

	valid := print.ParseStatusCodes(cfg.StatusCodeFlags)
	filterCode := print.ParseStatusCodes(cfg.FilterCodeFlags)

	if cfg.Threads <= 0 || cfg.Threads >= 250 {
		print.PrintError("Thread count must be between 1 and 249.")
		os.Exit(1)
	}

	delayStr := ""
	if minDelay == maxDelay {
		delayStr = fmt.Sprintf("%.1fs", minDelay)
	} else {
		delayStr = fmt.Sprintf("%.1fs-%.1fs", minDelay, maxDelay)
	}

	var rateStr string
	if cfg.Rate > 0 {
		rateStr = fmt.Sprintf("%-3dreq/s", cfg.Rate)
	} else {
		rateStr = "0"
	}
	
	print.PrintHeader(bannerString, cfg.URL, cfg.Wordlist, strconv.Itoa(cfg.Threads), delayStr, fmt.Sprintf("%ds", cfg.Timeout), customHeader, valid, cfg.Stealth, cfg.Proxy, cfg.Silence, cfg.Extension, rateStr, cfg.FilterBody, cfg.FilterTitle, cfg.FilterLine, cfg.FilterSize, cfg.Shuffle, cfg.RandomAgent, cfg.Live, cfg.RegexBody, cfg.RegexTitle, cfg.StatusOnly, cfg.Retries, cfg.Compare, cfg.RandomIp, chosenMethod, cfg.Payload, cfg.Userlist, cfg.Passlist, cfg.Redirect, cfg.NoBanner, cfg.FilterCodeFlags, cfg.Verbose, customCookie, cfg.Calibrate)

	if !cfg.Silence {
		fmt.Println()
		print.PrintLine("_", 80, "Results")
		fmt.Println()
	}

	var listExtension []string
	if cfg.Extension != "" {
		listExtension = strings.Split(cfg.Extension, ",")
	}

	if len(listExtension) > 0 && listExtension[0] != "" && !cfg.Stealth {
		if cfg.Rate == 0 {
			cfg.Rate = 20
		}
		if cfg.Threads == 30 {
			cfg.Threads = 35
		}
		if cfg.Timeout == 8 {
			cfg.Timeout = 6
		}
		minDelay = 0.4
		maxDelay = 0.4
	}

	var resultado []model.Resultado
	var temp time.Duration

	configGet := model.ParserConfigGet{
		Endereco:           cfg.URL,
		Threads:            cfg.Threads,
		Wordlist:           cfg.Wordlist,
		MinDelay:           minDelay,
		MaxDelay:           maxDelay,
		Timeout:            cfg.Timeout,
		CustomHeaders:      customHeader,
		Code:               valid,
		Stealth:            cfg.Stealth,
		Proxy:              cfg.Proxy,
		Silence:            cfg.Silence,
		Live:               cfg.Live,
		Extension:          listExtension,
		RateLimit:          cfg.Rate,
		FilterSize:         cfg.FilterSize,
		FilterLine:         cfg.FilterLine,
		FilterTitle:        cfg.FilterTitle,
		RandomAgent:        cfg.RandomAgent,
		Shuffle:            cfg.Shuffle,
		FilterBody:         cfg.FilterBody,
		RegexBody:          cfg.RegexBody,
		RegexTitle:         cfg.RegexTitle,
		Redirect:           cfg.Redirect,
		StatusOnly:         cfg.StatusOnly,
		Retries:            cfg.Retries,
		Compare:            cfg.Compare,
		RandomIp:           cfg.RandomIp,
		FilterCode:         filterCode,
		Verbose:            cfg.Verbose,
		Cookies:            customCookie,
		Calibrate:          cfg.Calibrate,
	}

	configPost := model.ParserConfigPost{
		Endereco:        cfg.URL,
		Threads:         cfg.Threads,
		Userlist:        cfg.Userlist,
		Passlist:        cfg.Passlist,
		PayloadTemplate: cfg.Payload,
		MinDelay:        minDelay,
		MaxDelay:        maxDelay,
		Timeout:         cfg.Timeout,
		CustomHeaders:   customHeader,
		RandomAgent:     cfg.RandomAgent,
		Shuffle:         cfg.Shuffle,
		Live:            cfg.Live,
		StatusOnly:      cfg.StatusOnly,
		RegexBody:       cfg.RegexBody,
		RegexTitle:      cfg.RegexTitle,
		Silence:         cfg.Silence,
		FilterCode:      filterCode,
		Verbose:         cfg.Verbose,
		Retries:         cfg.Retries,
		Proxy:           cfg.Proxy,
		RateLimit:       cfg.Rate,
		Cookies:         customCookie,
		Calibrate:       cfg.Calibrate,
	}

	switch chosenMethod {
	case "GET":
		resultado, temp = core.ParserGET(configGet)
	case "POST":
		resultado, temp = core.ParserPost(configPost)
	}
	resultadoJson := output.PrepareResultsForJSON(resultado)

	if cfg.StatusOnly && cfg.Live {
		print.SaveJson(*cfg, resultadoJson)

	} else if cfg.StatusOnly {
		print.PrintStatusOnly(resultado, chosenMethod)
		print.SaveJson(*cfg, resultadoJson)

	} else if !cfg.Live {
		if cfg.OutputFile != "" {
			print.SaveJson(*cfg, resultadoJson)
		} else {
			if chosenMethod == "GET" {
				print.PrintGet(resultado, cfg.Verbose)
			} else if chosenMethod == "POST" {
				print.PrintPost(resultado, cfg.Verbose)
			}
		}
	}

		if !cfg.Silence {
			print.PrintLine("_", 80)
			fmt.Printf("\n%s[✓]%s Scan completed in %s%s%s\n\n", print.Green, print.Reset, print.Blue, print.FormatDuration(temp), print.Reset)
		}

		if len(resultado) == 0 {
			fmt.Println(print.Red + "\n[!!] No path found\n" + print.Reset)
		}
	}