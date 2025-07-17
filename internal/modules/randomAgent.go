package modules

import (
	"fmt"
	"math/rand"
	"time"
)

func RandomUserAgent() string {
	rand.Seed(time.Now().UnixNano())

	platforms := []string{
		"Windows NT 10.0; Win64; x64",
		"Macintosh; Intel Mac OS X 10_15_7",
		"X11; Linux x86_64",
		"Linux; Android 11; SM-G991B",
		"iPhone; CPU iPhone OS 15_5 like Mac OS X",
	}

	browsers := []string{"Chrome", "Firefox", "Safari", "Edge"}

	platform := platforms[rand.Intn(len(platforms))]
	browser := browsers[rand.Intn(len(browsers))]

	switch browser {
	case "Chrome":
		return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%d.0.%d.0 Safari/537.36",
			platform,
			rand.Intn(30)+90,     // Chrome 90-120
			rand.Intn(4000)+1000, // Build
		)

	case "Firefox":
		return fmt.Sprintf("Mozilla/5.0 (%s; rv:%d.0) Gecko/20100101 Firefox/%d.0",
			platform,
			rand.Intn(20)+90,
			rand.Intn(20)+90,
		)

	case "Safari":
		return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/%d.1 Safari/605.1.15",
			platform,
			rand.Intn(5)+13, // Safari 13-18
		)

	case "Edge":
		return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%d.0.%d.0 Safari/537.36 Edg/%d.0.%d.0",
			platform,
			rand.Intn(30)+90,
			rand.Intn(4000)+1000,
			rand.Intn(30)+90,
			rand.Intn(4000)+1000,
		)

	default:
		return "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/99.0.0.0 Safari/537.36"
	}
}
