package httpapi

import (
	"errors"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

// A compiled frontend can share the API origin; the development proxy is not
// required in the private container. Never serve the repository or .env.
func FrontendHandler(directory string, api http.Handler) (http.Handler, error) {
	root, err := os.OpenRoot(directory)
	if err != nil {
		return nil, err
	}
	// The root descriptor lives for the server lifetime.
	if _, err := root.Stat("index.html"); err != nil {
		root.Close()
		return nil, errors.New("compiled frontend index.html not found")
	}
	files := http.FileServer(http.FS(root.FS()))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			api.ServeHTTP(w, r)
			return
		}
		if !allowedLocalRequest(r) {
			http.Error(w, "Origin/host not allowed", 403)
			return
		}
		if r.Method != "GET" && r.Method != "HEAD" {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "Method not allowed", 405)
			return
		}
		w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; style-src 'self'; img-src 'self' data:; font-src 'self'; connect-src 'self'; base-uri 'none'; form-action 'none'; frame-ancestors 'none'; object-src 'none'")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Cache-Control", "no-store")
		name := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if name == "" || name == "." {
			name = "index.html"
		}
		for _, part := range strings.Split(name, "/") {
			if strings.HasPrefix(part, ".") {
				http.NotFound(w, r)
				return
			}
		}
		info, err := fs.Stat(root.FS(), name)
		if err != nil || info.IsDir() {
			http.NotFound(w, r)
			return
		}
		files.ServeHTTP(w, r)
	}), nil
}
func allowedLocalRequest(r *http.Request) bool {
	allowed := map[string]bool{"127.0.0.1:8787": true, "localhost:8787": true, "127.0.0.1:5175": true, "localhost:5175": true}
	if configured := os.Getenv("APP_ORIGIN"); configured != "" {
		u, err := url.Parse(configured)
		if err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != "" && u.User == nil && u.RawQuery == "" && u.Fragment == "" && u.Path == "" {
			allowed[u.Host] = true
		}
	}
	if !allowed[r.Host] {
		return false
	}
	if origin := r.Header.Get("Origin"); origin != "" {
		u, err := url.Parse(origin)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" || !allowed[u.Host] {
			return false
		}
	}
	return true
}
