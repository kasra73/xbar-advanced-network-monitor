package monitor

import (
	"fmt"
	"math"
	"strings"

	"github.com/fatih/color"
	psnet "github.com/shirou/gopsutil/v3/net"
)

// ifaceStat holds per-interface bandwidth statistics.
type ifaceStat struct {
	name    string
	inRate  float64
	outRate float64
}

// ProcessAndDisplay computes the bandwidth differences, prints aggregated bars and interface stats.
func ProcessAndDisplay(startCounters, endCounters map[string]psnet.IOCountersStat) {
	totalInKBps, totalOutKBps := 0.0, 0.0
	var positiveStats []ifaceStat

	for name, start := range startCounters {
		end, ok := endCounters[name]
		if !ok {
			continue
		}
		inBytes := float64(end.BytesRecv - start.BytesRecv)
		outBytes := float64(end.BytesSent - start.BytesSent)
		inKBps := (inBytes * 2) / 1000.0
		outKBps := (outBytes * 2) / 1000.0

		totalInKBps += inKBps
		totalOutKBps += outKBps

		if inKBps > 0 || outKBps > 0 {
			positiveStats = append(positiveStats, ifaceStat{
				name:    name,
				inRate:  inKBps,
				outRate: outKBps,
			})
		}
	}

	printTotalBandwidth(totalInKBps, totalOutKBps)

	fmt.Println("---")

	for _, stat := range positiveStats {
		inStr := FormatInterfaceBW(stat.inRate)
		outStr := FormatInterfaceBW(stat.outRate)
		fmt.Printf("%s:\\t↓%s\\n\\t↑%s | color='#3487c7' font='Menlo' size=8\n", stat.name, inStr, outStr)
	}
}

func printTotalBandwidth(inTotal, outTotal float64) {
	inDisplay := FormatBW(inTotal)
	outDisplay := FormatBW(outTotal)

	currentMax := inTotal
	if outTotal > currentMax {
		currentMax = outTotal
	}
	storedMax := GetHistoricalMax(currentMax)

	maxBars := 4
	inBars := int(math.Round((inTotal / storedMax) * float64(maxBars)))
	outBars := int(math.Round((outTotal / storedMax) * float64(maxBars)))
	if inTotal > 0 && inBars == 0 {
		inBars = 1
	}
	if outTotal > 0 && outBars == 0 {
		outBars = 1
	}
	if inBars > maxBars {
		inBars = maxBars
	}
	if outBars > maxBars {
		outBars = maxBars
	}
	emptyIn := maxBars - inBars
	emptyOut := maxBars - outBars

	// Define bar color sequences: green, yellow, orange (HiMagenta), red.
	barColors := []*color.Color{
		color.New(color.FgGreen),
		color.New(color.FgYellow),
		color.New(color.FgHiMagenta), // used as orange
		color.New(color.FgRed),
	}
	gray := color.New(color.FgHiBlack)

	// Choose text color based on last filled bar (if any)
	inTextColor := gray
	if inBars > 0 {
		inTextColor = barColors[inBars-1]
	}
	outTextColor := gray
	if outBars > 0 {
		outTextColor = barColors[outBars-1]
	}

	// Print incoming bandwidth row with colored text and bars.
	fmt.Print("↓")
	inTextColor.Printf("%s ", inDisplay)
	for i := 0; i < inBars; i++ {
		barColors[i].Printf("▮")
	}
	gray.Print(strings.Repeat("▯", emptyIn))
	fmt.Print("\\n")

	// Print outgoing bandwidth row with colored text and bars.
	fmt.Print("↑")
	outTextColor.Printf("%s ", outDisplay)
	for i := 0; i < outBars; i++ {
		barColors[i].Printf("▮")
	}
	gray.Print(strings.Repeat("▯", emptyOut))
	fmt.Print(" | ansi=true font='Menlo' size=7\n")
}
