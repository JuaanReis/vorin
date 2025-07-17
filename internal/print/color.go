package print

import "fmt"

const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Reset  = "\033[0m"
)

func StatusColor(code int) (string, string) {
	switch code {
	case 200:
		return fmt.Sprintf("%s[OK]%s", Green, Reset), Green
	case 201:
		return fmt.Sprintf("%s[CREATED]%s", Green, Reset), Green
	case 204:
		return fmt.Sprintf("%s[NO CONTENT]%s", Green, Reset), Green
	case 301, 302, 307, 308:
		return fmt.Sprintf("%s[REDIRECT]%s", Blue, Reset), Blue
	case 400:
		return fmt.Sprintf("%s[BAD REQUEST]%s", Red, Reset), Red
	case 401:
		return fmt.Sprintf("%s[UNAUTHORIZED]%s", Yellow, Reset), Yellow
	case 403:
		return fmt.Sprintf("%s[FORBIDDEN]%s", Yellow, Reset), Yellow
	case 404:
		return fmt.Sprintf("%s[NOT FOUND]%s", Red, Reset), Red
	case 405:
		return fmt.Sprintf("%s[NOT ALLOWED]%s", Yellow, Reset), Yellow
	case 429:
		return fmt.Sprintf("%s[TOO MANY]%s", Yellow, Reset), Yellow
	case 500:
		return fmt.Sprintf("%s[SERVER ERROR]%s", Red, Reset), Red
	case 502:
		return fmt.Sprintf("%s[BAD GATEWAY]%s", Red, Reset), Red
	case 503:
		return fmt.Sprintf("%s[SERVICE UNAVAIL]%s", Red, Reset), Red
	case 504:
		return fmt.Sprintf("%s[TIMEOUT]%s", Red, Reset), Red
	default:
		switch {
		case code >= 200 && code < 300:
			return fmt.Sprintf("%s[2xx]%s", Green, Reset), Green
		case code >= 300 && code < 400:
			return fmt.Sprintf("%s[3xx]%s", Blue, Reset), Blue
		case code >= 400 && code < 500:
			return fmt.Sprintf("%s[4xx]%s", Yellow, Reset), Yellow
		case code >= 500:
			return fmt.Sprintf("%s[5xx]%s", Red, Reset), Red
		default:
			return fmt.Sprintf("%s[UNKNOWN]%s", Red, Reset), Red
		}
	}
}
