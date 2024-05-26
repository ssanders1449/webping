package webping

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	// Version describes application version
	Version   = "2.0.0"
	github    = "https://github.com/ssanders1449/webping"
	useragent = fmt.Sprintf("AwsPing/%s (+%s)", Version, github)
)

const (
	// ShowOnlyRegions describes a type of output when only region's name and code printed out
	ShowOnlyRegions = -1
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Duration2ms converts time.Duration to ms (float64)
func Duration2ms(d time.Duration) float64 {
	return float64(d.Nanoseconds()) / 1000 / 1000
}

// mkRandomString returns random string
func mkRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// LatencyOutput prints data into console
type LatencyOutput struct {
	Level   int
	Repeats int
	w       io.Writer
}

// NewOutput creates a new LatencyOutput instance
func NewOutput(level, repeats int) *LatencyOutput {
	return &LatencyOutput{
		Level:   level,
		Repeats: repeats,
		w:       os.Stdout,
	}
}

func (lo *LatencyOutput) show(regions *WebRegions) {
	for _, r := range *regions {
		fmt.Fprintf(lo.w, "%-15s %-s\n", r.Host, r.Name)
	}
}

func (lo *LatencyOutput) show0(regions *WebRegions) {
	for _, r := range *regions {
		fmt.Fprintf(lo.w, "%-25s %20s\n", r.Name, r.GetLatencyStr())
	}
}

func (lo *LatencyOutput) show1(regions *WebRegions) {
	outFmt := "%5v %-15s %-30s %20s\n"
	fmt.Fprintf(lo.w, outFmt, "", "Host", "Region", "Latency")
	for i, r := range *regions {
		fmt.Fprintf(lo.w, outFmt, i, r.Host, r.Name, r.GetLatencyStr())
	}
}

func (lo *LatencyOutput) show2(regions *WebRegions) {
	// format
	outFmt := "%5v %-15s %-25s"
	outFmt += strings.Repeat(" %15s", lo.Repeats) + " %15s\n"
	// header
	outStr := []interface{}{"", "Host", "Region"}
	for i := 0; i < lo.Repeats; i++ {
		outStr = append(outStr, "Try #"+strconv.Itoa(i+1))
	}
	outStr = append(outStr, "Avg Latency")

	// show header
	fmt.Fprintf(lo.w, outFmt, outStr...)

	// each region stats
	for i, r := range *regions {
		outData := []interface{}{strconv.Itoa(i), r.Host, r.Name}
		for n := 0; n < lo.Repeats; n++ {
			outData = append(outData, fmt.Sprintf("%.2f ms",
				Duration2ms(r.Latencies[n])))
		}
		outData = append(outData, fmt.Sprintf("%.2f ms", r.GetLatency()))
		fmt.Fprintf(lo.w, outFmt, outData...)
	}
}

// Show print data
func (lo *LatencyOutput) Show(regions *WebRegions) {
	switch lo.Level {
	case ShowOnlyRegions:
		lo.show(regions)
	case 0:
		lo.show0(regions)
	case 1:
		lo.show1(regions)
	case 2:
		lo.show2(regions)
	}
}

// GetRegions returns a list of regions
func GetRegions() WebRegions {
	return WebRegions{
		NewRegion("gcp Qatar", "gcp-qatar.latency.dev.streaming.synamedia.com"),
	}
}

// CalcLatency returns list of web regions sorted by Latency
func CalcLatency(regions WebRegions, repeats int, useHTTP bool, useHTTPS bool) {
	switch {
	case useHTTP:
		regions.SetCheckType(CheckTypeHTTP)
	case useHTTPS:
		regions.SetCheckType(CheckTypeHTTPS)
	default:
		regions.SetCheckType(CheckTypeTCP)
	}
	regions.SetDefaultTarget()

	var wg sync.WaitGroup
	for n := 1; n <= repeats; n++ {
		wg.Add(len(regions))
		for i := range regions {
			go regions[i].CheckLatency(&wg)
		}
		wg.Wait()
	}

	sort.Sort(regions)
}
