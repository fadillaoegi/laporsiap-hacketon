package integration

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestProviderVerifiesIdentityContractAndBoundary(t *testing.T) {
	c := Case{RawChat: "Saya belum membayar.", Messages: []Message{{ID: "M1", Text: "Saya belum membayar.", Sender: "PRIVATE-SENDER"}}, Transfers: []Transfer{{ID: "private", ToAccount: "PRIVATE-ACCOUNT", Amount: ptr(800)}}}
	good := `{"events":[{"source_id":"M1","description":"Pengirim menyatakan belum membayar."}]}`
	for _, tc := range []struct {
		name, model, content, finish string
		status                       int
		wantError                    bool
	}{
		{"valid", ApertusModel, good, "stop", 200, false},
		{"model mismatch", "other/model", good, "stop", 200, true},
		{"unsafe model metadata", "invalid\nPRIVATE-SENDER", good, "stop", 200, true},
		{"fake quote", ApertusModel, strings.Replace(good, `"description":`, `"evidence":"Saya sudah membayar.","description":`, 1), "stop", 200, true},
		{"schema injection", ApertusModel, strings.TrimSuffix(good, "}") + `,"transfers":[]}`, "stop", 200, true},
		{"truncated", ApertusModel, good, "length", 200, true},
		{"provider unavailable", ApertusModel, good, "stop", 504, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				if r.URL.Path != "/v1/chat/completions" || r.Header.Get("User-Agent") != "LaporSiap/0.1" {
					t.Error("provider integration path/header wrong")
				}
				b, _ := io.ReadAll(r.Body)
				if strings.Contains(string(b), "PRIVATE-") || strings.Contains(string(b), `"transfers"`) {
					t.Error("non-chat data crossed AI boundary")
				}
				w.WriteHeader(tc.status)
				_ = json.NewEncoder(w).Encode(map[string]any{"model": tc.model, "choices": []map[string]any{{"message": map[string]string{"content": tc.content}, "finish_reason": tc.finish}}})
			}))
			defer server.Close()
			p := NewProvider(server.URL, "test-credential", ApertusModel)
			d, err := p.Generate(context.Background(), c, "id")
			if (err != nil) != tc.wantError || calls != 1 {
				t.Fatalf("result error=%v calls=%d", err, calls)
			}
			if !tc.wantError && (d.Model != ApertusModel || d.Events[0].Evidence != c.Messages[0].Text) {
				t.Fatal("valid draft lost source")
			}
			if tc.model != ApertusModel {
				var identity *ModelIdentityError
				if !errors.As(err, &identity) || identity.RequestedModel != ApertusModel {
					t.Fatal("missing safe model diagnostic")
				}
				wantModel := tc.model
				if tc.name == "unsafe model metadata" {
					wantModel = "unrecognized"
				}
				if identity.ResponseModel != wantModel {
					t.Fatal("unsafe or incorrect model diagnostic")
				}
				w := httptest.NewRecorder()
				fail(w, 502, err)
				var payload struct {
					Diagnostic ModelIdentityError `json:"provider_diagnostic"`
				}
				if json.Unmarshal(w.Body.Bytes(), &payload) != nil || payload.Diagnostic.ResponseModel != wantModel || strings.Contains(w.Body.String(), "PRIVATE-") {
					t.Fatal("HTTP diagnostic omitted identity or exposed raw content")
				}
			}
			if tc.name == "fake quote" || tc.name == "schema injection" {
				var validation *DraftValidationError
				if !errors.As(err, &validation) || validation.ResponseModel != ApertusModel || validation.Rule == "" {
					t.Fatal("missing draft validation diagnostic")
				}
				w := httptest.NewRecorder()
				fail(w, 502, err)
				if !strings.Contains(w.Body.String(), `"validation_rule"`) || strings.Contains(w.Body.String(), "PRIVATE-") || strings.Contains(w.Body.String(), "Saya sudah membayar") {
					t.Fatal("validation diagnostic missing or exposed source content")
				}
			}
		})
	}
}

func TestProviderRejectsLostTimingWithoutRetryOrMutation(t *testing.T) {
	for _, language := range []string{"en", "id"} {
		t.Run(language, func(t *testing.T) {
			c := Case{RawChat: "The recruiter promises IDR 63,000 commission next week; it has not been paid.\n", Messages: []Message{{ID: "M1", Text: "The recruiter promises IDR 63,000 commission next week; it has not been paid.\n"}}}
			before, _ := json.Marshal(c)
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				_, _ = io.Copy(io.Discard, r.Body)
				description := "Sender reports a recruiter promised IDR 63,000 commission which has not been paid."
				if language == "id" {
					description = "Pengirim menyatakan perekrut menjanjikan komisi Rp63.000 yang belum dibayar."
				}
				content, _ := json.Marshal(map[string]any{"events": []map[string]string{{"source_id": "M1", "description": description}}})
				_ = json.NewEncoder(w).Encode(map[string]any{"model": ApertusModel, "choices": []map[string]any{{"message": map[string]string{"content": string(content)}, "finish_reason": "stop"}}})
			}))
			defer server.Close()
			draft, err := NewProvider(server.URL, "fictional-key", ApertusModel).Generate(context.Background(), c, language)
			var invalid *DraftValidationError
			if draft != nil || !errors.As(err, &invalid) || invalid.Rule != "chronology relative time differs from source" || calls != 1 {
				t.Fatal("time omission accepted or silently retried")
			}
			after, _ := json.Marshal(c)
			if string(after) != string(before) {
				t.Fatal("time omission changed source")
			}
		})
	}
}

func TestProviderTimeoutDoesNotRetryOrMutateCase(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		_, _ = io.Copy(io.Discard, r.Body)
		<-r.Context().Done()
	}))
	defer server.Close()
	p := NewProvider(server.URL, "test-credential", ApertusModel)
	if p.Client.Timeout != 40*time.Second {
		t.Fatal("production timeout changed")
	}
	p.Client.Timeout = time.Second // Same HTTP timeout path, bounded test duration.
	c := Case{RawChat: "No refund received.", Messages: []Message{{ID: "M1", Text: "No refund received."}}, Transfers: []Transfer{{ID: "T1", Status: "pending", Amount: ptr(100000)}}}
	before, _ := json.Marshal(c)
	draft, err := p.Generate(context.Background(), c, "en")
	if err == nil || draft != nil || calls.Load() != 1 {
		t.Fatalf("timeout produced draft/retry: draft=%v err=%v calls=%d", draft, err, calls.Load())
	}
	if err.Error() != "Apertus tidak merespons dalam batas waktu atau koneksi gagal. Coba lagi atau gunakan draft manual." {
		t.Fatal("timeout error did not explain recovery")
	}
	after, _ := json.Marshal(c)
	if string(before) != string(after) {
		t.Fatal("timeout changed original case")
	}
}
