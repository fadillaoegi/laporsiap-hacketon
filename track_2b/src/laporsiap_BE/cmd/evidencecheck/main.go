// Read-only checks of fictional evaluation state and a whitelisted provider
// catalog entry. Never prints access hashes, credentials or raw HTTP bodies.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"laporsiap/internal/domain"
	"laporsiap/internal/infrastructure/apertus"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	out := filepath.Join("../../output/hardening", "state-check-"+time.Now().UTC().Format("2006-01-02T15-04-05.000000000Z"))
	if os.MkdirAll(out, 0700) != nil {
		panic("Cannot create output")
	}
	write := func(name string, v any) {
		b, _ := json.MarshalIndent(v, "", "  ")
		f, e := os.OpenFile(filepath.Join(out, name), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if e != nil {
			panic("Cannot save evidence")
		}
		defer f.Close()
		if _, e = f.Write(append(b, '\n')); e != nil {
			panic("Cannot save evidence")
		}
	}
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		panic("Database unavailable")
	}
	defer pool.Close()
	fixtures, err := filepath.Glob("../../output/evaluation/*/fixture-frozen.json")
	if err != nil {
		panic("Fixtures unavailable")
	}
	rows := []map[string]any{}
	failed := false
	for _, filename := range fixtures {
		var frozen struct {
			Items []struct {
				ID    string       `json:"id"`
				Input domain.Input `json:"input"`
			}
		}
		b, e := os.ReadFile(filename)
		if e != nil || json.Unmarshal(b, &frozen) != nil {
			panic("Invalid evidence fixture")
		}
		for _, item := range frozen.Items {
			if !item.Input.Synthetic {
				panic("Fictional inputs only")
			}
			var result struct {
				CaseID   string `json:"case_id"`
				Contract bool   `json:"contract_passed"`
				Status   string `json:"status"`
			}
			b, e = os.ReadFile(filepath.Join(filepath.Dir(filename), item.ID+".json"))
			if e != nil || json.Unmarshal(b, &result) != nil {
				continue
			}
			if result.CaseID == "" {
				continue
			}
			var document []byte
			e = pool.QueryRow(ctx, "SELECT document FROM cases WHERE id=$1", result.CaseID).Scan(&document)
			var c domain.Case
			if e != nil || json.Unmarshal(document, &c) != nil {
				rows = append(rows, map[string]any{"id": item.ID, "status": "state_unavailable"})
				failed = true
				continue
			}
			messages, e := domain.ParseMessages(item.Input.RawChat, item.Input.Parser)
			unchanged := e == nil && c.RawChat == item.Input.RawChat && reflect.DeepEqual(messages, c.Messages) && reflect.DeepEqual(c.Transfers, item.Input.Transfers)
			failureUntouched := result.Contract || (c.Revision == 1 && c.Draft == nil && c.ReviewedRevision == nil)
			verified := unchanged && failureUntouched
			if !verified {
				failed = true
			}
			rows = append(rows, map[string]any{"id": item.ID, "generation_status": result.Status, "original_chat_messages_transfers_unchanged": unchanged, "failed_generation_has_no_draft_or_review": !result.Contract && failureUntouched, "revision": c.Revision, "state_verified": verified})
		}
	}
	write("evaluation-state.json", map[string]any{"checked_at": time.Now().UTC().Format(time.RFC3339), "read_only": true, "cases": rows, "all_verified": !failed})
	base, key, _ := apertus.ConfigFromEnv()
	u, e := url.Parse(base)
	catalog := map[string]any{"checked_at": time.Now().UTC().Format(time.RFC3339), "model": apertus.ApertusModel, "status": "not_checked", "live_generation_calls": 0}
	if e == nil && u.Scheme == "https" && u.Host != "" && u.User == nil && u.RawQuery == "" && u.Fragment == "" && key != "" {
		endpoint := strings.TrimRight(base, "/")
		if !strings.HasSuffix(endpoint, "/v1") {
			endpoint += "/v1"
		}
		endpoint += "/models"
		req, _ := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
		req.Header.Set("Authorization", "Bearer "+key)
		client := &http.Client{Timeout: 10 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
		res, e := client.Do(req)
		if e == nil {
			defer res.Body.Close()
			catalog["http_status"] = res.StatusCode
			if res.StatusCode == 200 {
				b, e := io.ReadAll(io.LimitReader(res.Body, 1000001))
				var data struct {
					Data []struct {
						ID      string `json:"id"`
						Pricing *struct {
							Input  float64 `json:"input"`
							Output float64 `json:"output"`
						} `json:"pricing"`
					} `json:"data"`
				}
				if e == nil && len(b) <= 1000000 && json.Unmarshal(b, &data) == nil {
					for _, m := range data.Data {
						if m.ID == apertus.ApertusModel {
							catalog["status"] = "listed"
							catalog["pricing"] = m.Pricing
						}
					}
				}
			}
		}
	}
	write("provider-catalog.json", catalog)
	fmt.Printf("Read-only evidence saved: %s; checked %d cases; all_verified=%t; no generation requests\n", out, len(rows), !failed)
	if failed {
		os.Exit(1)
	}
}
