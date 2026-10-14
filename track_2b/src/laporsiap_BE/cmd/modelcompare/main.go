// modelcompare runs the frozen fictional fixture through the Apertus provider
// directly (no API server, no database) to compare model sizes under the same
// prompt, privacy masking, and contract validation the app uses. Each case gets
// at most one attempt per model; there are no retries or fallbacks.
package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"laporsiap/internal/domain"
	"laporsiap/internal/infrastructure/apertus"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
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

type Result struct {
	ID             string                   `json:"id"`
	Language       string                   `json:"language"`
	Category       string                   `json:"category"`
	Model          string                   `json:"requested_model"`
	ResponseModel  string                   `json:"response_model,omitempty"`
	StartedAt      string                   `json:"started_at"`
	Status         string                   `json:"status"`
	Error          string                   `json:"error,omitempty"`
	Rule           string                   `json:"validation_rule,omitempty"`
	Seconds        float64                  `json:"request_seconds"`
	Usage          *domain.TokenUsage       `json:"usage,omitempty"`
	ContractPassed bool                     `json:"contract_passed"`
	TotalsMatch    bool                     `json:"totals_match_expected"`
	Sources        []domain.Message         `json:"original_sources"`
	Events         []domain.ChronologyEvent `json:"events,omitempty"`
	SemanticReview string                   `json:"semantic_review_status"`
	Criteria       []string                 `json:"semantic_review_criteria"`
}

func main() {
	live := flag.Bool("live", false, "Make real provider requests (one per case per model)")
	fixturePath := flag.String("fixture", "../../docs/evaluation/bilingual-safety-v1.json", "Original fictional evaluation fixture")
	models := flag.String("models", "swiss-ai/apertus-v1.5-8b,swiss-ai/apertus-v1.5-70b", "Comma-separated Apertus models")
	outRoot := flag.String("out", "../../output/evaluation", "Output root")
	only := flag.String("cases", "", "Optional comma-separated case IDs to run (e.g. a rerun of provider-rejected cases)")
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
	base, key, _ := apertus.ConfigFromEnv()
	var providers []*apertus.Provider
	for _, m := range strings.Split(*models, ",") {
		p := apertus.NewProvider(base, key, strings.TrimSpace(m))
		if err := p.ValidateConfiguration(); err != nil {
			panic("Provider configuration invalid: " + err.Error())
		}
		providers = append(providers, p)
	}
	if len(providers) == 0 || len(providers) > 3 {
		panic("Use 1–3 Apertus models")
	}
	cases := make([]domain.Case, len(fixture.Items))
	for i, item := range fixture.Items {
		if !item.Input.Synthetic || (item.Language != "en" && item.Language != "id") {
			panic("Non-fictional/invalid language fixture")
		}
		messages, err := domain.ValidateInput(item.Input)
		if err != nil {
			panic("Invalid fixture input")
		}
		transfers := item.Input.Transfers
		if transfers == nil {
			transfers = []domain.Transfer{}
		}
		cases[i] = domain.Case{ID: item.ID, Title: item.Input.Title, Synthetic: true, RawChat: item.Input.RawChat, Parser: item.Input.Parser, Messages: messages, Transfers: transfers, Attachments: []domain.Attachment{}, Revision: 1}
		privacy := item.Privacy
		if _, _, err := domain.PrepareAIInput(cases[i], &privacy, item.Language); err != nil {
			panic("Fixture privacy preparation failed")
		}
	}
	if !*live {
		fmt.Printf("Validated %d fictional cases × %d models = %d planned requests. SHA256=%s. No API requests. Use -live to run.\n", len(cases), len(providers), len(cases)*len(providers), hash)
		return
	}
	out := filepath.Join(*outRoot, "model-compare-"+time.Now().UTC().Format("2006-01-02T15-04-05Z"))
	if err = os.MkdirAll(out, 0700); err != nil {
		panic("Cannot create evaluation output")
	}
	write := func(name string, value any) {
		data, e := json.MarshalIndent(value, "", "  ")
		if e != nil {
			panic(e)
		}
		if e = os.MkdirAll(filepath.Dir(filepath.Join(out, name)), 0700); e != nil {
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
	write("prompt-frozen.json", map[string]any{"version": apertus.PromptVersion, "sha256": apertus.PromptSHA256(), "prompt": json.RawMessage(apertus.PromptSnapshot())})
	safeName := regexp.MustCompile(`[^A-Za-z0-9.-]+`)
	results := map[string][]Result{}
	// Interleave models per case so provider load drift affects both equally.
	selected := map[string]bool{}
	for _, id := range strings.Split(*only, ",") {
		if id = strings.TrimSpace(id); id != "" {
			selected[id] = true
		}
	}
	for i, item := range fixture.Items {
		if len(selected) > 0 && !selected[item.ID] {
			continue
		}
		for _, p := range providers {
			r := run(p, item, cases[i])
			write(safeName.ReplaceAllString(p.Model, "_")+"/"+item.ID+".json", r)
			results[p.Model] = append(results[p.Model], r)
			fmt.Println("END", item.ID, p.Model, r.Status, fmt.Sprintf("%.1fs", r.Seconds))
		}
	}
	summary := map[string]any{}
	for _, p := range providers {
		passed, in, outTok, secs := 0, 0, 0, 0.0
		rules := map[string]int{}
		for _, r := range results[p.Model] {
			if r.ContractPassed {
				passed++
			} else {
				reason := r.Rule
				if reason == "" {
					reason = r.Error
				}
				rules[reason]++
			}
			if r.Usage != nil {
				in += r.Usage.Input
				outTok += r.Usage.Output
			}
			secs += r.Seconds
		}
		summary[p.Model] = map[string]any{"cases": len(results[p.Model]), "contract_passed": passed, "failure_reasons": rules, "input_tokens": in, "output_tokens": outTok, "mean_request_seconds": secs / float64(len(results[p.Model]))}
	}
	write("summary.json", map[string]any{"fixture_version": fixture.Version, "fixture_sha256": hash, "prompt_version": apertus.PromptVersion, "prompt_sha256": apertus.PromptSHA256(), "path": "direct provider call (PrepareAIInput → Generate → RestoreEvidence); no API server or database", "timeout_seconds": 40, "case_filter": *only, "retry_count": 0, "fallback_count": 0, "models": summary, "semantic_review": "pending_assistant_review", "limitations": []string{"Development fixture of 12 fictional cases; not a statistical benchmark.", "API model identity is metadata, not a weights audit.", "Request timing reflects a shared hosted endpoint at one point in time.", "No retries; provider failures remain in the denominator."}})
	fmt.Println("Artifacts:", out)
}

func run(p *apertus.Provider, item Item, c domain.Case) Result {
	r := Result{ID: item.ID, Language: item.Language, Category: item.Category, Model: p.Model, StartedAt: time.Now().UTC().Format(time.RFC3339), Status: "failed_or_rejected", Sources: c.Messages, SemanticReview: "not_scored", Criteria: item.Criteria}
	privacy := item.Privacy
	masked, preview, err := domain.PrepareAIInput(c, &privacy, item.Language)
	if err != nil {
		r.Error = "privacy preparation failed"
		return r
	}
	before, _ := json.Marshal(c)
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	start := time.Now()
	draft, err := p.Generate(ctx, masked, item.Language)
	r.Seconds = time.Since(start).Seconds()
	if err == nil {
		r.ResponseModel, r.Usage = draft.Model, draft.Usage
		err = domain.RestoreEvidence(draft, masked, c, preview)
	}
	var identity *domain.ModelIdentityError
	var validation *domain.DraftValidationError
	switch {
	case errors.As(err, &identity):
		r.ResponseModel, r.Error = identity.ResponseModel, "response model differs from requested model"
		return r
	case errors.As(err, &validation):
		r.ResponseModel, r.Rule, r.Error = validation.ResponseModel, validation.Rule, "contract validation failed"
		return r
	case err != nil:
		r.Error = err.Error()
		return r
	}
	after, _ := json.Marshal(c)
	summary := domain.Summarize(c)
	r.TotalsMatch = reflect.DeepEqual(summary.SuccessfulKnownTotal, item.ExpectedTotal) && summary.SuccessfulUnknownAmount == item.ExpectedUnknown
	if string(before) != string(after) || !r.TotalsMatch {
		r.Error = "original input or totals changed"
		return r
	}
	r.Events, r.ContractPassed, r.Status, r.SemanticReview = draft.Events, true, "contract_passed_review_pending", "pending_assistant_review"
	return r
}
