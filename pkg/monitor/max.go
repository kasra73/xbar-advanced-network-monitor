package monitor

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var maxFileName = filepath.Join(userHomeDir(), ".xbar_bandwidth_max_go")

// GetHistoricalMax reads and updates the stored maximum value for scaling.
func GetHistoricalMax(currentMax float64) float64 {
	storedMax := currentMax
	info, err := os.Stat(maxFileName)
	if err == nil {
		age := time.Since(info.ModTime())
		if age <= time.Hour {
			data, err := os.ReadFile(maxFileName)
			if err == nil {
				if val, err := strconvParseFloat(strings.TrimSpace(string(data))); err == nil {
					storedMax = val
				}
			}
		}
	}
	if currentMax >= storedMax {
		storedMax = currentMax
		os.WriteFile(maxFileName, []byte(fmt.Sprintf("%.2f", storedMax)), 0644)
	}
	if storedMax == 0 {
		storedMax = 1000.0
	}
	return storedMax
}

// strconvParseFloat wraps strconv.ParseFloat.
func strconvParseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
