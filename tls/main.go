package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

// CSP represents a Content Security Policy
type CSP struct {
	Directives []Directive
}

// Directive represents a directive in a CSP with a list of sources (could be empty or *)
type Directive struct {
	Name    string
	Sources []string
}

func parseCSP(input string) (CSP, error) {
	return CSP{}, nil
}

func main() {
	dialer := &net.Dialer{
		Timeout: time.Second * 15,
	}
	conn, err := tls.DialWithDialer(dialer, "tcp", "google.com:443", nil)
	if err != nil {
		fmt.Printf("failed to connect: %#v\n", err)
		return
	}
	defer conn.Close()
	connState := conn.ConnectionState()
	fmt.Printf("DEBUG len peer certificates %d\n", len(connState.PeerCertificates))
	fmt.Printf("DEBUG len verified chains  %d\n", len(connState.VerifiedChains))
	for i := 0; i < len(connState.PeerCertificates); i++ {
		// fmt.Printf("DEBUG cert %#v", connState.PeerCertificates[i])
		cert := connState.PeerCertificates[i]
		fmt.Printf("DEBUG peer %#v %d\n", cert.Subject, i)
	}
	for i := 0; i < len(connState.VerifiedChains); i++ {
		chain := connState.VerifiedChains[i]
		for k := 0; k < len(chain); k++ {
			fmt.Printf("DEBUG chain %#v %d\n", chain[k].Subject, k)
		}
	}
}
