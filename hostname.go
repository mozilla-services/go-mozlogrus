package mozlogrus

import (
	"log"
	"os"
)

var hostname string

func Hostname() string {
	if hostname == "" {
		var err error
		hostname, err = os.Hostname()
		if err != nil {
			log.Printf("Can't resolve hostname: %v", err)
			hostname = "unknown"
		}
	}
	return hostname
}
