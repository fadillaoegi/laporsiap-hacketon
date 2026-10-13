// Bounded real-API evaluation. Tokens remain in memory; every case has at most
// one provider attempt. Contract results and semantic review are kept separate.
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"laporsiap/internal/domain"
	"laporsiap/internal/infrastructure/apertus"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"time"
)

type Item struct {
	ID              string                `json:"id"`
	Language        string                `json:"language"`
	Category        string                `json:"category"`
	Input           domain.Input          `json:"input"`
	Privacy         domain.PrivacyOptions `json:"privacy"`
	ExpectedTotal   *int64                `json:"expected_total"`
	ExpectedUnknown int                   `json:"expected_unknown_success"`
	Criteria        []string              `json:"semantic_review_criteria"`
}
type Envelope struct {
	Case    domain.Case    `json:"case"`
	Summary domain.Summary `json:"summary"`
	Token   string         `json:"token"`
}
type Diagnostic struct {
	RequestedModel string `json:"requested_model,omitempty"`
	ResponseModel  string `json:"response_model,omitempty"`
	Rule           string `json:"validation_rule,omitempty"`
}

func safeDiagnostic(data []byte) *Diagnostic {
	var body struct {
		Diagnostic *Diagnostic `json:"provider_diagnostic"`
	}
	if json.Unmarshal(data, &body) != nil || body.Diagnostic == nil {
		return nil
	}
	d := body.Diagnostic
	model := regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9./_:-]{0,199}$`)
	for _, value := range []string{d.RequestedModel, d.ResponseModel} {
		if value != "" && !model.MatchString(value) {
			return nil
		}
	}
	// Rules are server-owned strings; allow only fixed contract vocabulary.
	switch d.Rule {
	case "", "invalid chronology JSON schema", "trailing JSON after chronology", "chronology event count differs from sources", "chronology source order or identifier invalid", "chronology quote differs from original source", "chronology description empty or over 35 words", "chronology relative time differs from source":
	default:
		d.Rule = "contract_validation_failed"
	}
	return d
}

type Result struct {
	ProviderDiagnostic       *Diagnostic       `json:"provider_diagnostic,omitempty"`
	StateVerified            bool              `json:"state_verified"`
	StateError               string            `json:"state_error,omitempty"`
	OriginalSourcesUnchanged bool              `json:"original_sources_unchanged"`
	FailedDraftNotSaved      bool              `json:"failed_draft_not_saved"`
	ID                       string            `json:"id"`
	Language                 string            `json:"language"`
	Category                 string            `json:"category"`
	Status                   string            `json:"status"`
	Error                    string            `json:"error,omitempty"`
	APIStatus                int               `json:"api_status"`
	StartedAt                string            `json:"started_at"`
	CaseID                   string            `json:"case_id,omitempty"`
	Sources                  []domain.Message  `json:"original_sources,omitempty"`
	Preview                  *domain.AIPreview `json:"preview,omitempty"`
	Draft                    *domain.Draft     `json:"draft,omitempty"`
	ContractPassed           bool              `json:"contract_passed"`
	TransfersUnchanged       bool              `json:"transfers_unchanged"`
	RequestSeconds           float64           `json:"draft_request_seconds"`
	SemanticReview           string            `json:"semantic_review_status"`
	Criteria                 []string          `json:"semantic_review_criteria"`
}

func main() {
	live := flag.Bool("live", false, "One real AI attempt per fictional case, no retries")
	fixturePath := flag.String("fixture", "../../docs/evaluation/bilingual-safety-v1.json", "Original fictional evaluation fixture")
	apiURL := flag.String("api", "http://127.0.0.1:8787", "Local app origin only")
	flag.Parse()
	raw, err := os.ReadFile(*fixturePath)
	if err != nil {
		panic("Evaluation fixture unavailable")
	}
	var fixture struct {
		Version   string `json:"version"`
		Synthetic bool   `json:"synthetic"`
		Items     []Item `json:"items"`
	}
	if json.Unmarshal(raw, &fixture) != nil || !fixture.Synthetic || len(fixture.Items) == 0 || len(fixture.Items) > 12 {
		panic("Use 1–12 original fictional cases")
	}
	digest := sha256.Sum256(raw)
	hash := hex.EncodeToString(digest[:])
	for _, item := range fixture.Items {
		if !item.Input.Synthetic || (item.Language != "en" && item.Language != "id") {
			panic("Non-fictional/invalid language fixture")
		}
		if _, err := domain.ValidateInput(item.Input); err != nil {
			panic("Invalid fixture input")
		}
	}
	if !*live {
		fmt.Printf("Validated %d fictional cases. SHA256=%s. No API requests. Use -live to run.\n", len(fixture.Items), hash)
		return
	}
	u, err := url.Parse(*apiURL)
	if err != nil || u.Scheme != "http" || (u.Hostname() != "127.0.0.1" && u.Hostname() != "localhost") || u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.Path != "" {
		panic("Use a local app origin")
	}
	out := filepath.Join("../../output/evaluation", time.Now().UTC().Format("2006-01-02T15-04-05.000000000Z"))
	if err = os.MkdirAll(out, 0700); err != nil {
		panic("Cannot create evaluation output")
	}
	write := func(name string, value any) {
		data, e := json.MarshalIndent(value, "", "  ")
		if e != nil {
			panic(e)
		}
		f, e := os.OpenFile(filepath.Join(out, name), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if e != nil {
			panic(e)
		}
		defer f.Close()
		if _, e = f.Write(append(data, '\n')); e != nil {
			panic(e)
		}
	}
	write("fixture-frozen.json", map[string]any{"sha256": hash, "version": fixture.Version, "items": fixture.Items})
	client := &http.Client{Timeout: 50 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	send := func(endpoint, token, language string, body any) ([]byte, int, error) {
		data, _ := json.Marshal(body)
		req, e := http.NewRequest("POST", *apiURL+"/api"+endpoint, bytes.NewReader(data))
		if e != nil {
			return nil, 0, e
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept-Language", language)
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		res, e := client.Do(req)
		if e != nil {
			return nil, 0, fmt.Errorf("Local API unavailable")
		}
		defer res.Body.Close()
		data, e = io.ReadAll(io.LimitReader(res.Body, 512001))
		if e != nil || len(data) > 512000 {
			return nil, res.StatusCode, fmt.Errorf("Local response too large or unreadable")
		}
		if res.StatusCode >= 400 {
			var failure struct {
				Error string `json:"error"`
			}
			_ = json.Unmarshal(data, &failure)
			return data, res.StatusCode, fmt.Errorf("%s", failure.Error)
		}
		return data, res.StatusCode, nil
	}
	get := func(endpoint, token string) ([]byte, error) {
		req, e := http.NewRequest("GET", *apiURL+"/api"+endpoint, nil)
		if e != nil {
			return nil, fmt.Errorf("Invalid local request")
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		res, e := client.Do(req)
		if e != nil {
			return nil, fmt.Errorf("Local state check unavailable")
		}
		defer res.Body.Close()
		data, e := io.ReadAll(io.LimitReader(res.Body, 512001))
		if e != nil || len(data) > 512000 || res.StatusCode != 200 {
			return nil, fmt.Errorf("Local state check failed")
		}
		return data, nil
	}
	policyRaw, err := get("/policy", "")
	var policy struct {
		Version string `json:"prompt_version"`
		Hash    string `json:"prompt_sha256"`
	}
	if err != nil || json.Unmarshal(policyRaw, &policy) != nil || policy.Hash != apertus.PromptSHA256() || policy.Version != apertus.PromptVersion {
		panic("Running backend prompt differs from evaluator; no provider attempt made")
	}
	write("prompt-frozen.json", map[string]any{"version": apertus.PromptVersion, "sha256": apertus.PromptSHA256(), "prompt": json.RawMessage(apertus.PromptSnapshot()), "running_backend_verified": true})
	results := []Result{}
	for _, item := range fixture.Items {
		r := Result{ID: item.ID, Language: item.Language, Category: item.Category, StartedAt: time.Now().UTC().Format(time.RFC3339), Status: "not_started", SemanticReview: "not_scored", Criteria: item.Criteria}
		fmt.Println("START", item.ID)
		data, status, err := send("/cases", "", item.Language, item.Input)
		r.APIStatus = status
		var before Envelope
		if err == nil {
			err = json.Unmarshal(data, &before)
			if before.Token == "" {
				err = fmt.Errorf("Case access not returned")
			}
		}
		if err == nil {
			r.CaseID = before.Case.ID
			r.Sources = before.Case.Messages
			data, status, err = send("/cases/"+before.Case.ID+"/ai-preview", before.Token, item.Language, map[string]any{"revision": before.Case.Revision, "privacy": item.Privacy})
			r.APIStatus = status
		}
		var preview domain.AIPreview
		if err == nil {
			err = json.Unmarshal(data, &preview)
			r.Preview = &preview
		}
		if err == nil {
			start := time.Now()
			data, status, err = send("/cases/"+before.Case.ID+"/draft", before.Token, item.Language, map[string]any{"revision": before.Case.Revision, "mode": "apertus", "external_consent": true, "privacy": item.Privacy, "preview_hash": preview.Hash})
			r.APIStatus = status
			r.RequestSeconds = time.Since(start).Seconds()
			if err != nil {
				r.ProviderDiagnostic = safeDiagnostic(data)
			}
			var after Envelope
			if err == nil {
				err = json.Unmarshal(data, &after)
			}
			if err == nil {
				r.Draft = after.Case.Draft
				r.TransfersUnchanged = reflect.DeepEqual(before.Case.Transfers, after.Case.Transfers)
				if r.Draft == nil || r.Draft.Model != apertus.ApertusModel || r.Draft.Language != item.Language || r.Draft.PromptSHA256 != apertus.PromptSHA256() || r.Draft.PromptVersion != apertus.PromptVersion {
					err = fmt.Errorf("Draft model/language invalid")
				} else {
					sources := make([]domain.ChronologySource, len(before.Case.Messages))
					for i, m := range before.Case.Messages {
						sources[i] = domain.ChronologySource{ID: m.ID, Text: m.Text}
					}
					eventJSON, _ := json.Marshal(map[string]any{"events": r.Draft.Events})
					_, err = domain.ParseChronology(eventJSON, sources)
					if err == nil && (!r.TransfersUnchanged || before.Case.RawChat != after.Case.RawChat || !reflect.DeepEqual(before.Case.Messages, after.Case.Messages) || !reflect.DeepEqual(after.Summary.SuccessfulKnownTotal, item.ExpectedTotal) || after.Summary.SuccessfulUnknownAmount != item.ExpectedUnknown) {
						err = fmt.Errorf("Original input or totals changed")
					}
				}
			}
		}
		if before.Token != "" {
			state, stateErr := get("/cases/"+before.Case.ID, before.Token)
			var persisted Envelope
			if stateErr == nil {
				stateErr = json.Unmarshal(state, &persisted)
			}
			if stateErr == nil {
				r.TransfersUnchanged = reflect.DeepEqual(before.Case.Transfers, persisted.Case.Transfers)
				r.OriginalSourcesUnchanged = before.Case.RawChat == persisted.Case.RawChat && reflect.DeepEqual(before.Case.Messages, persisted.Case.Messages)
				r.FailedDraftNotSaved = persisted.Case.Revision == before.Case.Revision && persisted.Case.Draft == nil && persisted.Case.ReviewedRevision == nil
				if !r.TransfersUnchanged || !r.OriginalSourcesUnchanged || (r.APIStatus >= 400 && !r.FailedDraftNotSaved) {
					stateErr = fmt.Errorf("Persisted input or failed draft state changed")
				}
			}
			r.StateVerified = stateErr == nil
			if stateErr != nil {
				r.StateError = "Persisted state verification failed"
				if err == nil {
					err = stateErr
				}
			}
		}
		if err != nil {
			r.Status = "failed_or_rejected"
			r.Error = err.Error()
		} else {
			r.Status = "contract_passed_review_pending"
			r.ContractPassed = true
			r.SemanticReview = "pending_assistant_review"
		}
		// No response envelope/token is ever serialized.
		write(item.ID+".json", r)
		results = append(results, r)
		fmt.Println("END", item.ID, r.Status)
	}
	passed := 0
	for _, r := range results {
		if r.ContractPassed {
			passed++
		}
	}
	write("summary.json", map[string]any{"version": fixture.Version, "fixture_sha256": hash, "requested_model": apertus.ApertusModel, "prompt_version": apertus.PromptVersion, "prompt_sha256": apertus.PromptSHA256(), "cases": len(results), "contract_passed": passed, "results": results, "retry_count": 0, "fallback_count": 0, "independent_human_sessions": 0, "semantic_review": "Separate per-source decisions required; no automatic accuracy score.", "limitations": []string{"Development evaluation; not independent holdout.", "API model identity is metadata, not a weights audit.", "Request timing is not human work time or benefit.", "No retries; provider failures remain in the denominator."}})
	fmt.Printf("Evaluation artifacts: %s; contract %d/%d; semantic review pending\n", out, passed, len(results))
	if passed != len(results) {
		os.Exit(1)
	}
}
