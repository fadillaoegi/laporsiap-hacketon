package integration

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompiledFrontendHeadersAndBoundary(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("<main>LaporSiap</main>"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("do not serve"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(dir, "..", "outside"), filepath.Join(dir, "escape")); err != nil {
		t.Fatal(err)
	}
	h, err := FrontendHandler(dir, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) }))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		url    string
		status int
	}{{"http://127.0.0.1:8787/", 200}, {"http://127.0.0.1:8787/.env", 404}, {"http://127.0.0.1:8787/escape", 404}, {"http://evil.example/", 403}, {"http://127.0.0.1:8787/api/health", 204}} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", tc.url, nil))
		requireStatus(t, w, tc.status)
		if tc.status == 200 && (!strings.Contains(w.Header().Get("Content-Security-Policy"), "frame-ancestors 'none'") || w.Header().Get("X-Frame-Options") != "DENY") {
			t.Fatal("static security headers missing")
		}
	}
}
