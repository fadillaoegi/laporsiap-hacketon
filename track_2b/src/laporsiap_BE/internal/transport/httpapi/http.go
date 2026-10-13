package httpapi

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"laporsiap/internal/application"
	"laporsiap/internal/domain"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	Service *application.Service
}

func jsonResponse(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func ErrorResponse(w http.ResponseWriter, status int, err error) {
	var identity *domain.ModelIdentityError
	if errors.As(err, &identity) {
		jsonResponse(w, status, map[string]any{"error": err.Error(), "provider_diagnostic": identity})
		return
	}
	var validation *domain.DraftValidationError
	if errors.As(err, &validation) {
		jsonResponse(w, status, map[string]any{"error": err.Error(), "provider_diagnostic": validation})
		return
	}
	jsonResponse(w, status, map[string]string{"error": err.Error()})
}
func decode(w http.ResponseWriter, r *http.Request, out any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 256000)
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(out); err != nil {
		return errors.New("Input JSON tidak valid, memiliki field tambahan, atau terlalu besar.")
	}
	var extra any
	if d.Decode(&extra) != io.EOF {
		return errors.New("Input harus satu objek JSON.")
	}
	return nil
}
func respondCase(w http.ResponseWriter, c domain.Case) {
	jsonResponse(w, 200, map[string]any{"case": c, "summary": domain.Summarize(c)})
}
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if s.Service.Health(ctx) != nil {
			ErrorResponse(w, 503, errors.New("Database belum tersedia."))
			return
		}
		jsonResponse(w, 200, map[string]string{"status": "ok", "database": "postgresql"})
	})
	mux.HandleFunc("POST /api/cases", s.create)
	mux.HandleFunc("GET /api/cases/{id}", s.get)
	mux.HandleFunc("PUT /api/cases/{id}", s.update)
	mux.HandleFunc("DELETE /api/cases/{id}", s.deleteCase)
	mux.HandleFunc("POST /api/cases/{id}/access/rotate", s.rotateAccess)
	mux.HandleFunc("GET /api/policy", func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, struct {
			domain.SafetyLimits
			application.ProviderInfo
		}{s.Service.Limits(), s.Service.Info()})
	})
	mux.HandleFunc("POST /api/cases/{id}/draft", s.generate)
	mux.HandleFunc("POST /api/cases/{id}/ai-preview", s.previewAI)
	mux.HandleFunc("PUT /api/cases/{id}/draft", s.editDraft)
	mux.HandleFunc("POST /api/cases/{id}/review", s.review)
	mux.HandleFunc("POST /api/cases/{id}/attachments", s.upload)
	mux.HandleFunc("GET /api/cases/{id}/attachments/{attachment}", s.download)
	mux.HandleFunc("PUT /api/cases/{id}/attachments/{attachment}", s.linkAttachment)
	mux.HandleFunc("DELETE /api/cases/{id}/attachments/{attachment}", s.deleteAttachment)
	mux.HandleFunc("GET /api/cases/{id}/export", s.export)
	mux.HandleFunc("GET /api/cases/{id}/history", s.history)
	limiter := &localLimiter{buckets: map[string]rateBucket{}, limit: s.Service.Limits().RequestsPerMinute}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		if !allowedLocalRequest(r) {
			ErrorResponse(w, 403, errors.New("Host lokal tidak diizinkan."))
			return
		}
		if r.URL.Path != "/api/health" && !limiter.allow(r.RemoteAddr) {
			tooMany(w, errors.New("Terlalu banyak permintaan. Tunggu satu menit lalu coba lagi."))
			return
		}
		mux.ServeHTTP(w, r)
	})
}
func bearer(r *http.Request) string {
	s := r.Header.Get("Authorization")
	if !strings.HasPrefix(s, "Bearer ") {
		return ""
	}
	token := strings.TrimPrefix(s, "Bearer ")
	if len(token) != 48 {
		return ""
	}
	if _, err := hex.DecodeString(token); err != nil {
		return ""
	}
	return token
}
func (s *Server) authorized(w http.ResponseWriter, r *http.Request) (domain.Case, string, bool) {
	token := bearer(r)
	if token == "" {
		ErrorResponse(w, 404, domain.ErrNotFound)
		return domain.Case{}, "", false
	}
	c, err := s.Service.Get(r.Context(), r.PathValue("id"), token)
	if err != nil {
		s.storeError(w, err)
		return c, "", false
	}
	return c, token, true
}
func (s *Server) storeError(w http.ResponseWriter, err error) {
	var f *application.Failure
	switch {
	case errors.Is(err, domain.ErrNotFound):
		ErrorResponse(w, 404, err)
	case errors.Is(err, domain.ErrConflict):
		ErrorResponse(w, 409, err)
	case errors.Is(err, domain.ErrCapacity):
		ErrorResponse(w, 413, err)
	case errors.Is(err, domain.ErrAIQuota) || errors.Is(err, domain.ErrAIBusy):
		tooMany(w, err)
	case errors.As(err, &f):
		status := map[string]int{"input": 400, "conflict": 409, "provider": 502, "unavailable": 503, "internal": 500}[f.Kind]
		if status == 0 {
			status = 500
		}
		ErrorResponse(w, status, err)
	default:
		ErrorResponse(w, 500, errors.New("Penyimpanan gagal. Coba kembali; jangan ekspor sebelum review tersimpan."))
	}
}
func checkRevision(w http.ResponseWriter, c domain.Case, revision int64) bool {
	if revision != c.Revision {
		ErrorResponse(w, 409, domain.ErrConflict)
		return false
	}
	return true
}

func (s *Server) result(w http.ResponseWriter, c domain.Case, err error) {
	if err != nil {
		s.storeError(w, err)
		return
	}
	respondCase(w, c)
}
func (s *Server) create(w http.ResponseWriter, r *http.Request) {
	var in domain.Input
	if err := decode(w, r, &in); err != nil {
		ErrorResponse(w, 400, err)
		return
	}
	c, token, err := s.Service.Create(r.Context(), in)
	if err != nil {
		s.storeError(w, err)
		return
	}
	jsonResponse(w, 201, map[string]any{"case": c, "summary": domain.Summarize(c), "token": token})
}
func (s *Server) get(w http.ResponseWriter, r *http.Request) {
	c, _, ok := s.authorized(w, r)
	if ok {
		respondCase(w, c)
	}
}
func (s *Server) update(w http.ResponseWriter, r *http.Request) {
	c, token, ok := s.authorized(w, r)
	if !ok {
		return
	}
	var req struct {
		Revision int64        `json:"revision"`
		Input    domain.Input `json:"input"`
	}
	if err := decode(w, r, &req); err != nil {
		ErrorResponse(w, 400, err)
		return
	}
	saved, err := s.Service.Update(r.Context(), c.ID, token, req.Revision, req.Input)
	s.result(w, saved, err)
}
func (s *Server) generate(w http.ResponseWriter, r *http.Request) {
	c, token, ok := s.authorized(w, r)
	if !ok {
		return
	}
	var req application.GenerateRequest
	if err := decode(w, r, &req); err != nil {
		ErrorResponse(w, 400, err)
		return
	}
	req.Language = ExportLanguage(r)
	saved, err := s.Service.Generate(r.Context(), c.ID, token, req)
	s.result(w, saved, err)
}
func (s *Server) previewAI(w http.ResponseWriter, r *http.Request) {
	c, token, ok := s.authorized(w, r)
	if !ok {
		return
	}
	var req struct {
		Revision int64                  `json:"revision"`
		Privacy  *domain.PrivacyOptions `json:"privacy"`
	}
	if err := decode(w, r, &req); err != nil {
		ErrorResponse(w, 400, err)
		return
	}
	preview, err := s.Service.Preview(r.Context(), c.ID, token, req.Revision, req.Privacy, ExportLanguage(r))
	if err != nil {
		s.storeError(w, err)
		return
	}
	jsonResponse(w, 200, preview)
}
func (s *Server) editDraft(w http.ResponseWriter, r *http.Request) {
	c, token, ok := s.authorized(w, r)
	if !ok {
		return
	}
	var req struct {
		Revision int64                    `json:"revision"`
		Events   []domain.ChronologyEvent `json:"events"`
	}
	if err := decode(w, r, &req); err != nil {
		ErrorResponse(w, 400, err)
		return
	}
	saved, err := s.Service.EditDraft(r.Context(), c.ID, token, req.Revision, req.Events)
	s.result(w, saved, err)
}
func (s *Server) review(w http.ResponseWriter, r *http.Request) {
	c, token, ok := s.authorized(w, r)
	if !ok {
		return
	}
	var req struct {
		Revision    int64 `json:"revision"`
		Sources     bool  `json:"sources_reviewed"`
		Transfers   bool  `json:"transfers_reviewed"`
		Limitations bool  `json:"limitations_accepted"`
	}
	if err := decode(w, r, &req); err != nil {
		ErrorResponse(w, 400, err)
		return
	}
	saved, err := s.Service.Review(r.Context(), c.ID, token, req.Revision, req.Sources, req.Transfers, req.Limitations)
	s.result(w, saved, err)
}
func (s *Server) deleteCase(w http.ResponseWriter, r *http.Request) {
	c, token, ok := s.authorized(w, r)
	if !ok {
		return
	}
	var req struct {
		Revision int64 `json:"revision"`
		Confirm  bool  `json:"confirm_delete"`
	}
	if err := decode(w, r, &req); err != nil {
		ErrorResponse(w, 400, err)
		return
	}
	if err := s.Service.Delete(r.Context(), c.ID, token, req.Revision, req.Confirm); err != nil {
		s.storeError(w, err)
		return
	}
	jsonResponse(w, 200, map[string]bool{"deleted": true})
}
func (s *Server) rotateAccess(w http.ResponseWriter, r *http.Request) {
	c, token, ok := s.authorized(w, r)
	if !ok {
		return
	}
	var req struct {
		Revision int64 `json:"revision"`
	}
	if err := decode(w, r, &req); err != nil {
		ErrorResponse(w, 400, err)
		return
	}
	updated, next, err := s.Service.Rotate(r.Context(), c.ID, token, req.Revision)
	if err != nil {
		s.storeError(w, err)
		return
	}
	jsonResponse(w, 200, map[string]any{"case": updated, "summary": domain.Summarize(updated), "token": next})
}
func (s *Server) linkAttachment(w http.ResponseWriter, r *http.Request) {
	c, token, ok := s.authorized(w, r)
	if !ok {
		return
	}
	var req struct {
		Revision    int64    `json:"revision"`
		SourceIDs   []string `json:"source_ids"`
		TransferIDs []string `json:"transfer_ids"`
	}
	if err := decode(w, r, &req); err != nil {
		ErrorResponse(w, 400, err)
		return
	}
	saved, err := s.Service.LinkAttachment(r.Context(), c.ID, token, r.PathValue("attachment"), req.Revision, req.SourceIDs, req.TransferIDs)
	s.result(w, saved, err)
}
func (s *Server) upload(w http.ResponseWriter, r *http.Request) {
	c, token, ok := s.authorized(w, r)
	if !ok {
		return
	}
	revision, _ := strconv.ParseInt(r.Header.Get("X-Case-Revision"), 10, 64)
	if !checkRevision(w, c, revision) {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 5*1024*1024+65536)
	reader, err := r.MultipartReader()
	if err != nil {
		ErrorResponse(w, 400, errors.New("Unggah sebagai multipart/form-data."))
		return
	}
	part, err := reader.NextPart()
	if err != nil || part.FormName() != "file" || part.FileName() == "" {
		ErrorResponse(w, 400, errors.New("Satu file bukti diperlukan."))
		return
	}
	data, err := io.ReadAll(io.LimitReader(part, 5*1024*1024+1))
	if err != nil || len(data) == 0 || len(data) > 5*1024*1024 {
		ErrorResponse(w, 400, errors.New("Lampiran harus 1 byte–5 MB."))
		return
	}
	if _, err = reader.NextPart(); err != io.EOF {
		ErrorResponse(w, 400, errors.New("Unggah satu file pada setiap permintaan."))
		return
	}
	saved, err := s.Service.Upload(r.Context(), c.ID, token, revision, part.FileName(), http.DetectContentType(data), data)
	s.result(w, saved, err)
}

func (s *Server) download(w http.ResponseWriter, r *http.Request) {
	c, token, ok := s.authorized(w, r)
	if !ok {
		return
	}
	a, err := s.Service.Download(r.Context(), c.ID, token, r.PathValue("attachment"))
	if err != nil {
		s.storeError(w, err)
		return
	}
	w.Header().Set("Content-Type", a.Meta.ContentType)
	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": a.Meta.Filename}))
	_, _ = w.Write(a.Bytes)
}
func (s *Server) deleteAttachment(w http.ResponseWriter, r *http.Request) {
	c, token, ok := s.authorized(w, r)
	if !ok {
		return
	}
	revision, _ := strconv.ParseInt(r.Header.Get("X-Case-Revision"), 10, 64)
	saved, err := s.Service.DeleteAttachment(r.Context(), c.ID, token, r.PathValue("attachment"), revision)
	s.result(w, saved, err)
}
func (s *Server) export(w http.ResponseWriter, r *http.Request) {
	c, token, ok := s.authorized(w, r)
	if !ok {
		return
	}
	artifact, err := s.Service.Export(r.Context(), c.ID, token, r.URL.Query().Get("format"), ExportLanguage(r))
	if err != nil {
		s.storeError(w, err)
		return
	}
	w.Header().Set("Content-Type", artifact.ContentType)
	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": artifact.Filename}))
	_, _ = w.Write(artifact.Bytes)
}
func (s *Server) history(w http.ResponseWriter, r *http.Request) {
	c, token, ok := s.authorized(w, r)
	if !ok {
		return
	}
	history, err := s.Service.History(r.Context(), c.ID, token)
	if err != nil {
		s.storeError(w, err)
		return
	}
	jsonResponse(w, 200, map[string]any{"history": history})
}

// The UI sends its selected language. Requests without a language retain the
// existing Indonesian export format.
func ExportLanguage(r *http.Request) string {
	if strings.HasPrefix(strings.ToLower(r.Header.Get("Accept-Language")), "en") {
		return "en"
	}
	return "id"
}
