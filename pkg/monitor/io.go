package monitor

import (
	"os"
	"strings"
	"time"

	psnet "github.com/shirou/gopsutil/v3/net"
)

// SampleInterval defines the measurement interval.
var SampleInterval = 500 * time.Millisecond

// GetCounters returns IO counters for interfaces starting with "en".
func GetCounters() (map[string]psnet.IOCountersStat, error) {
	counters, err := psnet.IOCounters(true)
	if err != nil {
		return nil, err
	}
	result := make(map[string]psnet.IOCountersStat)
	for _, c := range counters {
		if strings.HasPrefix(c.Name, "en") {
			result[c.Name] = c
		}
	}
	return result, nil
}

// userHomeDir returns the user's home directory.
func userHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return home
}
