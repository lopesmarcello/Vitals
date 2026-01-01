package analyzer

import (
	"crypto/tls"
	"net/http"
	"net/http/httptrace"
	"time"
)

type Stats struct {
	URL           string
	DNSLookup     time.Duration // Time to resolve IP
	TCPConnection time.Duration // Time to stablish TCP Connection
	TLSHandshake  time.Duration // Time to negotiate HTTPS Security
	TTFB          time.Duration // Time from sending req to receive to first byte
	TotalTime     time.Duration
	StatusCode    int
}

func AnalyzeNetwork(targetURL string) (*Stats, error) {
	var (
		start        = time.Now()
		dnsStart     time.Time
		dnsDone      time.Time
		connectStart time.Time
		connectDone  time.Time
		wroteRequest time.Time
		firstByte    time.Time
	)

	stats := &Stats{URL: targetURL}

	trace := &httptrace.ClientTrace{
		// DNS Hooks
		DNSStart: func(i httptrace.DNSStartInfo) { dnsStart = time.Now() },
		DNSDone:  func(i httptrace.DNSDoneInfo) { dnsDone = time.Now() },

		// Connection (TCP) Hooks
		ConnectStart: func(network, addr string) { connectStart = time.Now() },
		ConnectDone:  func(network, addr string, err error) { connectDone = time.Now() },

		// The moment it finishes sending request headers and body
		WroteRequest: func(w httptrace.WroteRequestInfo) { wroteRequest = time.Now() },

		GotFirstResponseByte: func() { firstByte = time.Now() },
	}

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, err
	}

	contextWithTrace := httptrace.WithClientTrace(req.Context(), trace)
	req = req.WithContext(contextWithTrace)

	// disable Keep-Alives for cold-start experience
	transport := &http.Transport{
		DisableKeepAlives: true,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true}, // ignore bad certs
	}

	client := &http.Client{Transport: transport, Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	stats.StatusCode = resp.StatusCode

	if !dnsStart.IsZero() {
		stats.DNSLookup = dnsDone.Sub(dnsStart)
	}

	if !connectStart.IsZero() {
		stats.TCPConnection = connectDone.Sub(connectStart)
	}

	if !wroteRequest.IsZero() && !firstByte.IsZero() {
		stats.TTFB = firstByte.Sub(wroteRequest)
	}

	stats.TotalTime = time.Since(start)

	return stats, nil
}
