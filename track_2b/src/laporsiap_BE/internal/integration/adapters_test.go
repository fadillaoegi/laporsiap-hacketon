package integration

// Compatibility names are test-only. Production composition is in cmd/api;
// all preserved regressions exercise the new domain/use cases/adapters.
import (
	"laporsiap/internal/application"
	"laporsiap/internal/domain"
	"laporsiap/internal/infrastructure/apertus"
	"laporsiap/internal/infrastructure/config"
	"laporsiap/internal/infrastructure/postgres"
	"laporsiap/internal/infrastructure/reports"
	"laporsiap/internal/transport/httpapi"
	"net/http"
)

type Case = domain.Case
type Input = domain.Input
type Summary = domain.Summary
type Message = domain.Message
type Transfer = domain.Transfer
type Attachment = domain.Attachment
type Draft = domain.Draft
type ChronologySource = domain.ChronologySource
type ChronologyEvent = domain.ChronologyEvent
type ChronologyReview = domain.ChronologyReview
type PrivacyOptions = domain.PrivacyOptions
type AIPreview = domain.AIPreview
type attachmentData = domain.AttachmentData
type Store = postgres.Store
type Provider = apertus.Provider
type ModelIdentityError = domain.ModelIdentityError
type DraftValidationError = domain.DraftValidationError

var OpenStore = postgres.OpenStore
var ParseMessages = domain.ParseMessages
var ValidateInput = domain.ValidateInput
var Summarize = domain.Summarize
var ParseChronology = domain.ParseChronology
var ReviewChronology = domain.ReviewChronology
var PrepareAIInput = domain.PrepareAIInput
var restoreEvidence = domain.RestoreEvidence
var manualDraft = domain.ManualDraft
var manualDraftInLanguage = domain.ManualDraftInLanguage
var validateAttachment = domain.ValidateAttachment
var NewProvider = apertus.NewProvider
var PromptSnapshot = apertus.PromptSnapshot
var PromptSHA256 = apertus.PromptSHA256
var DefaultSafetyLimits = domain.DefaultSafetyLimits
var LoadSafetyLimits = config.LoadSafetyLimits
var ExportZIP = reports.ExportZIP
var ExportZIPInLanguage = reports.ExportZIPInLanguage
var Report = reports.Report
var ReportInLanguage = reports.ReportInLanguage
var ReportPDF = reports.ReportPDF
var pdfText = reports.PDFText
var safeFilename = application.SafeFilename
var fail = httpapi.ErrorResponse
var exportLanguage = httpapi.ExportLanguage
var FrontendHandler = httpapi.FrontendHandler
var ErrNotFound = domain.ErrNotFound
var ErrConflict = domain.ErrConflict
var ErrCapacity = domain.ErrCapacity
var ErrAIQuota = domain.ErrAIQuota
var ErrAIBusy = domain.ErrAIBusy

const ApertusModel = apertus.ApertusModel
const PromptVersion = apertus.PromptVersion

type Server struct {
	Store     *postgres.Store
	Generator application.Generator
}

func (s *Server) Handler() http.Handler {
	return (&httpapi.Server{Service: application.NewService(s.Store, s.Generator, reports.Renderer{})}).Handler()
}
