package apertus

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCatalogCheckIsReadOnlyPinnedBoundedAndDoesNotFollowRedirects(t *testing.T) {
	for _, tc := range []struct {
		name, body string
		status     int
		wantError  bool
	}{
		{"listed", `{"data":[{"id":"swiss-ai/apertus-v1.5-70b"}]}`, 200, false},
		{"wrong model", `{"data":[{"id":"other/model"}]}`, 200, true},
		{"unavailable", `{"private":"DO_NOT_DISCLOSE"}`, 503, true},
		{"malformed", `{"data":`, 200, true},
		{"oversized", strings.Repeat("X", 256001), 200, true},
		{"redirect", ``, 302, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				if r.Method != "GET" || r.URL.Path != "/v1/models" || r.Header.Get("Authorization") != "Bearer DO_NOT_DISCLOSE" || r.Header.Get("User-Agent") != "LaporSiap/0.1" {
					t.Error("unexpected catalog request")
				}
				w.Header().Set("Location", "/redirect-target")
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer server.Close()
			err := NewProvider(server.URL+"/v1/", "DO_NOT_DISCLOSE", ApertusModel).CheckAvailability(context.Background())
			if (err != nil) != tc.wantError || calls != 1 {
				t.Fatalf("catalog result err=%v calls=%d", err, calls)
			}
			if err != nil && (strings.Contains(err.Error(), "DO_NOT_DISCLOSE") || strings.Contains(err.Error(), server.URL)) {
				t.Fatal("private config exposed")
			}
		})
	}
}

func TestInvalidConfigurationNeverMakesProviderRequest(t *testing.T) {
	for _, tc := range []struct{ base, key, model string }{
		{"https://example.test/v1", "", ApertusModel},
		{"https://example.test/v1", "key", "other/model"},
		{"http://example.test/v1", "key", ApertusModel},
		{"http://8.8.8.8/v1", "key", ApertusModel},
		{"https://example.test/v1", "key", "meta-llama/Llama-3.1-70B"},
		{"https://user:DO_NOT_DISCLOSE@example.test/v1", "key", ApertusModel},
		{"https://example.test/v1?token=DO_NOT_DISCLOSE", "key", ApertusModel},
	} {
		p := NewProvider(tc.base, tc.key, tc.model)
		if err := p.ValidateConfiguration(); err == nil || strings.Contains(err.Error(), "DO_NOT_DISCLOSE") {
			t.Fatal("invalid configuration accepted or exposed")
		}
	}
	for _, tc := range []struct{ base, model string }{
		{"http://inference:9000", "swiss-ai/Apertus-v1.5-70B"},
		{"http://10.0.0.5:8000/v1", "swiss-ai/apertus-v1.5-8b"},
		{"http://127.0.0.1:9000", ApertusModel},
		{"https://example.test/v1", "Apertus-v1.5-8B"},
	} {
		if err := NewProvider(tc.base, "key", tc.model).ValidateConfiguration(); err != nil {
			t.Fatalf("sovereign/Apertus configuration %q %q rejected: %v", tc.base, tc.model, err)
		}
	}
	var snapshot map[string]any
	if json.Unmarshal(PromptSnapshot(), &snapshot) != nil || snapshot["max_tokens"] != float64(2400) {
		t.Fatal("prompt invalid")
	}
}
