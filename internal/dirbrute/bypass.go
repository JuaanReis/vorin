package dirbrute

import (
	"strings"
	"math/rand"
	"fmt"
)

func ApplyBypassTechniques(path string) []string {
    var bypassedPaths []string
    bypassedPaths = append(bypassedPaths, path)
    bypassedPaths = append(bypassedPaths, strings.ToUpper(path))
    bypassedPaths = append(bypassedPaths, path+"?")
    bypassedPaths = append(bypassedPaths, "/./"+path)
    bypassedPaths = append(bypassedPaths, "/..;/"+path)
    bypassedPaths = append(bypassedPaths, "/%2e/"+path)
    bypassedPaths = append(bypassedPaths, "/"+path+"/.")

    return bypassedPaths
}

func randomIP() string {
	return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
}

func BuildBypassHeaders(path string) map[string]string {
	headers := make(map[string]string)

	headers["X-Original-URL"] = path
	headers["X-Rewrite-URL"] = path
	headers["X-Forwarded-For"] = randomIP()
	headers["X-Client-IP"] = randomIP()
	headers["X-Http-Method-Override"] = "GET"
	headers["X-Requested-With"] = "XMLHttpRequest"
	headers["Referer"] = "https://google.com"
	headers["Origin"] = "https://google.com"
	headers["X-Forwarded-Proto"] = "https"

	return headers
}
