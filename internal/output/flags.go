package output

import (
	"strings"
	"os"
)

func NormalizeFlags() {
    longToShort := map[string]string{
        "--url":           "-u",
        "--threads":       "-t",
        "--wordlist":      "-w",
        "--payload":       "-P",
        "--delay":         "-d",
        "--timeout":       "-timeout",
        "--header":        "-H",
        "--status":        "-s",
        "--stealth":       "-stealth",
        "--proxy":         "-proxy",
        "--silence":       "-silence",
        "--live":          "-live",
        "--save-json":     "-save-json",
        "--bypass":        "-bypass",
        "--ext":           "-ext",
        "--rate":          "-rate",
        "--filter-size":   "-filter-size",
        "--filter-line":   "-filter-line",
        "--filter-body":   "-filter-body",
        "--filter-title":  "-filter-title",
        "--random-agent":  "-random-agent",
        "--shuffle":       "-shuffle",
        "--title-contains":"-title-contains",
        "--body-contains": "-body-contains",
        "--regex-body":    "-regex-body",
        "--regex-title":   "-regex-title",
        "--redirect":      "-redirect",
        "--status-only":   "-status-only",
        "--retries":       "-retries",
        "--compare":       "-compare",
        "--random-ip":     "-random-ip",
        "--method":        "-method",
        "--userlist":      "-userlist",
        "--passlist":      "-passlist",
        "--help":          "-h",
    }

    var normalized []string
    for _, arg := range os.Args {
        if strings.HasPrefix(arg, "--") {
            if short, ok := longToShort[arg]; ok {
                normalized = append(normalized, short)
            } else {
                normalized = append(normalized, arg)
            }
        } else {
            normalized = append(normalized, arg)
        }
    }
    os.Args = normalized
}