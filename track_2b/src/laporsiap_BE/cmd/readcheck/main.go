// Small private-container read check, not an AI throughput or scale benchmark.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

func main() {
	const base = "http://127.0.0.1:8877/api/cases"
	client := &http.Client{Timeout: 5 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	in := map[string]any{"title": "Fictional bounded read check", "synthetic": true, "parser": "plain", "raw_chat": "This is a fictional read check. No transfer or AI call.", "transfers": []any{}}
	b, _ := json.Marshal(in)
	resp, e := client.Post(base, "application/json", bytes.NewReader(b))
	if e != nil {
		panic("Private app unavailable")
	}
	var c struct {
		Token string `json:"token"`
		Case  struct {
			ID       string `json:"id"`
			Revision int64  `json:"revision"`
		} `json:"case"`
	}
	e = json.NewDecoder(io.LimitReader(resp.Body, 100000)).Decode(&c)
	resp.Body.Close()
	if e != nil || resp.StatusCode != 201 || c.Token == "" {
		panic("Fictional case creation failed")
	}
	endpoint := base + "/" + c.Case.ID
	remove := func() bool {
		b, _ := json.Marshal(map[string]any{"revision": c.Case.Revision, "confirm_delete": true})
		r, _ := http.NewRequest("DELETE", endpoint, bytes.NewReader(b))
		r.Header.Set("Authorization", "Bearer "+c.Token)
		r.Header.Set("Content-Type", "application/json")
		res, e := client.Do(r)
		if e != nil {
			return false
		}
		defer res.Body.Close()
		_, _ = io.Copy(io.Discard, res.Body)
		return res.StatusCode == 200
	}
	defer func() {
		if c.Token != "" {
			_ = remove()
		}
	}()
	const requests = 20
	const concurrency = 4
	jobs := make(chan int, requests)
	for i := 0; i < requests; i++ {
		jobs <- i
	}
	close(jobs)
	timings := make([]float64, requests)
	codes := make([]int, requests)
	var wg sync.WaitGroup
	start := time.Now()
	for worker := 0; worker < concurrency; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				r, _ := http.NewRequest("GET", endpoint, nil)
				r.Header.Set("Authorization", "Bearer "+c.Token)
				begin := time.Now()
				res, e := client.Do(r)
				if e == nil {
					codes[i] = res.StatusCode
					_, _ = io.Copy(io.Discard, res.Body)
					res.Body.Close()
				}
				timings[i] = float64(time.Since(begin).Microseconds()) / 1000
			}
		}()
	}
	wg.Wait()
	elapsed := time.Since(start).Seconds()
	deleted := remove()
	if deleted {
		c.Token = ""
	}
	sort.Float64s(timings)
	passed := 0
	for _, code := range codes {
		if code == 200 {
			passed++
		}
	}
	out := filepath.Join("../../output/hardening", "readcheck-"+time.Now().UTC().Format("2006-01-02T15-04-05.000000000Z")+".json")
	report := map[string]any{"requests": requests, "concurrency": concurrency, "successful_reads": passed, "http_statuses": codes, "wall_seconds": elapsed, "median_request_ms": (timings[9] + timings[10]) / 2, "p95_request_ms": timings[18], "request_ms_sorted": timings, "created_fictional_case_deleted": deleted, "live_ai_calls": 0, "environment": "loopback private app container on developer Mac", "limitations": "Only twenty authenticated reads at concurrency four; not production capacity, AI throughput, availability or a hardware-independent benchmark."}
	b, _ = json.MarshalIndent(report, "", "  ")
	f, e := os.OpenFile(out, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if e != nil {
		panic("Cannot save read check")
	}
	defer f.Close()
	if _, e = f.Write(append(b, '\n')); e != nil {
		panic("Cannot save read check")
	}
	fmt.Printf("Bounded reads: %d/%d; cleaned_up=%t; %s\n", passed, requests, deleted, out)
	if passed != requests || !deleted {
		os.Exit(1)
	}
}
