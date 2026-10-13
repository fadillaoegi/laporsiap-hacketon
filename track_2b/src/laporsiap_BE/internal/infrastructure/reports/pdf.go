package reports

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"codeberg.org/go-pdf/fpdf"
	"laporsiap/internal/domain"
)

// Font files and their OFL license ship with the backend; export is offline.
//
//go:embed fonts/dm-sans-400.ttf
var pdfRegular []byte

//go:embed fonts/dm-sans-700.ttf
var pdfBold []byte

//go:embed fonts/coverage.json
var pdfCoverageJSON []byte

var pdfCoverage = func() map[rune]bool {
	var codes []int
	_ = json.Unmarshal(pdfCoverageJSON, &codes)
	m := map[rune]bool{}
	for _, n := range codes {
		m[rune(n)] = true
	}
	return m
}()

// Missing glyphs are shown explicitly rather than silently dropped or drawn
// as blank squares. TXT and the JSON snapshot retain the original Unicode.
func PDFText(text string) (string, bool) {
	var b strings.Builder
	escaped := false
	for _, r := range text {
		if r == '\n' || pdfCoverage[r] {
			b.WriteRune(r)
		} else if r == '\r' {
			continue
		} else if r == '\t' {
			b.WriteString("    ")
		} else {
			fmt.Fprintf(&b, "[U+%04X]", r)
			escaped = true
		}
	}
	return b.String(), escaped
}

func ReportPDF(c domain.Case, language string) ([]byte, error) {
	if c.Draft == nil {
		return nil, errors.New("draft required for PDF export")
	}
	label := func(id, en string) string {
		if language == "en" {
			return en
		}
		return id
	}
	value := func(v string) string {
		if v == "" {
			return label("Belum diketahui", "Unknown")
		}
		return v
	}
	amountDisplay := func(v *int64) string {
		if v == nil {
			return label("Belum diketahui", "Unknown")
		}
		text := amountText(v)
		separator := ","
		if language == "id" {
			separator = "."
		}
		for i := len(text) - 3; i > 0; i -= 3 {
			text = text[:i] + separator + text[i:]
		}
		return text
	}
	transferLabels := map[string]string{}
	for i, t := range c.Transfers {
		transferLabels[t.ID] = fmt.Sprintf("Transfer %02d", i+1)
	}
	p := fpdf.New("P", "mm", "A4", "")
	p.SetMargins(18, 25, 18)
	p.SetAutoPageBreak(false, 20)
	p.AddUTF8FontFromBytes("DM", "", pdfRegular)
	p.AddUTF8FontFromBytes("DM", "B", pdfBold)
	p.SetTitle("LaporSiap - "+c.Title, true)
	p.SetAuthor("LaporSiap", true)
	p.AliasNbPages("")
	p.SetHeaderFunc(func() {
		p.SetFont("DM", "B", 11)
		p.SetTextColor(23, 91, 80)
		p.SetXY(18, 12)
		p.CellFormat(174, 6, "LaporSiap", "", 0, "L", false, 0, "")
		p.SetDrawColor(202, 221, 216)
		p.Line(18, 21, 192, 21)
		p.SetY(27)
	})
	p.SetFooterFunc(func() {
		p.SetXY(18, 280)
		p.SetFont("DM", "", 8)
		p.SetTextColor(82, 101, 110)
		p.CellFormat(125, 5, label("Bahan bantu laporan - diperiksa pengguna", "Supporting report materials - reviewed by user"), "", 0, "L", false, 0, "")
		p.CellFormat(49, 5, fmt.Sprintf("%d / {nb}", p.PageNo()), "", 0, "R", false, 0, "")
	})
	p.AddPage()
	room := func(height float64) {
		if p.GetY()+height > 271 {
			p.AddPage()
		}
	}
	block := func(text string, bold bool, size float64, quote bool) {
		style := ""
		if bold {
			style = "B"
		}
		p.SetFont("DM", style, size)
		encoded, _ := PDFText(text)
		width := 174.0
		if quote {
			width = 166
		}
		lines := p.SplitText(encoded, width)
		lineHeight := size * 0.3528 * 1.5
		for _, line := range lines {
			room(lineHeight)
			p.SetFont("DM", style, size)
			p.SetX(18)
			if quote {
				p.SetFillColor(241, 246, 247)
				p.SetTextColor(58, 78, 85)
				p.SetX(22)
			} else if bold {
				p.SetTextColor(24, 60, 58)
			} else {
				p.SetTextColor(36, 60, 67)
			}
			p.CellFormat(width, lineHeight, line, "", 1, "L", quote, 0, "")
		}
		p.SetY(p.GetY() + 3)
	}
	section := func(title string) {
		room(24)
		p.SetY(p.GetY() + 6)
		p.SetDrawColor(202, 221, 216)
		p.Line(18, p.GetY(), 192, p.GetY())
		p.SetY(p.GetY() + 5)
		block(title, true, 13, false)
	}
	block(label("Bahan bantu laporan", "Supporting report materials"), true, 20, false)
	block(c.Title, true, 14, false)
	block(fmt.Sprintf(label("Revisi %d | Review: %s", "Revision %d | Reviewed: %s"), c.Revision, value(c.ReviewedAt)), false, 9, false)
	if c.Synthetic {
		block(label("DATA SIMULASI - BUKAN KASUS ATAU TRANSAKSI NYATA", "SIMULATED DATA - NOT A REAL CASE OR TRANSACTION"), true, 10, true)
	}
	block(label("Dokumen ini bukan laporan resmi, tiket IASC, atau verifikasi bank. Tidak ada pengiriman otomatis. Gunakan kanal resmi: https://iasc.ojk.go.id/.", "This document is not an official report, an IASC ticket, or bank verification. There is no automatic submission. Use the official channel: https://iasc.ojk.go.id/."), false, 9, false)
	if _, escaped := PDFText(ReportInLanguage(c, language)); escaped {
		block(label("Karakter di luar cakupan font ditampilkan sebagai [U+XXXX]. Teks asli lengkap tersedia dalam TXT dan ZIP.", "Characters outside the font coverage are shown as [U+XXXX]. Full original text is available in TXT and ZIP."), false, 9, true)
	}
	if c.Draft != nil {
		origin := label("Manual", "Manual")
		if c.Draft.Origin == "apertus" {
			origin = "Apertus"
		} else if c.Draft.Origin == "apertus_edited" {
			origin = label("Apertus, dikoreksi pengguna", "Apertus, edited by user")
		}
		draftLanguage := label("Belum diketahui", "Unknown")
		if c.Draft.Language == "en" {
			draftLanguage = "English"
		} else if c.Draft.Language == "id" {
			draftLanguage = "Indonesia"
		}
		block(label("Asal draft: ", "Draft origin: ")+origin+" | "+label("Bahasa saat dibuat: ", "Language when generated: ")+draftLanguage, false, 9, false)
		if c.Draft.Model != "" {
			block(label("Model (metadata API): ", "Model (API metadata): ")+c.Draft.Model, false, 9, false)
		}
	}
	s := domain.Summarize(c)
	section(label("Ringkasan catatan transfer", "Transfer record summary"))
	amount := amountDisplay(s.SuccessfulKnownTotal)
	block(label("Nominal transfer berhasil yang dilaporkan pengguna: IDR ", "Successful transfer amount reported by user: IDR ")+amount, true, 12, false)
	block(fmt.Sprintf(label("Berhasil %d | Pending %d | Gagal %d | Belum dipastikan %d", "Successful %d | Pending %d | Failed %d | Unconfirmed %d"), s.StatusCounts["successful"], s.StatusCounts["pending"], s.StatusCounts["failed"], s.StatusCounts["unknown"]), false, 10, false)
	block(fmt.Sprintf(label("%d transfer berhasil memiliki nominal belum diketahui. Total hanya mencakup nominal yang diketahui; bukan kerugian bersih atau jumlah yang diverifikasi bank.", "%d successful transfers have unknown amounts. The total covers known amounts only; it is not net losses or a bank-verified amount."), s.SuccessfulUnknownAmount), false, 9, false)
	section(label("Kronologi dan sumber", "Chronology and sources"))
	if c.Draft != nil {
		for i, e := range c.Draft.Events {
			// Keep short entries with their quote; long sources may continue on
			// the next page without clipping or truncation.
			p.SetFont("DM", "", 11)
			description, _ := PDFText(e.Description)
			descriptionLines := len(p.SplitText(description, 174))
			p.SetFont("DM", "", 10)
			evidence, _ := PDFText(e.Evidence)
			height := float64(descriptionLines)*5.82 + float64(len(p.SplitText(evidence, 166)))*5.29 + 38
			for _, a := range c.Attachments {
				if containsID(a.SourceIDs, e.SourceID) {
					height += 9
				}
			}
			if height < 244 {
				room(height)
			} else {
				room(28)
			}
			block(fmt.Sprintf("%02d / %s", i+1, e.SourceID), true, 10, false)
			block(e.Description, false, 11, false)
			if i < len(c.Messages) {
				m := c.Messages[i]
				block(label("Pengirim: ", "Sender: ")+value(m.Sender)+" | "+label("Waktu pesan: ", "Message time: ")+value(m.SentAt), false, 8, false)
			}
			block(label("Waktu pesan bukan tanggal transfer. Sumber asli:", "Message time is not the transfer date. Original source:"), false, 8, false)
			block(e.Evidence, false, 10, true)
			for _, a := range c.Attachments {
				for _, id := range a.SourceIDs {
					if id == e.SourceID {
						block(label("Lampiran terkait (pilihan pengguna): ", "Linked attachment (user selected): ")+a.Filename, false, 8, false)
					}
				}
			}
		}
	}
	section(label("Rincian transfer dari input pengguna", "Transfer details from user input"))
	if len(c.Transfers) == 0 {
		block(label("Tidak ada catatan transfer.", "No transfer records."), false, 10, false)
	}
	for i, t := range c.Transfers {
		room(28)
		block(fmt.Sprintf("Transfer %02d", i+1), true, 11, false)
		block(value(t.FromBank)+" / "+value(t.FromAccount)+" -> "+value(t.ToBank)+" / "+value(t.ToAccount), false, 10, false)
		block(label("Nominal: IDR ", "Amount: IDR ")+amountDisplay(t.Amount)+" | "+reportText(statusLabels[t.Status], language), false, 10, false)
		block(label("Tanggal transfer: ", "Transfer date: ")+value(t.OccurredAt)+" | "+label("Catatan: ", "Note: ")+value(t.Note), false, 9, false)
		for _, a := range c.Attachments {
			for _, id := range a.TransferIDs {
				if id == t.ID {
					block(label("Lampiran terkait (pilihan pengguna): ", "Linked attachment (user selected): ")+a.Filename, false, 8, false)
				}
			}
		}
	}
	section(label("Data yang perlu dilengkapi atau diperiksa", "Details to complete or review"))
	for _, w := range s.Warnings {
		block("- "+reportText(w, language), false, 10, false)
	}
	block(label("Setiap ringkasan telah ditinjau pengguna. Pemeriksaan otomatis dapat melewatkan kesalahan; tidak menjamin kebenaran.", "Each summary has been reviewed by the user. Automated checks may miss errors and do not guarantee accuracy."), false, 9, false)
	section(label("Daftar lampiran referensi", "Reference attachment index"))
	block(label("Hubungan bukti dipilih pengguna. Isi dan keaslian file tidak diverifikasi; belum termasuk OCR. File asli disertakan dalam ZIP.", "Evidence links are selected by the user. File contents and authenticity are not verified; OCR is not included. Original files are included in the ZIP."), false, 9, false)
	for i, a := range c.Attachments {
		room(22)
		block(fmt.Sprintf("%02d / %s", i+1, a.Filename), true, 10, false)
		block(fmt.Sprintf("%s | %d bytes", a.ContentType, a.Size), false, 9, false)
		linkedTransfers := []string{}
		for _, id := range a.TransferIDs {
			linkedTransfers = append(linkedTransfers, value(transferLabels[id]))
		}
		block(label("Pesan: ", "Messages: ")+value(strings.Join(a.SourceIDs, ", "))+" | "+label("Transfer: ", "Transfers: ")+value(strings.Join(linkedTransfers, ", ")), false, 9, false)
		block("SHA-256: "+a.SHA256, false, 8, false)
	}
	var out bytes.Buffer
	if err := p.Output(&out); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func containsID(ids []string, id string) bool {
	for _, value := range ids {
		if value == id {
			return true
		}
	}
	return false
}
