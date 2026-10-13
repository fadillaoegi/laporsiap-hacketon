package integration

import (
	"encoding/json"
	"os"
	"testing"
)

// Keep the user-facing test pack usable with the actual MVP parser and totals.
// Semantic expectations in the pack are for human review, not truth scores.
func TestTestingPackInputsAndExpectedTotals(t *testing.T) {
	raw, err := os.ReadFile("../../../../docs/testing/cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var pack struct {
		Synthetic bool `json:"synthetic"`
		Cases     []struct {
			ID       string `json:"id"`
			TXTFile  string `json:"txt_file"`
			Input    Input  `json:"input"`
			Expected struct {
				MessageCount int            `json:"message_count"`
				Total        *int64         `json:"successful_known_total"`
				Unknown      int            `json:"successful_unknown_amount"`
				Counts       map[string]int `json:"status_counts"`
			} `json:"expected"`
		} `json:"cases"`
	}
	if err := json.Unmarshal(raw, &pack); err != nil {
		t.Fatal(err)
	}
	if !pack.Synthetic || len(pack.Cases) != 8 {
		t.Fatal("expected eight synthetic scenarios")
	}
	for _, scenario := range pack.Cases {
		t.Run(scenario.ID, func(t *testing.T) {
			input := scenario.Input
			if !input.Synthetic {
				t.Fatal("scenario must be labeled synthetic")
			}
			chat, err := os.ReadFile("../../../../docs/testing/" + scenario.TXTFile)
			if err != nil || string(chat) != input.RawChat {
				t.Fatal("TXT differs from scenario input", err)
			}
			messages, err := ValidateInput(input)
			if err != nil || len(messages) != scenario.Expected.MessageCount {
				t.Fatal("invalid input/message count", err)
			}
			summary := Summarize(Case{Messages: messages, Transfers: input.Transfers})
			if (summary.SuccessfulKnownTotal == nil) != (scenario.Expected.Total == nil) || (summary.SuccessfulKnownTotal != nil && *summary.SuccessfulKnownTotal != *scenario.Expected.Total) || summary.SuccessfulUnknownAmount != scenario.Expected.Unknown {
				t.Fatalf("unexpected total: %+v", summary)
			}
			for status, count := range scenario.Expected.Counts {
				if summary.StatusCounts[status] != count {
					t.Fatalf("%s count mismatch", status)
				}
			}
			if input.Parser == "plain" {
				for _, message := range messages {
					if message.Sender != "" || message.SentAt != "" {
						t.Fatal("plain text invented metadata")
					}
				}
			}
		})
	}
}
