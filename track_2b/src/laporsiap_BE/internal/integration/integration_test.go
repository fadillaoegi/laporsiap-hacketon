package integration

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func testStore(t *testing.T) *Store {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL required for PostgreSQL integration tests")
	}
	s, err := OpenStore(context.Background(), url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(s.Pool.Close)
	s.Limits = DefaultSafetyLimits()
	s.Limits.AIGlobalDaily = 1000 // Isolated DB includes other deterministic provider tests.
	return s
}
func call(t *testing.T, h http.Handler, method, path, token string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var data []byte
	if body != nil {
		data, _ = json.Marshal(body)
	}
	r := httptest.NewRequest(method, "http://localhost:8787"+path, bytes.NewReader(data))
	if token != "" {
		r.Header.Set("Authorization", "Bearer "+token)
	}
	if body != nil {
		r.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}
func requireStatus(t *testing.T, w *httptest.ResponseRecorder, status int) {
	t.Helper()
	if w.Code != status {
		t.Fatalf("expected %d got %d: %s", status, w.Code, w.Body.String())
	}
}

type caseResponse struct {
	Case    Case    `json:"case"`
	Summary Summary `json:"summary"`
	Token   string  `json:"token"`
}

func readCase(t *testing.T, w *httptest.ResponseRecorder) caseResponse {
	t.Helper()
	var c caseResponse
	if err := json.Unmarshal(w.Body.Bytes(), &c); err != nil {
		t.Fatal(err)
	}
	return c
}
func createCase(t *testing.T, s *Store, h http.Handler, title string) caseResponse {
	t.Helper()
	w := call(t, h, "POST", "/api/cases", "", Input{Title: title, Synthetic: true, RawChat: "Saya baru mengirim Rp180.000.\nTambah deposit Rp420.000.", Parser: "plain", Transfers: []Transfer{{ID: "T1", Status: "successful", Amount: ptr(180000), FromAccount: "DEMO-A", ToAccount: "DEMO-B"}, {ID: "T2", Status: "pending", Amount: ptr(420000)}}})
	requireStatus(t, w, 201)
	c := readCase(t, w)
	t.Cleanup(func() {
		_, _ = s.Pool.Exec(context.Background(), "DELETE FROM ai_attempts WHERE case_id=$1", c.Case.ID)
		_, err := s.Pool.Exec(context.Background(), "DELETE FROM cases WHERE id=$1", c.Case.ID)
		if err != nil {
			t.Error(err)
		}
	})
	return c
}
func reviewCase(t *testing.T, h http.Handler, c caseResponse) caseResponse {
	t.Helper()
	w := call(t, h, "POST", "/api/cases/"+c.Case.ID+"/review", c.Token, map[string]any{"revision": c.Case.Revision, "sources_reviewed": true, "transfers_reviewed": true, "limitations_accepted": true})
	requireStatus(t, w, 200)
	out := readCase(t, w)
	out.Token = c.Token
	return out
}
func TestPostgresAPIIsolationReviewInvalidationAttachmentsAndHistory(t *testing.T) {
	s := testStore(t)
	h := (&Server{Store: s}).Handler()
	a := createCase(t, s, h, "Kasus A")
	b := createCase(t, s, h, "Kasus B")
	path := "/api/cases/" + a.Case.ID
	for _, endpoint := range []string{"", "/export?format=zip", "/history"} {
		requireStatus(t, call(t, h, "GET", path+endpoint, b.Token, nil), 404)
	}
	requireStatus(t, call(t, h, "PUT", path, b.Token, map[string]any{"revision": 1, "input": toInputForTest(a.Case)}), 404)
	requireStatus(t, call(t, h, "GET", path+"/export", a.Token, nil), 409)
	requireStatus(t, call(t, h, "POST", path+"/draft", a.Token, map[string]any{"revision": 1, "mode": "apertus", "external_consent": false}), 400)
	w := call(t, h, "POST", path+"/draft", a.Token, map[string]any{"revision": 1, "mode": "manual", "external_consent": false})
	requireStatus(t, w, 200)
	token := a.Token
	a = readCase(t, w)
	a.Token = token
	a = reviewCase(t, h, a)
	requireStatus(t, call(t, h, "GET", path+"/export", a.Token, nil), 200)
	// Editing a reviewed draft invalidates review and preserves the original source.
	a.Case.Draft.Events[0].Description = "Pengirim menyatakan baru mengirim Rp180.000."
	w = call(t, h, "PUT", path+"/draft", a.Token, map[string]any{"revision": a.Case.Revision, "events": a.Case.Draft.Events})
	requireStatus(t, w, 200)
	a = readCase(t, w)
	a.Token = token
	if a.Case.ReviewedRevision != nil {
		t.Fatal("edit retained confirmation")
	}
	requireStatus(t, call(t, h, "GET", path+"/export", token, nil), 409)
	a = reviewCase(t, h, a)
	// A transfer edit must never infer a date from message time or include pending.
	input := toInputForTest(a.Case)
	input.Transfers = append(input.Transfers, Transfer{ID: "T3", Status: "successful"}, Transfer{ID: "T4", Status: "failed", Amount: ptr(900000)})
	oldRevision := a.Case.Revision
	w = call(t, h, "PUT", path, token, map[string]any{"revision": oldRevision, "input": input})
	requireStatus(t, w, 200)
	a = readCase(t, w)
	a.Token = token
	if a.Summary.SuccessfulKnownTotal == nil || *a.Summary.SuccessfulKnownTotal != 180000 || a.Summary.SuccessfulUnknownAmount != 1 || a.Case.Transfers[0].OccurredAt != "" {
		t.Fatalf("transfer invariant broken: %+v", a)
	}
	requireStatus(t, call(t, h, "POST", path+"/review", token, map[string]any{"revision": oldRevision, "sources_reviewed": true, "transfers_reviewed": true, "limitations_accepted": true}), 409)
	a = reviewCase(t, h, a)
	var data bytes.Buffer
	mw := multipart.NewWriter(&data)
	part, _ := mw.CreateFormFile("file", "../../bukti.pdf")
	_, _ = part.Write(validAttachmentPDF(t))
	_ = mw.Close()
	req := httptest.NewRequest("POST", "http://localhost:8787"+path+"/attachments", &data)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Case-Revision", stringRevision(a.Case.Revision))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	requireStatus(t, rec, 200)
	a = readCase(t, rec)
	a.Token = token
	if a.Case.ReviewedRevision != nil || len(a.Case.Attachments) != 1 || a.Case.Attachments[0].Filename != "bukti.pdf" {
		t.Fatal("attachment failed/invalidation missing")
	}
	fileID := a.Case.Attachments[0].ID
	requireStatus(t, call(t, h, "GET", path+"/attachments/"+fileID, b.Token, nil), 404)
	requireStatus(t, call(t, h, "GET", path+"/attachments/"+fileID, token, nil), 200)
	a = reviewCase(t, h, a)
	rec = call(t, h, "GET", path+"/export?format=zip", token, nil)
	requireStatus(t, rec, 200)
	z, err := zip.NewReader(bytes.NewReader(rec.Body.Bytes()), int64(rec.Body.Len()))
	if err != nil {
		t.Fatal(err)
	}
	if len(z.File) != 6 {
		t.Fatal("zip missing attachment", len(z.File))
	}
	rec = call(t, h, "GET", path+"/history", token, nil)
	requireStatus(t, rec, 200)
	if !strings.Contains(rec.Body.String(), "draft_edited") || !strings.Contains(rec.Body.String(), "human_review_confirmed") {
		t.Fatal("history missing audit trail")
	}
	req = httptest.NewRequest("DELETE", "http://localhost:8787"+path+"/attachments/"+fileID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Case-Revision", stringRevision(a.Case.Revision))
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	requireStatus(t, rec, 200)
	requireStatus(t, call(t, h, "GET", path+"/export", token, nil), 409)
	requireStatus(t, call(t, h, "GET", path+"/attachments/"+fileID, token, nil), 404)
}
func toInputForTest(c Case) Input {
	return Input{Title: c.Title, Synthetic: c.Synthetic, RawChat: c.RawChat, Parser: c.Parser, Transfers: c.Transfers}
}
func stringRevision(n int64) string { b, _ := json.Marshal(n); return string(b) }

type generatorFunc func(context.Context, Case, string) (*Draft, error)

func (f generatorFunc) Generate(ctx context.Context, c Case, language string) (*Draft, error) {
	return f(ctx, c, language)
}
func TestProviderFailureAndConcurrentInputChangeKeepCaseSafe(t *testing.T) {
	s := testStore(t)
	started := make(chan struct{})
	release := make(chan struct{})
	h := (&Server{Store: s, Generator: generatorFunc(func(_ context.Context, c Case, _ string) (*Draft, error) {
		close(started)
		<-release
		return manualDraft(c), nil
	})}).Handler()
	a := createCase(t, s, h, "Concurrency")
	path := "/api/cases/" + a.Case.ID
	done := make(chan *httptest.ResponseRecorder, 1)
	go func() {
		done <- call(t, h, "POST", path+"/draft", a.Token, map[string]any{"revision": 1, "mode": "apertus", "external_consent": true})
	}()
	<-started
	in := toInputForTest(a.Case)
	in.RawChat = "Sumber sudah diganti."
	requireStatus(t, call(t, h, "PUT", path, a.Token, map[string]any{"revision": 1, "input": in}), 200)
	close(release)
	requireStatus(t, <-done, 409)
	got := readCase(t, call(t, h, "GET", path, a.Token, nil))
	if got.Case.Draft != nil || got.Case.RawChat != in.RawChat {
		t.Fatal("stale AI overwrote newer input")
	}
	failing := (&Server{Store: s, Generator: generatorFunc(func(context.Context, Case, string) (*Draft, error) { return nil, errors.New("provider unavailable") })}).Handler()
	requireStatus(t, call(t, failing, "POST", path+"/draft", a.Token, map[string]any{"revision": got.Case.Revision, "mode": "apertus", "external_consent": true}), 502)
	requireStatus(t, call(t, failing, "POST", path+"/draft", a.Token, map[string]any{"revision": got.Case.Revision, "mode": "manual", "external_consent": false}), 200)
}
func TestRejectsUntrustedOriginAndHost(t *testing.T) {
	h := (&Server{}).Handler()
	for _, host := range []string{"evil.example:8787", "127.0.0.1.evil.example:8787"} {
		r := httptest.NewRequest("GET", "http://"+host+"/api/health", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		requireStatus(t, w, 403)
	}
	r := httptest.NewRequest("POST", "http://localhost:8787/api/cases", nil)
	r.Header.Set("Origin", "https://evil.example")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	requireStatus(t, w, 403)
}
