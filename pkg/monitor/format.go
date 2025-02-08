package monitor

import (
	"fmt"
	"math"
)

// FormatInterfaceBW returns a string for individual interface bandwidth with two decimals.
func FormatInterfaceBW(kbps float64) string {
	if kbps < 800.0 {
		return fmt.Sprintf("%.2fKB/s", kbps)
	} else if kbps < 1000000.0 {
		return fmt.Sprintf("%.2fMB/s", kbps/1000.0)
	}
	return fmt.Sprintf("%.2fGB/s", kbps/1000000.0)
}

// FormatBW returns a rounded string (no decimals) for total bandwidth.
func FormatBW(kbps float64) string {
	if kbps < 800.0 {
		return fmt.Sprintf("%4dKB/s", int(math.Round(kbps)))
	} else if kbps < 1000000.0 {
		return fmt.Sprintf("%4dMB/s", int(math.Round(kbps/1000.0)))
	}
	return fmt.Sprintf("%4dGB/s", int(math.Round(kbps/1000000.0)))
}
