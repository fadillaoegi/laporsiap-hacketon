package domain

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestSummaryBindingPreservesEvidenceAndRejectsProviderSourceChanges(t *testing.T) {
	sources := []ChronologySource{{ID: "M1", Text: "  I sent IDR 264,000.\r\nNo refund.\n"}, {ID: "M2", Text: "Perekrut menjanjikan komisi Rp63.000. 🌱\t"}}
	before := append([]ChronologySource(nil), sources...)
	good := `{"events":[{"source_id":"M1","description":"The sender reports sending IDR 264,000 without a refund."},{"source_id":"M2","description":"The sender reports a recruiter promise of IDR 63,000 commission."}]}`
	events, err := ParseChronologySummaries([]byte(good), sources)
	if err != nil {
		t.Fatal(err)
	}
	for i, event := range events {
		if event.SourceID != sources[i].ID || event.Evidence != sources[i].Text {
			t.Fatal("source bytes changed during binding")
		}
	}
	for name, raw := range map[string]string{
		"invented ID":          strings.Replace(good, "M1", "M9", 1),
		"duplicate ID":         strings.Replace(good, "M2", "M1", 1),
		"reordered ID":         strings.Replace(strings.Replace(good, "M1", "M3", 1), "M2", "M1", 1),
		"missing event":        `{"events":[{"source_id":"M1","description":"A claim."}]}`,
		"empty list":           `{"events":[]}`,
		"injected quote":       strings.Replace(good, `"description":`, `"evidence":"rewritten source","description":`, 1),
		"injected transaction": strings.TrimSuffix(good, "}") + `,"transfers":[]}`,
		"missing description":  `{"events":[{"source_id":"M1"},{"source_id":"M2","description":"A claim."}]}`,
		"null description":     strings.Replace(good, `"The sender reports sending IDR 264,000 without a refund."`, "null", 1),
		"overlong description": strings.Replace(good, "The sender reports sending IDR 264,000 without a refund.", strings.Repeat("word ", 36), 1),
		"trailing JSON":        good + `{}`,
		"markdown":             "```json\n" + good + "\n```",
	} {
		t.Run(name, func(t *testing.T) {
			if result, err := ParseChronologySummaries([]byte(raw), sources); err == nil || result != nil {
				t.Fatal("invalid provider response accepted")
			}
		})
	}
	if !reflect.DeepEqual(sources, before) {
		t.Fatal("binding changed original sources")
	}
	// Saved/event-edit contract still rejects altered evidence.
	events[0].Evidence = strings.TrimSpace(events[0].Evidence)
	raw, _ := json.Marshal(map[string]any{"events": events})
	if _, err := ParseChronology(raw, sources); err == nil {
		t.Fatal("saved draft quote validator was weakened")
	}
}

func TestGeneratedSummariesMustPreserveRelativeTime(t *testing.T) {
	for _, tc := range []struct {
		name, source, summary string
		invalid               bool
	}{
		{"observed English omission", "The recruiter promises IDR 63,000 commission next week; it has not been paid.\n", "Sender reports a recruiter promised IDR 63,000 commission which has not been paid.", true},
		{"corrected English", "The recruiter promises IDR 63,000 commission next week; it has not been paid.\n", "The sender reports the recruiter promises IDR 63,000 commission next week; it has not been paid.", false},
		{"Indonesian to English", "Perekrut menjanjikan komisi Rp63.000 pekan depan; belum dibayar.", "The sender reports the recruiter promises IDR 63,000 commission next week; it has not been paid.", false},
		{"English to Indonesian", "Commission is promised next week.", "Pengirim mengatakan komisi dijanjikan minggu depan.", false},
		{"changed tomorrow", "Commission is promised next week.", "The sender reports commission is promised tomorrow.", true},
		{"invented today", "No refund received.", "The sender reports no refund received today.", true},
		{"multiple timings", "I sent IDR 10,000 yesterday; commission is promised tomorrow.", "The sender reports IDR 10,000 sent yesterday and commission promised tomorrow.", false},
		{"multiple timings one omitted", "I sent IDR 10,000 yesterday; commission is promised tomorrow.", "The sender reports IDR 10,000 sent yesterday and commission promised.", true},
		{"case insensitive", "NEXT WEEK, commission is promised.", "The sender reports commission is promised next week.", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			raw, _ := json.Marshal(map[string]any{"events": []map[string]string{{"source_id": "M1", "description": tc.summary}}})
			events, err := ParseChronologySummaries(raw, []ChronologySource{{ID: "M1", Text: tc.source}})
			if (err != nil) != tc.invalid {
				t.Fatalf("time contract result err=%v", err)
			}
			if tc.invalid && (events != nil || err.Error() != "chronology relative time differs from source") {
				t.Fatal("incorrect time diagnostic or partial draft returned")
			}
			if !tc.invalid && (events[0].Description != tc.summary || events[0].Evidence != tc.source) {
				t.Fatal("validation changed source or summary")
			}
		})
	}
}
