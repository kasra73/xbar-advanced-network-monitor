// <xbar.title>Bandwidth (KB/s or MB/s)</xbar.title>
// <xbar.version>v0.0.1</xbar.version>
// <xbar.author>Uy Nguyen</xbar.author>
// <xbar.author.github>nguyenvanuyn96</xbar.author.github>
// <xbar.desc>Displays bandwidth usage for the primary interface in MegaBytes/s or KiloBytes/s</xbar.desc>
// <xbar.dependencies>gopsutil</xbar.dependencies>
// <xbar.image>https://user-images.githubusercontent.com/13082464/113498791-ba3ef380-9542-11eb-82e4-76e78cac98b7.png</xbar.image>

package main

import (
	"fmt"
	"os"
	"time"

	"xbar-advanced-network-monitor/pkg/monitor"
)

func main() {
	startCounters, err := monitor.GetCounters()
	if err != nil {
		fmt.Println("Error reading network counters:", err)
		os.Exit(1)
	}

	time.Sleep(monitor.SampleInterval)

	endCounters, err := monitor.GetCounters()
	if err != nil {
		fmt.Println("Error reading network counters:", err)
		os.Exit(1)
	}

	monitor.ProcessAndDisplay(startCounters, endCounters)
}
