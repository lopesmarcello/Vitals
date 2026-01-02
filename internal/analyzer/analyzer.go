package analyzer

import (
	"context"
	"crypto/tls"
	_ "embed"
	"encoding/json"
	"net/http"
	"net/http/httptrace"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

//go:embed script.js
var browserScript string

type Stats struct {
	URL           string        `json:"url"`
	DNSLookup     time.Duration `json:"dns_lookup"`     // Time to resolve IP
	TCPConnection time.Duration `json:"tcp_connection"` // Time to stablish TCP Connection
	TLSHandshake  time.Duration `json:"tls_handshake"`  // Time to negotiate HTTPS Security
	TTFB          time.Duration `json:"ttfb"`           // Time from sending req to receive to first byte
	TotalTime     time.Duration `json:"total_time"`
	StatusCode    int           `json:"status_code"`
}

type BrowserResult struct {
	FCP   float64  `json:"fcp"`
	Links []string `json:"links"`
}

type LinkHealth struct {
	URL        string        `json:"url"`
	StatusCode int           `json:"status_code"`
	Duration   time.Duration `json:"duration"`
	Error      string        `json:"error,omitempty"`
}

type FullReport struct {
	Network     *Stats         `json:"network"`
	Browser     *BrowserResult `json:"browser"`
	LinksHealth []LinkHealth   `json:"links_health"`
}

func Analyze(ctx context.Context, url string) (*FullReport, error) {
	var (
		netStats *Stats
		netErr   error
		browser  *BrowserResult
		brErr    error
	)

	done := make(chan bool)

	go func() {
		netStats, netErr = AnalyzeNetwork(url)
		done <- true
	}()

	go func() {
		browser, brErr = AnalyzeBrowser(ctx, url)
		done <- true
	}()

	// wait for both
	<-done
	<-done

	if netErr != nil {
		return nil, netErr
	}

	if brErr != nil {
		return nil, brErr
	}

	linkResults := checkLinks(browser.Links)

	return &FullReport{
		Network:     netStats,
		Browser:     browser,
		LinksHealth: linkResults,
	}, nil
}

func AnalyzeBrowser(parentContext context.Context, targetURL string) (*BrowserResult, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true), // required in docker
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-software-rasterizer", true),
		chromedp.Flag("mute-audio", true),
		chromedp.Flag("hide-scrollbars", true),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(parentContext, opts...)
	defer cancelAlloc()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	var result BrowserResult
	var jsonString string

	err := chromedp.Run(ctx,
		chromedp.Navigate(targetURL),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.EvaluateAsDevTools(browserScript, &jsonString),
	)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal([]byte(jsonString), &result); err != nil {
		return nil, err
	}

	return &result, nil
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

func checkLinks(links []string) []LinkHealth {
	var wg sync.WaitGroup
	results := make([]LinkHealth, len(links))

	semaphore := make(chan struct{}, 10)

	for i, link := range links {
		wg.Add(1)

		go func(index int, url string) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			start := time.Now()
			client := http.Client{Timeout: 30 * time.Second}
			resp, err := client.Head(url)

			duration := time.Since(start)

			health := LinkHealth{
				URL:      url,
				Duration: duration,
			}

			if err != nil {
				health.Error = err.Error()
				health.StatusCode = 0
			} else {
				health.StatusCode = resp.StatusCode
				resp.Body.Close()
			}

			results[index] = health
		}(i, link)
	}

	wg.Wait()
	return results
}
