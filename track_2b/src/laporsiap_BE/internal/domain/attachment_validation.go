package domain

import (
	"bytes"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"regexp"
	"strconv"
)

const maxAttachmentImagePixels = 24_000_000

var (
	pdfHeader     = regexp.MustCompile(`^%PDF-[12]\.[0-9](?:\r|\n)`)
	pdfXrefObject = regexp.MustCompile(`^[0-9]+\s+[0-9]+\s+obj\b`)
	pdfXrefType   = regexp.MustCompile(`/Type\s*/XRef\b`)
)

// Images are decoded after a pixel bound to catch damaged payloads without
// allocating an arbitrarily large bitmap. PDF checks are an envelope check,
// not a full PDF parser, malware scan, OCR, or verification of evidence.
func ValidateAttachment(data []byte, contentType string) error {
	switch contentType {
	case "image/png", "image/jpeg":
		config, format, err := image.DecodeConfig(bytes.NewReader(data))
		if err != nil || (contentType == "image/png" && format != "png") || (contentType == "image/jpeg" && format != "jpeg") {
			return errors.New("Gambar rusak atau tidak lengkap. Pilih file PNG/JPEG yang dapat dibuka.")
		}
		if config.Width <= 0 || config.Height <= 0 || int64(config.Width)*int64(config.Height) > maxAttachmentImagePixels {
			return errors.New("Dimensi gambar terlalu besar. Gunakan gambar maksimal 24 juta piksel.")
		}
		if _, _, err := image.Decode(bytes.NewReader(data)); err != nil {
			return errors.New("Gambar rusak atau tidak lengkap. Pilih file PNG/JPEG yang dapat dibuka.")
		}
	case "application/pdf":
		if !CompletePDFEnvelope(data) {
			return errors.New("PDF tidak lengkap atau struktur dasarnya tidak valid. Pilih PDF yang dapat dibuka.")
		}
	}
	return nil
}

func CompletePDFEnvelope(data []byte) bool {
	if !pdfHeader.Match(data) {
		return false
	}
	end := bytes.TrimRight(data, " \t\r\n\f\x00")
	if !bytes.HasSuffix(end, []byte("%%EOF")) {
		return false
	}
	position := bytes.LastIndex(end, []byte("startxref"))
	if position < 0 || position+len("startxref") >= len(end)-len("%%EOF") {
		return false
	}
	offset, err := strconv.ParseInt(string(bytes.TrimSpace(end[position+len("startxref"):len(end)-len("%%EOF")])), 10, 64)
	if err != nil || offset < 0 || offset >= int64(position) {
		return false
	}
	xref := end[offset:position]
	if bytes.HasPrefix(xref, []byte("xref\n")) || bytes.HasPrefix(xref, []byte("xref\r")) || bytes.HasPrefix(xref, []byte("xref ")) || bytes.HasPrefix(xref, []byte("xref\t")) {
		return bytes.Contains(xref, []byte("trailer")) && bytes.Contains(xref, []byte("/Root"))
	}
	// PDF 1.5+ can point directly to an xref-stream object instead of a table.
	stream := bytes.Index(xref, []byte("stream"))
	return stream > 0 && pdfXrefObject.Match(xref) && pdfXrefType.Match(xref[:stream]) && bytes.Contains(xref[:stream], []byte("/Root")) && bytes.Contains(xref[stream:], []byte("endstream"))
}
