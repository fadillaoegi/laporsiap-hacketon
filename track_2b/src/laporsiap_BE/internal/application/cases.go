package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"laporsiap/internal/domain"
	"path/filepath"
	"strings"
	"time"
)

type Service struct {
	Repository CaseRepository
	Generator  Generator
	Renderer   ReportRenderer
}

func NewService(repository CaseRepository, generator Generator, renderer ReportRenderer) *Service {
	return &Service{repository, generator, renderer}
}
func (s *Service) Get(ctx context.Context, id, token string) (domain.Case, error) {
	return s.Repository.Get(ctx, id, token)
}
func (s *Service) Create(ctx context.Context, in domain.Input) (domain.Case, string, error) {
	messages, err := domain.ValidateInput(in)
	if err != nil {
		return domain.Case{}, "", failure("input", err)
	}
	limits := s.Limits()
	created := time.Now().UTC()
	c := domain.Case{ID: domain.RandomID(), Title: in.Title, Synthetic: in.Synthetic, RawChat: in.RawChat, Parser: in.Parser, Messages: messages, Transfers: in.Transfers, Attachments: []domain.Attachment{}, Revision: 1, CreatedAt: created.Format(time.RFC3339), UpdatedAt: created.Format(time.RFC3339), ExpiresAt: created.Add(time.Duration(limits.RetentionDays) * 24 * time.Hour).Format(time.RFC3339), AccessExpiresAt: created.Add(time.Duration(limits.AccessHours) * time.Hour).Format(time.RFC3339)}
	if c.Transfers == nil {
		c.Transfers = []domain.Transfer{}
	}
	token := domain.RandomID()
	saved, err := s.Repository.Create(ctx, c, token)
	if err != nil {
		return saved, "", err
	}
	return saved, token, nil
}
func (s *Service) Limits() domain.SafetyLimits {
	if s.Repository == nil {
		return domain.DefaultSafetyLimits()
	}
	return s.Repository.PolicyLimits()
}
func (s *Service) Info() ProviderInfo {
	if p, ok := s.Generator.(ProviderMetadata); ok {
		return p.Info()
	}
	return ProviderInfo{Label: "Configured Apertus provider"}
}
func (s *Service) Health(ctx context.Context) error {
	if s.Repository == nil {
		return errors.New("Database belum tersedia.")
	}
	return s.Repository.Ping(ctx)
}
func (s *Service) current(ctx context.Context, id, token string, revision int64) (domain.Case, error) {
	c, err := s.Get(ctx, id, token)
	if err == nil && revision != c.Revision {
		err = domain.ErrConflict
	}
	return c, err
}
func (s *Service) Update(ctx context.Context, id, token string, revision int64, in domain.Input) (domain.Case, error) {
	c, err := s.current(ctx, id, token, revision)
	if err != nil {
		return c, err
	}
	msgs, err := domain.ValidateInput(in)
	if err != nil {
		return c, failure("input", err)
	}
	if c.RawChat != in.RawChat || c.Parser != in.Parser {
		c.Draft = nil
		for i := range c.Attachments {
			c.Attachments[i].SourceIDs = nil
		}
	}
	c.Title = in.Title
	c.RawChat = in.RawChat
	c.Parser = in.Parser
	c.Synthetic = in.Synthetic
	c.Transfers = in.Transfers
	c.Messages = msgs
	domain.PruneTransferLinks(&c)
	if c.Transfers == nil {
		c.Transfers = []domain.Transfer{}
	}
	return s.persist(ctx, c, token, "input_updated", false, nil, "")
}

type GenerateRequest struct {
	Revision        int64                  `json:"revision"`
	Mode            string                 `json:"mode"`
	ExternalConsent bool                   `json:"external_consent"`
	Privacy         *domain.PrivacyOptions `json:"privacy,omitempty"`
	PreviewHash     string                 `json:"preview_hash,omitempty"`
	Language        string                 `json:"-"`
}

func (s *Service) Generate(ctx context.Context, id, token string, req GenerateRequest) (domain.Case, error) {
	c, err := s.current(ctx, id, token, req.Revision)
	if err != nil {
		return c, err
	}
	switch req.Mode {
	case "manual":
		c.Draft = domain.ManualDraftInLanguage(c, req.Language)
	case "apertus":
		if !req.ExternalConsent {
			return c, failure("input", errors.New("Persetujuan pengiriman teks ke provider eksternal diperlukan."))
		}
		if s.Generator == nil {
			return c, failure("unavailable", errors.New("AI belum tersedia. Gunakan draft manual."))
		}
		masked, preview, err := domain.PrepareAIInput(c, req.Privacy, req.Language)
		if err != nil {
			return c, failure("input", err)
		}
		if req.Privacy != nil && req.PreviewHash != preview.Hash {
			return c, failure("conflict", errors.New("Preview AI berubah. Periksa teks dan setujui kembali."))
		}
		if len(c.Messages) > 20 || len(c.RawChat) > 12000 || len(masked.RawChat) > 12000 {
			return c, failure("input", errors.New("Draft AI dibatasi 20 pesan / 12 KB untuk menjaga kutipan lengkap. Gunakan draft manual untuk kasus lebih panjang."))
		}
		attempt, err := s.Repository.ReserveAI(ctx, c.ID, token, c.Revision)
		if err != nil {
			return c, err
		}
		defer s.Repository.FinishAI(attempt)
		bounded, cancel := context.WithTimeout(ctx, 40*time.Second)
		defer cancel()
		draft, err := s.Generator.Generate(bounded, masked, req.Language)
		if err != nil {
			return c, failure("provider", err)
		}
		if err = domain.RestoreEvidence(draft, masked, c, preview); err != nil {
			return c, failure("provider", err)
		}
		c.Draft = draft
	default:
		return c, failure("input", errors.New("Pilih Apertus atau manual."))
	}
	return s.persist(ctx, c, token, "draft_"+req.Mode, false, nil, "")
}
func (s *Service) Preview(ctx context.Context, id, token string, revision int64, privacy *domain.PrivacyOptions, language string) (domain.AIPreview, error) {
	c, err := s.current(ctx, id, token, revision)
	if err != nil {
		return domain.AIPreview{}, err
	}
	_, preview, err := domain.PrepareAIInput(c, privacy, language)
	if err != nil {
		err = failure("input", err)
	}
	return preview, err
}
func (s *Service) EditDraft(ctx context.Context, id, token string, revision int64, events []domain.ChronologyEvent) (domain.Case, error) {
	c, err := s.current(ctx, id, token, revision)
	if err != nil {
		return c, err
	}
	if c.Draft == nil {
		return c, failure("input", errors.New("Buat draft terlebih dahulu."))
	}
	if err = domain.ValidateEditedEvents(events, c); err != nil {
		return c, failure("input", err)
	}
	c.Draft.Events = events
	c.Draft.Checks = domain.ReviewChronology(events)
	if c.Draft.Origin == "apertus" {
		c.Draft.Origin = "apertus_edited"
	}
	return s.persist(ctx, c, token, "draft_edited", false, nil, "")
}
func (s *Service) Review(ctx context.Context, id, token string, revision int64, sources, transfers, limits bool) (domain.Case, error) {
	c, err := s.current(ctx, id, token, revision)
	if err != nil {
		return c, err
	}
	if c.Draft == nil || !sources || !transfers || !limits {
		return c, failure("input", errors.New("Periksa sumber, transfer, dan data belum lengkap sebelum mengonfirmasi review."))
	}
	if err = domain.ValidateEditedEvents(c.Draft.Events, c); err != nil {
		return c, failure("input", err)
	}
	return s.persist(ctx, c, token, "human_review_confirmed", true, nil, "")
}
func (s *Service) Delete(ctx context.Context, id, token string, revision int64, confirmed bool) error {
	if _, err := s.current(ctx, id, token, revision); err != nil {
		return err
	}
	if !confirmed {
		return failure("input", errors.New("Konfirmasi penghapusan kasus diperlukan."))
	}
	return s.Repository.DeleteCase(ctx, id, token, revision)
}
func (s *Service) Rotate(ctx context.Context, id, token string, revision int64) (domain.Case, string, error) {
	if _, err := s.current(ctx, id, token, revision); err != nil {
		return domain.Case{}, "", err
	}
	return s.Repository.RotateAccess(ctx, id, token, revision)
}
func (s *Service) LinkAttachment(ctx context.Context, id, token, attachment string, revision int64, sources, transfers []string) (domain.Case, error) {
	c, err := s.current(ctx, id, token, revision)
	if err != nil {
		return c, err
	}
	index := -1
	for i, a := range c.Attachments {
		if a.ID == attachment {
			index = i
			break
		}
	}
	if index < 0 {
		return c, domain.ErrNotFound
	}
	if err = domain.ValidateLinks(c, sources, transfers); err != nil {
		return c, failure("input", err)
	}
	c.Attachments[index].SourceIDs = sources
	c.Attachments[index].TransferIDs = transfers
	return s.persist(ctx, c, token, "attachment_links_updated", false, nil, "")
}
func SafeFilename(name string) string {
	name = filepath.Base(strings.ReplaceAll(name, "\\", "/"))
	var b strings.Builder
	for _, r := range name {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '.' || r == '-' || r == '_' {
			b.WriteRune(r)
		} else {
			b.WriteByte('_')
		}
	}
	out := strings.Trim(b.String(), ".")
	if out == "" {
		out = "bukti"
	}
	if len(out) > 100 {
		ext := filepath.Ext(out)
		if len(ext) > 10 {
			ext = ""
		}
		out = out[:80] + ext
	}
	return out
}
func (s *Service) Upload(ctx context.Context, id, token string, revision int64, name, contentType string, data []byte) (domain.Case, error) {
	c, err := s.current(ctx, id, token, revision)
	if err != nil {
		return c, err
	}
	if len(c.Attachments) >= 10 {
		return c, failure("input", errors.New("Maksimal 10 lampiran per kasus."))
	}
	if len(data) == 0 || len(data) > 5*1024*1024 {
		return c, failure("input", errors.New("Lampiran harus 1 byte–5 MB."))
	}
	ext := map[string]string{"image/png": ".png", "image/jpeg": ".jpg", "application/pdf": ".pdf"}[contentType]
	if ext == "" {
		return c, failure("input", errors.New("Bukti harus berupa PNG, JPEG, atau PDF. OCR belum tersedia."))
	}
	if err = domain.ValidateAttachment(data, contentType); err != nil {
		return c, failure("input", err)
	}
	name = SafeFilename(name)
	if !strings.EqualFold(filepath.Ext(name), ext) && !(ext == ".jpg" && strings.EqualFold(filepath.Ext(name), ".jpeg")) {
		name += ext
	}
	hash := sha256.Sum256(data)
	a := domain.AttachmentData{Meta: domain.Attachment{ID: domain.RandomID(), Filename: name, ContentType: contentType, Size: len(data), SHA256: hex.EncodeToString(hash[:])}, Bytes: data}
	c.Attachments = append(c.Attachments, a.Meta)
	return s.persist(ctx, c, token, "attachment_added", false, &a, "")
}
func (s *Service) DeleteAttachment(ctx context.Context, id, token, attachment string, revision int64) (domain.Case, error) {
	c, err := s.current(ctx, id, token, revision)
	if err != nil {
		return c, err
	}
	found := false
	for i, a := range c.Attachments {
		if a.ID == attachment {
			c.Attachments = append(c.Attachments[:i], c.Attachments[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return c, domain.ErrNotFound
	}
	return s.persist(ctx, c, token, "attachment_removed", false, nil, attachment)
}
func (s *Service) Download(ctx context.Context, id, token, attachment string) (domain.AttachmentData, error) {
	if _, err := s.Get(ctx, id, token); err != nil {
		return domain.AttachmentData{}, err
	}
	return s.Repository.Attachment(ctx, id, attachment, token)
}
func (s *Service) History(ctx context.Context, id, token string) ([]domain.HistoryEntry, error) {
	return s.Repository.History(ctx, id, token)
}
func (s *Service) Export(ctx context.Context, id, token, format, language string) (Artifact, error) {
	if _, err := s.Get(ctx, id, token); err != nil {
		return Artifact{}, err
	}
	c, files, err := s.Repository.ExportSnapshot(ctx, id, token)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return Artifact{}, err
		}
		return Artifact{}, failure("conflict", errors.New("Ekspor belum tersedia. Konfirmasikan review revisi terbaru."))
	}
	if c.Draft == nil || c.ReviewedRevision == nil || *c.ReviewedRevision != c.Revision {
		return Artifact{}, failure("conflict", errors.New("Ekspor belum tersedia. Konfirmasikan review revisi terbaru."))
	}
	if format == "" {
		format = "txt"
	}
	if format != "txt" && format != "pdf" && format != "zip" {
		return Artifact{}, failure("input", errors.New("Pilih format TXT, ZIP, atau PDF."))
	}
	artifact, err := s.Renderer.Render(c, files, format, language)
	if err != nil && format == "pdf" {
		return Artifact{}, failure("internal", errors.New("PDF gagal dibuat. TXT dan ZIP tetap tersedia."))
	}
	artifact.Filename = fmt.Sprintf("laporsiap-revisi-%d.%s", c.Revision, format)
	return artifact, err
}

// Review status belongs to the use case. The repository atomically persists the
// resulting aggregate against the explicit old revision.
func (s *Service) persist(ctx context.Context, c domain.Case, token, action string, confirmed bool, attachment *domain.AttachmentData, deleteID string) (domain.Case, error) {
	expected := c.Revision
	if expected >= s.Limits().MaxRevisions {
		return c, domain.ErrCapacity
	}
	c.Revision++
	c.UpdatedAt = domain.Now()
	domain.Invalidate(&c)
	if confirmed {
		revision := c.Revision
		c.ReviewedRevision = &revision
		c.ReviewedAt = domain.Now()
	}
	return s.Repository.Save(ctx, c, token, action, expected, attachment, deleteID)
}
