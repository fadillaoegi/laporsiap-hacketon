package reports

import (
	"laporsiap/internal/application"
	"laporsiap/internal/domain"
)

type Renderer struct{}

func (Renderer) Render(c domain.Case, files []domain.AttachmentData, format, language string) (application.Artifact, error) {
	var bytes []byte
	var err error
	kind := "text/plain; charset=utf-8"
	switch format {
	case "zip":
		kind = "application/zip"
		bytes, err = ExportZIPInLanguage(c, files, language)
	case "pdf":
		kind = "application/pdf"
		bytes, err = ReportPDF(c, language)
	default:
		bytes = []byte(ReportInLanguage(c, language))
	}
	return application.Artifact{Bytes: bytes, ContentType: kind}, err
}

var _ application.ReportRenderer = Renderer{}
