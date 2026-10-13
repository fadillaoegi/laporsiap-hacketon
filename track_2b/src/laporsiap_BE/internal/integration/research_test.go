package integration

import (
	"encoding/json"
	"os"
	"testing"
)

// A development regression, not a fresh independent accuracy evaluation.
func TestGuardRetainsObservedResearchFlags(t *testing.T) {
	raw, err := os.ReadFile("testdata/research-recovery-answers.json")
	if err != nil {
		t.Fatal(err)
	}
	var results []struct {
		ID       string `json:"id"`
		Verified bool   `json:"model_verified"`
		Answer   string `json:"answer"`
	}
	if json.Unmarshal(raw, &results) != nil || len(results) != 8 {
		t.Fatal("historical fixture incomplete")
	}
	flagged := map[string][]ChronologyReview{}
	for _, r := range results {
		if !r.Verified {
			t.Fatal("historical response model not verified")
		}
		var response struct {
			Events []ChronologyEvent `json:"events"`
		}
		if err = json.Unmarshal([]byte(r.Answer), &response); err != nil {
			t.Fatal(err)
		}
		flagged[r.ID] = ReviewChronology(response.Events)
	}
	if len(flagged["N01"]) != 1 || flagged["N01"][0].SourceID != "M3" || len(flagged["N05"]) != 2 || len(flagged["N04"]) == 0 {
		t.Fatalf("observed risk regression: %+v", flagged)
	}
}
