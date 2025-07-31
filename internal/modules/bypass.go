package modules

import (
	"strings"
)

// maybe use: Apply techniques for bypass

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


