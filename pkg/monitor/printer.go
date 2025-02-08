package monitor

import (
	"fmt"
	"math"
	"strings"

	psnet "github.com/shirou/gopsutil/v3/net"
)

// New helper type and functions to print colors manually.
type myColor struct {
	r, g, b int
}

func RGB(r, g, b int) myColor {
	return myColor{r, g, b}
}

func (c myColor) Sprint(s string) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", c.r, c.g, c.b, s)
}

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

	// Define bar colors using our RGB helper: green, yellow, orange, red.
	barColors := []myColor{
		RGB(0, 255, 0),   // green
		RGB(255, 255, 0), // yellow
		RGB(255, 165, 0), // orange
		RGB(255, 0, 0),   // red
	}
	gray := RGB(128, 128, 128)

	// Choose text color based on last filled bar (if any)
	inTextColor := gray
	if inBars > 0 {
		inTextColor = barColors[inBars-1]
	}
	outTextColor := gray
	if outBars > 0 {
		outTextColor = barColors[outBars-1]
	}

	// Capture output on a strings.Builder instead of printing directly.
	var sb strings.Builder

	// Build incoming bandwidth row.
	sb.WriteString(inTextColor.Sprint(inDisplay) + " ")
	sb.WriteString("▼ ")
	for i := 0; i < inBars; i++ {
		sb.WriteString(barColors[i].Sprint("▮"))
	}
	sb.WriteString(gray.Sprint(strings.Repeat("▯", emptyIn)))
	sb.WriteString("\\n")

	// Build outgoing bandwidth row.
	sb.WriteString(outTextColor.Sprint(outDisplay) + " ")
	sb.WriteString("▲ ")
	for i := 0; i < outBars; i++ {
		sb.WriteString(barColors[i].Sprint("▮"))
	}
	sb.WriteString(gray.Sprint(strings.Repeat("▯", emptyOut)))
	// Print the captured string with quoted formatting.
	esc := literalEscape(sb.String())
	esc = strings.Trim(esc, "\"")
	fmt.Printf("%s | size=8\n", esc)
}

func literalEscape(s string) string {
	// Replace every actual ESC character with the literal text "\x1b"
	return strings.ReplaceAll(s, "\x1b", `\x1b`)
}
