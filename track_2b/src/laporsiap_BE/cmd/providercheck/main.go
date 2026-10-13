// providercheck performs configuration/catalog/policy/quota reads only.
// It never requests a completion, creates a case, or prints credentials.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v5"
	"laporsiap/internal/infrastructure/apertus"
	"laporsiap/internal/infrastructure/config"
	"net/http"
	"os"
	"time"
)

func main() {
	output := flag.String("output", "", "Optional new JSON report path (never overwrite)")
	flag.Parse()
	p := apertus.NewProvider(apertus.ConfigFromEnv())
	report := map[string]any{"checked_at": time.Now().UTC().Format(time.RFC3339), "provider_label": p.Info().Label, "requested_model": p.Model, "prompt_version": apertus.PromptVersion, "prompt_sha256": apertus.PromptSHA256(), "completion_requests": 0, "case_writes": 0, "configuration_valid": false, "catalog_model_listed": false, "runtime_prompt_matches": false, "quota_verified": false, "limitations": []string{"Catalog listing does not verify completion routing or model weights.", "Quota values are a read-only snapshot, not a reservation."}}
	ok := true
	if err := p.ValidateConfiguration(); err != nil {
		report["configuration_error"] = err.Error()
		ok = false
	} else {
		report["configuration_valid"] = true
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		err := p.CheckAvailability(ctx)
		cancel()
		if err != nil {
			report["catalog_error"] = err.Error()
			ok = false
		} else {
			report["catalog_model_listed"] = true
		}
	}
	client := &http.Client{Timeout: 5 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	resp, err := client.Get("http://127.0.0.1:8787/api/policy")
	if err == nil {
		var policy struct {
			Hash    string `json:"prompt_sha256"`
			Version string `json:"prompt_version"`
		}
		if resp.StatusCode == 200 && json.NewDecoder(http.MaxBytesReader(nil, resp.Body, 65536)).Decode(&policy) == nil {
			report["runtime_prompt_matches"] = policy.Hash == apertus.PromptSHA256() && policy.Version == apertus.PromptVersion
		}
		resp.Body.Close()
	}
	if report["runtime_prompt_matches"] != true {
		report["runtime_error"] = "Running backend unavailable or embedded prompt differs; restart before evaluation."
		ok = false
	}
	limits, err := config.LoadSafetyLimits()
	if err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		conn, e := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
		if e == nil {
			var attempts, active, cases int
			e = conn.QueryRow(ctx, `SELECT count(*) FILTER (WHERE created_at>=now()-interval '24 hours'), count(*) FILTER (WHERE finished_at IS NULL AND lease_until>now()), (SELECT count(*) FROM cases) FROM ai_attempts`).Scan(&attempts, &active, &cases)
			if e == nil {
				report["quota_verified"] = true
				report["attempts_last_24h"] = attempts
				report["global_daily_limit"] = limits.AIGlobalDaily
				report["active_leases"] = active
				report["concurrent_limit"] = limits.AIConcurrent
				report["stored_cases"] = cases
				report["max_cases"] = limits.MaxCases
				room := attempts+2 <= limits.AIGlobalDaily && active < limits.AIConcurrent && cases+2 <= limits.MaxCases
				report["room_for_two_attempts"] = room
				if !room {
					ok = false
				}
			}
			conn.Close(ctx)
		}
		cancel()
	}
	if report["quota_verified"] != true {
		report["quota_error"] = "Cannot read configured database/quota; no completion permitted by this check."
		ok = false
	}
	report["preflight_passed"] = ok
	raw, _ := json.MarshalIndent(report, "", "  ")
	raw = append(raw, '\n')
	if *output != "" {
		f, e := os.OpenFile(*output, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if e != nil {
			fmt.Fprintln(os.Stderr, "Cannot create new provider report.")
			os.Exit(1)
		}
		_, e = f.Write(raw)
		f.Close()
		if e != nil {
			fmt.Fprintln(os.Stderr, "Cannot write provider report.")
			os.Exit(1)
		}
	}
	fmt.Print(string(raw))
	if !ok {
		os.Exit(1)
	}
}
