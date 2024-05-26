package webping

import (
	"fmt"
	"net"
)

// Targetter is an interface to get target's IP or URL
type Targetter interface {
	GetURL() string
	GetIP() (*net.TCPAddr, error)
}

// WebTarget implements Targetter for Web
type WebTarget struct {
	HTTPS   bool
	Code    string
	Service string
	Rnd     string
}

// GetURL return URL for Web target
func (r *WebTarget) GetURL() string {
	proto := "http"
	if r.HTTPS {
		proto = "https"
	}
	hostname := fmt.Sprintf("%s.%s.amazonweb.com", r.Service, r.Code)
	url := fmt.Sprintf("%s://%s/ping?x=%s", proto, hostname, r.Rnd)
	return url
}

// GetIP return IP for Web target
func (r *WebTarget) GetIP() (*net.TCPAddr, error) {
	tcpURI := fmt.Sprintf("%s.%s.amazonweb.com:80", r.Service, r.Code)
	return net.ResolveTCPAddr("tcp4", tcpURI)
}
