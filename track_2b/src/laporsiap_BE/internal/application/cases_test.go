package application

import (
	"context"
	"errors"
	"laporsiap/internal/domain"
	"testing"
)

// Embed the port for unused operations; exercised methods are explicitly fake.
// These tests execute without HTTP, PostgreSQL or provider infrastructure.
type memoryRepository struct {
	CaseRepository
	c                             domain.Case
	saves, reservations, finished int
}

func (r *memoryRepository) Get(context.Context, string, string) (domain.Case, error) { return r.c, nil }
func (r *memoryRepository) Save(_ context.Context, c domain.Case, _, _ string, _ int64, _ *domain.AttachmentData, _ string) (domain.Case, error) {
	r.saves++
	r.c = c
	return c, nil
}
func (r *memoryRepository) ReserveAI(context.Context, string, string, int64) (string, error) {
	r.reservations++
	return "one-attempt", nil
}
func (r *memoryRepository) FinishAI(string) { r.finished++ }

type generatorStub struct {
	calls      int
	fail       bool
	omitTiming bool
}

func (g *generatorStub) Generate(_ context.Context, c domain.Case, language string) (*domain.Draft, error) {
	g.calls++
	if g.fail {
		return nil, errors.New("provider unavailable")
	}
	draft := domain.ManualDraftInLanguage(c, language)
	if g.omitTiming {
		for i := range draft.Events {
			draft.Events[i].Description = "The sender reports unpaid commission."
		}
	}
	return draft, nil
}

func TestGeneratedTimingFailureDoesNotSaveOrInvalidateExistingDraft(t *testing.T) {
	original := domain.Case{ID: "case", Revision: 2, RawChat: "Commission is promised next week.", Messages: []domain.Message{{ID: "M1", Text: "Commission is promised next week."}}}
	original.Draft = domain.ManualDraftInLanguage(original, "en")
	repo := &memoryRepository{c: original}
	gen := &generatorStub{omitTiming: true}
	service := NewService(repo, gen, nil)
	_, err := service.Generate(context.Background(), "case", "token", GenerateRequest{Revision: 2, Mode: "apertus", ExternalConsent: true, Language: "en"})
	var invalid *domain.DraftValidationError
	if !errors.As(err, &invalid) || invalid.Rule != "chronology relative time differs from source" || repo.saves != 0 || repo.finished != 1 || gen.calls != 1 || repo.reservations != 1 {
		t.Fatal("lost timing saved/retried or lease not released")
	}
	if repo.c.Revision != 2 || repo.c.Draft != original.Draft || repo.c.RawChat != original.RawChat {
		t.Fatal("rejected new draft replaced existing draft or changed source")
	}
}
func TestPreviewConsentAndFailureProtectPersistence(t *testing.T) {
	original := domain.Case{ID: "case", Revision: 3, RawChat: "I sent IDR 264,000.", Messages: []domain.Message{{ID: "M1", Text: "I sent IDR 264,000."}}}
	repo := &memoryRepository{c: original}
	gen := &generatorStub{fail: true}
	service := NewService(repo, gen, nil)
	ctx := context.Background()
	privacy := &domain.PrivacyOptions{Mode: "redacted"}
	_, preview, err := domain.PrepareAIInput(original, privacy, "en")
	if err != nil {
		t.Fatal(err)
	}
	req := GenerateRequest{Revision: 3, Mode: "apertus", ExternalConsent: true, Privacy: privacy, PreviewHash: "stale", Language: "en"}
	if _, err = service.Generate(ctx, "case", "token", req); err == nil || repo.reservations != 0 || gen.calls != 0 {
		t.Fatal("stale preview crossed provider boundary")
	}
	req.PreviewHash = preview.Hash
	req.ExternalConsent = false
	if _, err = service.Generate(ctx, "case", "token", req); err == nil || repo.reservations != 0 {
		t.Fatal("missing consent reserved AI")
	}
	req.ExternalConsent = true
	if _, err = service.Generate(ctx, "case", "token", req); err == nil || repo.reservations != 1 || gen.calls != 1 || repo.finished != 1 || repo.saves != 0 || repo.c.Draft != nil || repo.c.RawChat != original.RawChat {
		t.Fatal("failed provider changed persisted case or retried")
	}
}
func TestUpdateInvalidatesDraftAndEvidenceLinks(t *testing.T) {
	repo := &memoryRepository{c: domain.Case{ID: "case", Revision: 2, RawChat: "old", Parser: "plain", Draft: &domain.Draft{}, Attachments: []domain.Attachment{{ID: "A1", SourceIDs: []string{"M1"}, TransferIDs: []string{"T1"}}}}}
	service := NewService(repo, nil, nil)
	saved, err := service.Update(context.Background(), "case", "token", 2, domain.Input{Title: "Fictional", RawChat: "new", Parser: "plain", Transfers: []domain.Transfer{}})
	if err != nil || saved.Draft != nil || len(saved.Attachments[0].SourceIDs) != 0 || len(saved.Attachments[0].TransferIDs) != 0 || repo.saves != 1 {
		t.Fatal("source/transfer changes retained stale AI/evidence links")
	}
	if _, err = service.Update(context.Background(), "case", "token", 1, domain.Input{}); !errors.Is(err, domain.ErrConflict) || repo.saves != 1 {
		t.Fatal("stale update bypassed revision guard")
	}
}

func (r *memoryRepository) PolicyLimits() domain.SafetyLimits { return domain.DefaultSafetyLimits() }
