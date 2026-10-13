package integration

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBilingualReportPreservesEvidenceAndUnknownAmounts(t *testing.T) {
	c := Case{Title: "Belum diketahui", Synthetic: true, Revision: 3,
		RawChat:   "Belum diketahui\nSaya baru mengirim Rp180.000.",
		Messages:  []Message{{ID: "M1", Text: "Belum diketahui"}, {ID: "M2", Text: "Saya baru mengirim Rp180.000."}},
		Transfers: []Transfer{{ID: "T1", Status: "successful", Note: "Dalam proses"}, {ID: "T2", Status: "pending", Amount: ptr(420000)}}}
	c.Draft = manualDraft(c)
	for _, language := range []string{"en", "id"} {
		report := ReportInLanguage(c, language)
		marker := "DATA SIMULASI"
		if language == "en" {
			marker = "SIMULATED DATA"
		}
		if !strings.Contains(report, marker) || strings.Contains(report, "DEMO SINTETIS") || strings.Contains(report, "DEMO DATA") {
			t.Fatal("simulation disclosure missing or outdated demo label remains")
		}
		for _, original := range []string{c.Title, c.Messages[0].Text, c.Messages[1].Text, c.Draft.Events[0].Description, c.Transfers[0].Note} {
			if !strings.Contains(report, original) {
				t.Fatalf("%s report changed original content %q", language, original)
			}
		}
		if language == "en" && (!strings.Contains(report, "(known amounts): Unknown rupiah") || !strings.Contains(report, "Sender: Unknown") || !strings.Contains(report, "Transfer 1: missing amount, transfer date, sending account, receiving account.")) {
			t.Fatal("English report lost unknown values or localized warnings")
		}
		if language == "id" && report != Report(c) {
			t.Fatal("Indonesian report compatibility changed")
		}
		data, err := ExportZIPInLanguage(c, nil, language)
		if err != nil {
			t.Fatal(err)
		}
		z, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			t.Fatal(err)
		}
		contents := map[string]string{}
		for _, f := range z.File {
			r, err := f.Open()
			if err != nil {
				t.Fatal(err)
			}
			b, err := io.ReadAll(r)
			r.Close()
			if err != nil {
				t.Fatal(err)
			}
			contents[f.Name] = string(b)
		}
		if contents["kronologi.txt"] != report || contents["percakapan-asli.txt"] != c.RawChat {
			t.Fatal("localized ZIP changed raw evidence or report")
		}
		var snapshot Case
		if err := json.Unmarshal([]byte(contents["kasus-direview.json"]), &snapshot); err != nil || snapshot.Draft.Events[0] != c.Draft.Events[0] {
			t.Fatal("localized ZIP changed draft snapshot")
		}
		if language == "en" && !strings.HasPrefix(contents["transfer.csv"], "id,sending_bank,sending_account,") {
			t.Fatal("English CSV headers missing")
		}
		if !strings.Contains(contents["transfer.csv"], ",,,successful,Dalam proses") {
			t.Fatal("CSV lost unknown amount or original user note")
		}
	}
}

func TestExportLanguageUsesRequestedLocale(t *testing.T) {
	for _, tc := range []struct{ header, want string }{{"en", "en"}, {"en-US,en;q=0.9", "en"}, {"id", "id"}, {"id-ID,en;q=0.9", "id"}, {"", "id"}, {"fr", "id"}} {
		r := httptest.NewRequest("GET", "/api/cases/example/export", nil)
		r.Header.Set("Accept-Language", tc.header)
		if got := exportLanguage(r); got != tc.want {
			t.Fatalf("%q locale: got %q want %q", tc.header, got, tc.want)
		}
	}
}
