package cli

import (
	"flag"
	"strings"
)

var defaultURL = "http://ipv4.download.thinkbroadband.com/5MB.zip"
var defaultString = strings.Join([]string{defaultURL, defaultURL, defaultURL}, ",")


// GetURLs collects the URLS from the CLI
func GetURLs() []string {
    URLs := flag.String("url", defaultString, "Comma separated string with URLs")
    flag.Parse()

	// Split and trim input
    arr := strings.Split(*URLs , ",")
	for ind, value := range arr {
      arr[ind] = strings.TrimSpace(value)
    }

    return arr
}
