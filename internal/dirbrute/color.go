package dirbrute

import "fmt"

const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Reset  = "\033[0m"
)

func StatusColor(code int) (string, string) {
	switch {
	case code >= 200 && code < 300:
		return fmt.Sprintf("%s[OK]%s", Green, Reset), Green
	case code >= 300 && code < 400:
		return fmt.Sprintf("%s[REDIRECT]%s", Blue, Reset), Blue
	case code >= 400 && code < 500:
		return fmt.Sprintf("%s[NOT FOUND]%s", Red, Reset), Red
	case code >= 500:
		return fmt.Sprintf("%s[SERVER ERROR]%s", Red, Reset), Red
	default:
		return fmt.Sprintf("%s[UNKNOWN]%s", Red, Reset), Red
	}
}