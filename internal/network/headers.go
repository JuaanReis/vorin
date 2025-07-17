package network

import (
	"math/rand"
	"time"
)

func GetRandomHeaders() map[string]string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/124.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 Version/15.1 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x86_64) Gecko/20100101 Firefox/113.0",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 15_5 like Mac OS X) AppleWebKit/605.1.15 Mobile/15E148",
		"Mozilla/5.0 (Linux; Android 11; SM-G991B) AppleWebKit/537.36 Chrome/91.0.4472.120 Mobile Safari/537.36",
		"Googlebot/2.1 (+http://www.google.com/bot.html)",
	}

	accepts := []string{
		"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"application/json, text/plain, */*",
		"text/html",
		"*/*",
	}

	languages := []string{
		"en-US,en;q=0.9",
		"pt-BR,pt;q=0.8,en-US;q=0.6,en;q=0.4",
		"fr-FR,fr;q=0.9",
		"de-DE,de;q=0.9,en;q=0.8",
	}

	cacheControl := []string{
		"no-cache",
		"max-age=0",
		"no-store",
	}

	upInsecure := []string{
		"1",
		"0",
	}

	dnt := []string{
		"1",
		"0",
	}

	rand.Seed(time.Now().UnixNano())

	return map[string]string{
		"User-Agent":                userAgents[rand.Intn(len(userAgents))],
		"Accept":                    accepts[rand.Intn(len(accepts))],
		"Accept-Language":           languages[rand.Intn(len(languages))],
		"Accept-Encoding":           "gzip, deflate",
		"Connection":                "keep-alive",
		"Cache-Control":             cacheControl[rand.Intn(len(cacheControl))],
		"Upgrade-Insecure-Requests": upInsecure[rand.Intn(len(upInsecure))],
		"DNT":                       dnt[rand.Intn(len(dnt))],
		"Sec-Fetch-Dest":            "document",
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-Site":            "none",
		"Sec-Fetch-User":            "?1",
	}
}
