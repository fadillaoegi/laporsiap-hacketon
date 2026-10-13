package integration

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestHTTPProviderReceivesPreviewAndLocalEvidenceRemainsOriginal(t *testing.T) {
	s := testStore(t)
	var wanted AIPreview
	calls := 0
	private := []string{"eleanor@example.test", "+6281234567890", "1234567890", "Jalan Mawar 7", "PRIVATE-TITLE", "PRIVATE-TRANSFER"}
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		for _, term := range private {
			if strings.Contains(string(body), term) {
				t.Errorf("fictional private field crossed HTTP boundary: %s", term)
			}
		}
		var payload struct {
			Messages []struct {
				Content string `json:"content"`
			} `json:"messages"`
		}
		if json.Unmarshal(body, &payload) != nil || len(payload.Messages) != 2 {
			t.Error("bad provider payload")
			w.WriteHeader(500)
			return
		}
		parts := strings.Split(payload.Messages[1].Content, "\nPesan sumber:\n")
		var received []ChronologySource
		if len(parts) != 2 || json.Unmarshal([]byte(parts[1]), &received) != nil {
			t.Error("missing provider sources")
			w.WriteHeader(500)
			return
		}
		expected, _ := json.Marshal(wanted.Sources)
		actual, _ := json.Marshal(received)
		if string(expected) != string(actual) {
			t.Error("provider received something other than preview")
		}
		content, _ := json.Marshal(map[string]any{"events": []map[string]string{{"source_id": "M1", "description": "The sender states no payment was made."}}})
		_ = json.NewEncoder(w).Encode(map[string]any{"model": ApertusModel, "choices": []map[string]any{{"message": map[string]string{"content": string(content)}, "finish_reason": "stop"}}})
	}))
	defer provider.Close()
	h := (&Server{Store: s, Generator: NewProvider(provider.URL, "fictional-provider-credential", ApertusModel)}).Handler()
	c := createCase(t, s, h, "PRIVATE-TITLE")
	input := Input{Title: "PRIVATE-TITLE", Synthetic: true, Parser: "whatsapp_bracket_dmy", RawChat: "[10/10/2026, 10:00] Eleanor: eleanor@example.test +6281234567890 account 1234567890 address Jalan Mawar 7. I have not paid.", Transfers: []Transfer{{ID: "T1", ToAccount: "PRIVATE-TRANSFER", Status: "pending", Amount: ptr(100000)}}}
	path := "/api/cases/" + c.Case.ID
	w := call(t, h, "PUT", path, c.Token, map[string]any{"revision": c.Case.Revision, "input": input})
	requireStatus(t, w, 200)
	updated := readCase(t, w)
	settings := PrivacyOptions{Mode: "redacted", Terms: []string{"Jalan Mawar 7"}}
	w = call(t, h, "POST", path+"/ai-preview", c.Token, map[string]any{"revision": updated.Case.Revision, "privacy": settings})
	requireStatus(t, w, 200)
	if json.Unmarshal(w.Body.Bytes(), &wanted) != nil {
		t.Fatal("preview unavailable")
	}
	w = call(t, h, "POST", path+"/draft", c.Token, map[string]any{"revision": updated.Case.Revision, "mode": "apertus", "privacy": settings, "preview_hash": wanted.Hash, "external_consent": true})
	requireStatus(t, w, 200)
	after := readCase(t, w)
	if calls != 1 || after.Case.RawChat != input.RawChat || after.Case.Draft.Events[0].Evidence != updated.Case.Messages[0].Text || after.Case.Transfers[0].ToAccount != "PRIVATE-TRANSFER" {
		t.Fatal("HTTP privacy did not preserve original local data")
	}
}

func TestLifecycleRotationDeletionAndExpiry(t *testing.T) {
	s := testStore(t)
	h := (&Server{Store: s}).Handler()
	a := createCase(t, s, h, "Lifecycle fictional case")
	b := createCase(t, s, h, "Unrelated fictional case")
	path := "/api/cases/" + a.Case.ID
	if a.Case.ExpiresAt == "" || a.Case.AccessExpiresAt == "" {
		t.Fatal("expiry missing")
	}
	originalToken := a.Token
	w := call(t, h, "POST", path+"/draft", a.Token, map[string]any{"revision": a.Case.Revision, "mode": "manual"})
	requireStatus(t, w, 200)
	a = readCase(t, w)
	a.Token = originalToken
	w = attachmentRequest(t, h, a, "fictional.pdf", validAttachmentPDF(t))
	requireStatus(t, w, 200)
	a = readCase(t, w)
	a.Token = originalToken
	a = reviewCase(t, h, a)
	requireStatus(t, call(t, h, "POST", path+"/access/rotate", b.Token, map[string]any{"revision": a.Case.Revision}), 404)
	w = call(t, h, "POST", path+"/access/rotate", a.Token, map[string]any{"revision": a.Case.Revision})
	requireStatus(t, w, 200)
	rotated := readCase(t, w)
	if rotated.Token == "" || rotated.Token == originalToken || rotated.Case.Revision != a.Case.Revision || rotated.Case.ReviewedRevision == nil || *rotated.Case.ReviewedRevision != a.Case.Revision || rotated.Case.ExpiresAt != a.Case.ExpiresAt {
		t.Fatal("rotation changed content/review/retention or failed revocation")
	}
	for _, endpoint := range []string{"", "/history", "/export?format=zip", "/attachments/" + a.Case.Attachments[0].ID} {
		requireStatus(t, call(t, h, "GET", path+endpoint, originalToken, nil), 404)
	}
	requireStatus(t, call(t, h, "GET", path+"/export", rotated.Token, nil), 200)
	requireStatus(t, call(t, h, "DELETE", path, rotated.Token, map[string]any{"revision": 1, "confirm_delete": true}), 409)
	requireStatus(t, call(t, h, "DELETE", path, rotated.Token, map[string]any{"revision": a.Case.Revision, "confirm_delete": false}), 400)
	lease, err := s.ReserveAI(context.Background(), a.Case.ID, rotated.Token, a.Case.Revision)
	if err != nil {
		t.Fatal(err)
	}
	s.FinishAI(lease)
	t.Cleanup(func() { _, _ = s.Pool.Exec(context.Background(), "DELETE FROM ai_attempts WHERE id=$1", lease) })
	requireStatus(t, call(t, h, "DELETE", path, rotated.Token, map[string]any{"revision": a.Case.Revision, "confirm_delete": true}), 200)
	for _, table := range []string{"cases", "case_history", "attachments"} {
		column := "case_id"
		if table == "cases" {
			column = "id"
		}
		var n int
		err := s.Pool.QueryRow(context.Background(), "SELECT count(*) FROM "+table+" WHERE "+column+"=$1", a.Case.ID).Scan(&n)
		if err != nil || n != 0 {
			t.Fatalf("cascade deletion failed for %s", table)
		}
	}
	var detached bool
	if err := s.Pool.QueryRow(context.Background(), "SELECT case_id IS NULL FROM ai_attempts WHERE id=$1", lease).Scan(&detached); err != nil || !detached {
		t.Fatal("anonymous usage log still linked to deleted case")
	}
	requireStatus(t, call(t, h, "GET", "/api/cases/"+b.Case.ID, b.Token, nil), 200)
	expired := createCase(t, s, h, "Expired synthetic")
	accessExpired := createCase(t, s, h, "Expired access synthetic")
	_, err = s.Pool.Exec(context.Background(), "UPDATE cases SET expires_at=now()-interval '1 second' WHERE id=$1", expired.Case.ID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.Pool.Exec(context.Background(), "UPDATE cases SET access_expires_at=now()-interval '1 second' WHERE id=$1", accessExpired.Case.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []caseResponse{expired, accessExpired} {
		requireStatus(t, call(t, h, "GET", "/api/cases/"+c.Case.ID, c.Token, nil), 404)
		if _, _, err := s.ExportSnapshot(context.Background(), c.Case.ID, c.Token); !errors.Is(err, ErrNotFound) {
			t.Fatal("expired snapshot accessible")
		}
	}
	if n, err := s.PurgeExpired(context.Background()); err != nil || n < 1 {
		t.Fatal("expiry cleanup failed")
	}
	var retained bool
	if err := s.Pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM cases WHERE id=$1)", accessExpired.Case.ID).Scan(&retained); err != nil || !retained {
		t.Fatal("access expiry deleted case before retention")
	}
	deadline, _ := time.Parse(time.RFC3339, b.Case.ExpiresAt)
	if time.Until(deadline) < 29*24*time.Hour {
		t.Fatal("wrong retention interval")
	}
}

func TestPrivacyMasksPayloadAndRestoresOriginalEvidence(t *testing.T) {
	source := "Nadia: email nadia@example.test, call +6281234567890, account 1234567890, NIK 1234567890123456. I sent IDR 180,000. Address Jalan Mawar 7."
	c := Case{ID: "fictional", Revision: 2, RawChat: source, Messages: []Message{{ID: "M1", Sender: "Nadia", Text: source}}, Transfers: []Transfer{{ID: "T1", ToAccount: "DO_NOT_SEND"}}, Title: "DO_NOT_SEND", Attachments: []Attachment{{Filename: "DO_NOT_SEND"}}}
	masked, preview, err := PrepareAIInput(c, &PrivacyOptions{Mode: "redacted", Terms: []string{"Jalan Mawar 7"}}, "en")
	if err != nil {
		t.Fatal(err)
	}
	for _, private := range []string{"Nadia", "nadia@example.test", "+6281234567890", "1234567890", "1234567890123456", "Jalan Mawar 7"} {
		if strings.Contains(masked.RawChat, private) {
			t.Fatalf("private term leaked: %s", private)
		}
	}
	if !strings.Contains(masked.RawChat, "IDR 180,000") || preview.Redactions != 6 || masked.Title != "" || len(masked.Transfers) != 0 || len(masked.Attachments) != 0 || masked.Messages[0].Sender != "" {
		t.Fatal("payload boundary or amount changed")
	}
	d := &Draft{Model: ApertusModel, Events: []ChronologyEvent{{SourceID: "M1", Evidence: masked.Messages[0].Text, Description: "The sender states they sent IDR 180,000."}}}
	if err := restoreEvidence(d, masked, c, preview); err != nil {
		t.Fatal(err)
	}
	if d.Events[0].Evidence != source || d.PrivacyMode != "redacted" || d.Redactions != 6 {
		t.Fatal("original evidence not restored")
	}
	d.Events[0].Evidence = "altered quote"
	if restoreEvidence(d, masked, c, preview) == nil {
		t.Fatal("altered masked quote accepted")
	}
	_, same, _ := PrepareAIInput(c, &PrivacyOptions{Mode: "redacted", Terms: []string{"Jalan Mawar 7"}}, "en")
	if same.Hash != preview.Hash {
		t.Fatal("preview not deterministic")
	}
	_, other, _ := PrepareAIInput(c, &PrivacyOptions{Mode: "redacted", Terms: []string{"Jalan Mawar 7"}}, "id")
	if other.Hash == preview.Hash {
		t.Fatal("language change did not invalidate preview")
	}
	original, _, _ := PrepareAIInput(c, &PrivacyOptions{Mode: "original"}, "en")
	if original.Messages[0].Text != source {
		t.Fatal("explicit original mode changed source")
	}
}

func TestPrivacyRejectsBadSettingsAndHandlesPlaceholderCollision(t *testing.T) {
	c := Case{ID: "fixture", RawChat: "[LS_PRIVATE_0001] call 081234567890", Messages: []Message{{ID: "M1", Text: "[LS_PRIVATE_0001] call 081234567890"}}}
	masked, p, err := PrepareAIInput(c, nil, "en")
	if err != nil {
		t.Fatal(err)
	}
	if p.Redactions != 1 || !strings.Contains(masked.RawChat, "[[LS_PRIVATE_") {
		t.Fatal("placeholder collision not handled")
	}
	for _, settings := range []PrivacyOptions{{Mode: "bad"}, {Mode: "redacted", Terms: []string{"x"}}, {Mode: "redacted", Terms: make([]string, 51)}} {
		if _, _, err := PrepareAIInput(c, &settings, "en"); err == nil {
			t.Fatal("bad privacy settings accepted")
		}
	}
}

func TestRateLimiterAndInvalidSafetyConfig(t *testing.T) {
	s := &Store{Limits: DefaultSafetyLimits()}
	s.Limits.RequestsPerMinute = 2
	h := (&Server{Store: s}).Handler()
	for i := 0; i < 2; i++ {
		requireStatus(t, call(t, h, "GET", "/api/policy", "", nil), 200)
	}
	w := call(t, h, "GET", "/api/policy", "", nil)
	requireStatus(t, w, 429)
	if w.Header().Get("Retry-After") == "" {
		t.Fatal("missing retry hint")
	}
	t.Setenv("AI_CONCURRENT_LIMIT", "0")
	if _, err := LoadSafetyLimits(); err == nil {
		t.Fatal("unbounded configuration accepted")
	}
}

func TestProviderDisclosureUsesConfiguredEndpointWithoutLeakingKey(t *testing.T) {
	for _, tc := range []struct{ base, label string }{
		{"https://api.publicai.co/v1", "Public AI"},
		{"http://127.0.0.1:9000", "Local Apertus server"},
		{"https://private.example.test/v1", "Configured Apertus provider"},
	} {
		h := (&Server{Store: &Store{Limits: DefaultSafetyLimits()}, Generator: NewProvider(tc.base, "DO_NOT_DISCLOSE", ApertusModel)}).Handler()
		w := call(t, h, "GET", "/api/policy", "", nil)
		requireStatus(t, w, 200)
		if !strings.Contains(w.Body.String(), tc.label) || strings.Contains(w.Body.String(), "DO_NOT_DISCLOSE") || strings.Contains(w.Body.String(), tc.base) {
			t.Fatal("incorrect provider label or private config disclosed")
		}
	}
}

func TestAIPreviewConsentQuotaConcurrencyAndRestart(t *testing.T) {
	s := testStore(t)
	s.Limits = DefaultSafetyLimits()
	s.Limits.AICaseDaily = 2
	s.Limits.AIConcurrent = 1
	s.Limits.AIGlobalDaily = 1000 // Other test cases share a test database.
	calls := 0
	gen := generatorFunc(func(_ context.Context, c Case, lang string) (*Draft, error) {
		calls++
		d := manualDraftInLanguage(c, lang)
		for i := range d.Events {
			d.Events[i].Description = "The sender makes a statement."
		}
		return d, nil
	})
	h := (&Server{Store: s, Generator: gen}).Handler()
	a := createCase(t, s, h, "Quota synthetic")
	p := "/api/cases/" + a.Case.ID
	settings := PrivacyOptions{Mode: "redacted"}
	w := call(t, h, "POST", p+"/ai-preview", a.Token, map[string]any{"revision": 1, "privacy": settings})
	requireStatus(t, w, 200)
	var preview AIPreview
	if err := json.Unmarshal(w.Body.Bytes(), &preview); err != nil {
		t.Fatal(err)
	}
	requireStatus(t, call(t, h, "POST", p+"/draft", a.Token, map[string]any{"revision": 1, "mode": "apertus", "privacy": settings, "preview_hash": "stale", "external_consent": true}), 409)
	if calls != 0 {
		t.Fatal("stale preview called AI")
	}
	lease, err := s.ReserveAI(context.Background(), a.Case.ID, a.Token, 1)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.ReserveAI(context.Background(), a.Case.ID, a.Token, 1); !errors.Is(err, ErrAIBusy) {
		t.Fatal("concurrency guard bypassed")
	}
	s.FinishAI(lease)
	body := map[string]any{"revision": 1, "mode": "apertus", "privacy": settings, "preview_hash": preview.Hash, "external_consent": true}
	requireStatus(t, call(t, h, "POST", p+"/draft", a.Token, body), 200)
	fresh, err := OpenStore(context.Background(), os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer fresh.Pool.Close()
	fresh.Limits = s.Limits
	if _, err = fresh.ReserveAI(context.Background(), a.Case.ID, a.Token, 2); !errors.Is(err, ErrAIQuota) {
		t.Fatal("restart reset persistent quota")
	}
	if calls != 1 {
		t.Fatal("unexpected retry")
	}
	requireStatus(t, call(t, h, "POST", p+"/draft", a.Token, map[string]any{"revision": 2, "mode": "manual"}), 200)
	t.Cleanup(func() {
		_, _ = s.Pool.Exec(context.Background(), "DELETE FROM ai_attempts WHERE case_id=$1", a.Case.ID)
	})
}

func TestConcurrentReservationsDoNotOverspend(t *testing.T) {
	s := testStore(t)
	s.Limits = DefaultSafetyLimits()
	s.Limits.AICaseDaily = 1
	s.Limits.AIGlobalDaily = 1000
	h := (&Server{Store: s}).Handler()
	a := createCase(t, s, h, "Concurrent synthetic")
	var wg sync.WaitGroup
	var mu sync.Mutex
	accepted := []string{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id, err := s.ReserveAI(context.Background(), a.Case.ID, a.Token, 1)
			if err == nil {
				mu.Lock()
				accepted = append(accepted, id)
				mu.Unlock()
			} else if !errors.Is(err, ErrAIQuota) {
				t.Errorf("unexpected reserve error: %v", err)
			}
		}()
	}
	wg.Wait()
	if len(accepted) != 1 {
		t.Fatalf("accepted %d requests", len(accepted))
	}
	s.FinishAI(accepted[0])
	t.Cleanup(func() {
		_, _ = s.Pool.Exec(context.Background(), "DELETE FROM ai_attempts WHERE case_id=$1", a.Case.ID)
	})
}

func TestStorageAndRevisionCapacityKeepsSavedCase(t *testing.T) {
	s := testStore(t)
	h := (&Server{Store: s}).Handler()
	a := createCase(t, s, h, "Capacity synthetic")
	var count int
	var used int64
	if err := s.Pool.QueryRow(context.Background(), `SELECT count(*) FROM cases`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	s.Limits = DefaultSafetyLimits()
	s.Limits.MaxCases = count
	requireStatus(t, call(t, h, "POST", "/api/cases", "", Input{Title: "Must reject", RawChat: "Synthetic message", Parser: "plain"}), 413)
	s.Limits.MaxCases = count + 10
	s.Limits.MaxRevisions = 1
	requireStatus(t, call(t, h, "POST", "/api/cases/"+a.Case.ID+"/draft", a.Token, map[string]any{"revision": 1, "mode": "manual"}), 413)
	if err := s.Pool.QueryRow(context.Background(), `SELECT (SELECT COALESCE(sum(octet_length(document::text)),0) FROM cases)+(SELECT COALESCE(sum(octet_length(document::text)),0) FROM case_history)+(SELECT COALESCE(sum(octet_length(data)),0) FROM attachments)`).Scan(&used); err != nil {
		t.Fatal(err)
	}
	s.Limits.MaxRevisions = 200
	s.Limits.StorageBytes = used
	requireStatus(t, call(t, h, "POST", "/api/cases/"+a.Case.ID+"/draft", a.Token, map[string]any{"revision": 1, "mode": "manual"}), 413)
	got := readCase(t, call(t, h, "GET", "/api/cases/"+a.Case.ID, a.Token, nil))
	if got.Case.Revision != 1 || got.Case.Draft != nil {
		t.Fatal("quota rejection changed case")
	}
	requireStatus(t, call(t, h, "DELETE", "/api/cases/"+a.Case.ID, a.Token, map[string]any{"revision": 1, "confirm_delete": true}), 200)
}

func TestGlobalQuotaSurvivesDeletion(t *testing.T) {
	s := testStore(t)
	s.Limits = DefaultSafetyLimits()
	h := (&Server{Store: s}).Handler()
	a := createCase(t, s, h, "Global quota A")
	b := createCase(t, s, h, "Global quota B")
	var used int
	if err := s.Pool.QueryRow(context.Background(), "SELECT count(*) FROM ai_attempts WHERE created_at>=now()-interval '24 hours'").Scan(&used); err != nil {
		t.Fatal(err)
	}
	s.Limits.AIGlobalDaily = used + 1
	lease, err := s.ReserveAI(context.Background(), a.Case.ID, a.Token, 1)
	if err != nil {
		t.Fatal(err)
	}
	s.FinishAI(lease)
	requireStatus(t, call(t, h, "DELETE", "/api/cases/"+a.Case.ID, a.Token, map[string]any{"revision": 1, "confirm_delete": true}), 200)
	if _, err = s.ReserveAI(context.Background(), b.Case.ID, b.Token, 1); !errors.Is(err, ErrAIQuota) {
		t.Fatal("case deletion bypassed global daily limit")
	}
	t.Cleanup(func() { _, _ = s.Pool.Exec(context.Background(), "DELETE FROM ai_attempts WHERE id=$1", lease) })
}
