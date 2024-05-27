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
	Host    string
	Rnd     string
}

// GetURL return URL for Web target
func (r *WebTarget) GetURL() string {
	proto := "http"
	if r.HTTPS {
		//proto = "https"
		fmt.Printf("HTTPS not supported, using http\n")
	}
	url := fmt.Sprintf("%s://%s.dev.streaming.synamedia.com/ping?x=%s", proto, r.Host, r.Rnd)
	return url
}

// GetIP return IP for Web target
func (r *WebTarget) GetIP() (*net.TCPAddr, error) {
	tcpURI := fmt.Sprintf("%s.dev.streaming.synamedia.com:80", r.Host)
	return net.ResolveTCPAddr("tcp4", tcpURI)
}
