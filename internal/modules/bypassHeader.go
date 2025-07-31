package modules

// maybe use: Apply headers to bypass mode

func BuildBypassHeaders(path string) map[string]string {
	headers := make(map[string]string)

	headers["X-Original-URL"] = path
	headers["X-Rewrite-URL"] = path
	// headers["X-Forwarded-For"] = RandomIP()
	// headers["X-Client-IP"] = RandomIP()
	headers["X-Http-Method-Override"] = "GET"
	headers["X-Requested-With"] = "XMLHttpRequest"
	headers["Referer"] = "https://google.com"
	headers["Origin"] = "https://google.com"
	headers["X-Forwarded-Proto"] = "https"

	return headers
}