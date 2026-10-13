package integration

import "testing"

func TestManualDraftLanguagePreservesOriginalSource(t *testing.T) {
	c := Case{Messages: []Message{{ID: "M1", Text: "Saya belum mengirim uang."}}}
	for _, tc := range []struct{ language, prefix string }{{"id", "Kutipan pesan: "}, {"en", "Message quote: "}} {
		draft := manualDraftInLanguage(c, tc.language)
		if draft.Origin != "manual" || len(draft.Events) != 1 {
			t.Fatalf("%s: manual origin or entry missing", tc.language)
		}
		e := draft.Events[0]
		if e.SourceID != "M1" || e.Evidence != c.Messages[0].Text || e.Description != tc.prefix+c.Messages[0].Text {
			t.Fatalf("%s: source evidence was translated or altered: %+v", tc.language, e)
		}
	}
}
