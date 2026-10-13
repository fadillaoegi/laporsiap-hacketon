package application

import (
	"context"
	"laporsiap/internal/domain"
)

// CaseRepository owns atomic authenticated persistence and optimistic revisions.
// Implementations must never overwrite a newer revision or an expired token.
type CaseRepository interface {
	Create(context.Context, domain.Case, string) (domain.Case, error)
	Get(context.Context, string, string) (domain.Case, error)
	Save(context.Context, domain.Case, string, string, int64, *domain.AttachmentData, string) (domain.Case, error)
	Attachment(context.Context, string, string, string) (domain.AttachmentData, error)
	History(context.Context, string, string) ([]domain.HistoryEntry, error)
	ExportSnapshot(context.Context, string, string) (domain.Case, []domain.AttachmentData, error)
	DeleteCase(context.Context, string, string, int64) error
	RotateAccess(context.Context, string, string, int64) (domain.Case, string, error)
	ReserveAI(context.Context, string, string, int64) (string, error)
	FinishAI(string)
	PolicyLimits() domain.SafetyLimits
	Ping(context.Context) error
}
type Generator interface {
	Generate(context.Context, domain.Case, string) (*domain.Draft, error)
}
type ProviderInfo struct {
	Label         string `json:"provider_label"`
	PromptVersion string `json:"prompt_version"`
	PromptSHA256  string `json:"prompt_sha256"`
}
type ProviderMetadata interface{ Info() ProviderInfo }
type Artifact struct {
	Bytes                 []byte
	ContentType, Filename string
}
type ReportRenderer interface {
	Render(domain.Case, []domain.AttachmentData, string, string) (Artifact, error)
}

type Failure struct {
	Kind  string
	Cause error
}

func (e *Failure) Error() string           { return e.Cause.Error() }
func (e *Failure) Unwrap() error           { return e.Cause }
func failure(kind string, err error) error { return &Failure{Kind: kind, Cause: err} }
