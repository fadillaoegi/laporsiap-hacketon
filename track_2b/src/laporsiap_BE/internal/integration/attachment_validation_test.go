package integration

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func validAttachmentPDF(t *testing.T) []byte {
	t.Helper()
	c := Case{Title: "Fictional attachment reference", Synthetic: true, Revision: 1, RawChat: "Fictional source.", Messages: []Message{{ID: "M1", Text: "Fictional source."}}}
	c.Draft = manualDraft(c)
	data, err := ReportPDF(c, "en")
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func attachmentRequest(t *testing.T, h http.Handler, c caseResponse, name string, data []byte) *httptest.ResponseRecorder {
	t.Helper()
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	part, err := mw.CreateFormFile("file", name)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write(data)
	_ = mw.Close()
	r := httptest.NewRequest("POST", "http://localhost:8787/api/cases/"+c.Case.ID+"/attachments", &body)
	r.Header.Set("Content-Type", mw.FormDataContentType())
	r.Header.Set("Authorization", "Bearer "+c.Token)
	r.Header.Set("X-Case-Revision", stringRevision(c.Case.Revision))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func TestAttachmentValidationChecksPayloadAndPDFEnvelope(t *testing.T) {
	var pngData, jpegData bytes.Buffer
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	_ = png.Encode(&pngData, img)
	_ = jpeg.Encode(&jpegData, img, nil)
	pdfData := validAttachmentPDF(t)
	for _, tc := range []struct {
		name, contentType string
		data              []byte
	}{
		{"PNG", "image/png", pngData.Bytes()},
		{"JPEG", "image/jpeg", jpegData.Bytes()},
		{"PDF", "application/pdf", pdfData},
		{"PDF trailing whitespace", "application/pdf", append(bytes.Clone(pdfData), '\r', '\n', ' ')},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := validateAttachment(tc.data, tc.contentType); err != nil {
				t.Fatal(err)
			}
		})
	}
	corruptPNG := bytes.Clone(pngData.Bytes())
	corruptPNG[len(corruptPNG)-5] ^= 1 // DecodeConfig succeeds; complete decoding must reject the CRC.
	largePNG := bytes.Clone(pngData.Bytes())
	binary.BigEndian.PutUint32(largePNG[16:20], 6000)
	binary.BigEndian.PutUint32(largePNG[20:24], 6000)
	binary.BigEndian.PutUint32(largePNG[29:33], crc32.ChecksumIEEE(largePNG[12:29]))
	for _, tc := range []struct {
		name, contentType string
		data              []byte
	}{
		{"truncated PNG", "image/png", pngData.Bytes()[:32]},
		{"PNG corrupted payload", "image/png", corruptPNG},
		{"image pixel bound", "image/png", largePNG},
		{"truncated JPEG", "image/jpeg", jpegData.Bytes()[:16]},
		{"truncated PDF", "application/pdf", pdfData[:80]},
		{"PDF missing EOF", "application/pdf", pdfData[:len(pdfData)-8]},
		{"PDF fake offset", "application/pdf", []byte("%PDF-1.4\nstartxref\n0\n%%EOF")},
		{"PDF offset outside file", "application/pdf", []byte("%PDF-1.4\nstartxref\n999999999\n%%EOF")},
		{"PDF incomplete table", "application/pdf", []byte("%PDF-1.4\nxref\nstartxref\n9\n%%EOF")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := validateAttachment(tc.data, tc.contentType); err == nil {
				t.Fatal("damaged or oversized attachment accepted")
			}
		})
	}
}

func TestRejectedUploadsLeaveReviewedCaseAndExistingAttachmentUntouched(t *testing.T) {
	s := testStore(t)
	h := (&Server{Store: s}).Handler()
	c := createCase(t, s, h, "Fictional rejected uploads")
	token := c.Token
	w := call(t, h, "POST", "/api/cases/"+c.Case.ID+"/draft", token, map[string]any{"revision": c.Case.Revision, "mode": "manual"})
	requireStatus(t, w, 200)
	c = readCase(t, w)
	c.Token = token
	w = attachmentRequest(t, h, c, "valid.pdf", validAttachmentPDF(t))
	requireStatus(t, w, 200)
	c = readCase(t, w)
	c.Token = token
	c = reviewCase(t, h, c)
	before := c.Case
	for _, tc := range []struct {
		name string
		data []byte
	}{
		{"unsupported.txt", []byte("Fictional unsupported text.")},
		{"empty.png", nil},
		{"truncated.png", []byte("\x89PNG\r\n\x1a\n")},
		{"truncated.jpg", []byte("\xff\xd8\xff")},
		{"truncated.pdf", []byte("%PDF-1.4\nFictional truncated PDF")},
		{"oversize.pdf", bytes.Repeat([]byte{'x'}, 5*1024*1024+1)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			requireStatus(t, attachmentRequest(t, h, c, tc.name, tc.data), 400)
			got := readCase(t, call(t, h, "GET", "/api/cases/"+c.Case.ID, token, nil))
			if !reflect.DeepEqual(before, got.Case) {
				t.Fatal("rejection changed saved data, attachment or review")
			}
		})
	}
}
