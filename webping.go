package webping

import (
	"fmt"
	"sync"
	"time"
)

// CheckType describes a type for a check
type CheckType int

const (
	// CheckTypeTCP is TCP type of check
	CheckTypeTCP CheckType = iota
	// CheckTypeHTTP is HTTP type of check
	CheckTypeHTTP
	// CheckTypeHTTPS is HTTPS type of check
	CheckTypeHTTPS
)

// --------------------------------------------

// WebRegion description of the Web EC2 region
type WebRegion struct {
	Name      string
	Host      string
	Latencies []time.Duration
	Error     error
	CheckType CheckType

	Target  Targetter
	Request Requester
}

// NewRegion creates a new region with a name and code
func NewRegion(name, host string) WebRegion {
	return WebRegion{
		Name:      name,
		Host:      host,
		CheckType: CheckTypeTCP,
		Request:   NewWebRequest(),
	}
}

// CheckLatency does a latency check for a region
func (r *WebRegion) CheckLatency(wg *sync.WaitGroup) {
	defer wg.Done()

	if r.CheckType == CheckTypeHTTP || r.CheckType == CheckTypeHTTPS {
		r.checkLatencyHTTP(r.CheckType == CheckTypeHTTPS)
	} else {
		r.checkLatencyTCP()
	}
}

// checkLatencyHTTP Test Latency via HTTP
func (r *WebRegion) checkLatencyHTTP(https bool) {
	url := r.Target.GetURL()
	l, err := r.Request.Do(useragent, url, RequestTypeHTTP)
	if err != nil {
		r.Error = err
		return
	}
	r.Latencies = append(r.Latencies, l)
}

// checkLatencyTCP Test Latency via TCP
func (r *WebRegion) checkLatencyTCP() {
	tcpAddr, err := r.Target.GetIP()
	if err != nil {
		r.Error = err
		return
	}

	l, err := r.Request.Do(useragent, tcpAddr.String(), RequestTypeTCP)
	if err != nil {
		r.Error = err
		return
	}
	r.Latencies = append(r.Latencies, l)
}

// GetLatency returns Latency in ms
func (r *WebRegion) GetLatency() float64 {
	sum := float64(0)
	for _, l := range r.Latencies {
		sum += Duration2ms(l)
	}
	return sum / float64(len(r.Latencies))
}

// GetLatencyStr returns Latency in string
func (r *WebRegion) GetLatencyStr() string {
	if r.Error != nil {
		return r.Error.Error()
	}
	return fmt.Sprintf("%.2f ms", r.GetLatency())
}

// --------------------------------------------

// WebRegions slice of the WebRegion
type WebRegions []WebRegion

// Len returns a count of regions
func (rs WebRegions) Len() int {
	return len(rs)
}

// Less return a result of latency compare between two regions
func (rs WebRegions) Less(i, j int) bool {
	return rs[i].GetLatency() < rs[j].GetLatency()
}

// Swap two regions by index
func (rs WebRegions) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

// SetCheckType sets Check Type for all regions
func (rs WebRegions) SetCheckType(checkType CheckType) {
	for i := range rs {
		rs[i].CheckType = checkType
	}
}

// SetDefaultTarget sets default target instance
func (rs WebRegions) SetDefaultTarget() {
	rs.SetTarget(func(r *WebRegion) {
		r.Target = &WebTarget{
			HTTPS:   r.CheckType == CheckTypeHTTPS,
			Host:    r.Host,
			Rnd:     mkRandomString(13),
		}
	})
}

// SetTarget sets default target instance for all regions
func (rs WebRegions) SetTarget(fn func(r *WebRegion)) {
	for i := range rs {
		fn(&rs[i])
	}
}
