package reports

import (
	"archive/zip"
	"bytes"

	"encoding/csv"
	"encoding/json"

	"fmt"

	"laporsiap/internal/domain"
	"sort"
	"strings"
)

func amountText(v *int64) string {
	if v == nil {
		return "Belum diketahui"
	}
	return fmt.Sprintf("%d", *v)
}
func filled(v string) string {
	if v == "" {
		return "Belum diketahui"
	}
	return v
}

var statusLabels = map[string]string{"successful": "Berhasil (dilaporkan pengguna)", "pending": "Dalam proses", "failed": "Gagal", "unknown": "Belum dipastikan"}

func Report(c domain.Case) string { return ReportInLanguage(c, "id") }

func ReportInLanguage(c domain.Case, language string) string {
	tr := func(text string) string { return reportText(text, language) }
	value := func(text string) string {
		if text == "" {
			return tr("Belum diketahui")
		}
		return text
	}
	amount := func(v *int64) string {
		if v == nil {
			return tr("Belum diketahui")
		}
		return amountText(v)
	}
	var b strings.Builder
	fmt.Fprintf(&b, tr("LAPORSIAP — BAHAN BANTU LAPORAN\n%s\nRevisi: %d | Diperiksa: %s\n"), c.Title, c.Revision, c.ReviewedAt)
	if c.Synthetic {
		b.WriteString(tr("DATA SIMULASI — BUKAN KASUS ATAU TRANSAKSI NYATA\n"))
	}
	b.WriteString(tr("\nDokumen ini bukan tiket IASC, bukti laporan diterima, atau verifikasi keaslian transaksi. Tidak ada pengiriman otomatis. Gunakan kanal resmi https://iasc.ojk.go.id/.\n\nKRONOLOGI YANG DITINJAU PENGGUNA\n"))
	fmt.Fprintf(&b, tr("Asal draft: %s | Model (metadata API): %s | Bahasa saat dibuat: %s\n"), c.Draft.Origin, value(c.Draft.Model), value(c.Draft.Language))
	for i, e := range c.Draft.Events {
		m := c.Messages[i]
		fmt.Fprintf(&b, tr("\n%d. [%s] %s\nPengirim: %s | Waktu pesan: %s (bukan waktu transfer)\nSumber asli: %s\n"), i+1, e.SourceID, e.Description, value(m.Sender), value(m.SentAt), e.Evidence)
	}
	s := domain.Summarize(c)
	fmt.Fprintf(&b, tr("\nRINCIAN TRANSFER DARI INPUT PENGGUNA\nNominal transfer berhasil yang dilaporkan pengguna (nominal diketahui): %s rupiah\nTransfer berhasil dengan nominal belum diketahui: %d\nTotal ini bukan kerugian bersih atau jumlah yang diverifikasi bank. Pending/gagal/unknown tidak dihitung.\n"), amount(s.SuccessfulKnownTotal), s.SuccessfulUnknownAmount)
	groups := map[string][]domain.Transfer{}
	for _, t := range c.Transfers {
		key := fmt.Sprintf("%s / %s → %s / %s", value(t.FromBank), value(t.FromAccount), value(t.ToBank), value(t.ToAccount))
		groups[key] = append(groups[key], t)
	}
	keys := []string{}
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Fprintf(&b, tr("\nPasangan rekening: %s\n"), key)
		for _, t := range groups[key] {
			fmt.Fprintf(&b, tr("- %s | %s rupiah | Tanggal transfer: %s | %s | Catatan: %s\n"), t.ID, amount(t.Amount), value(t.OccurredAt), tr(statusLabels[t.Status]), value(t.Note))
		}
	}
	b.WriteString(tr("\nDATA YANG MASIH PERLU DILENGKAPI / DIPERIKSA\n"))
	for _, w := range s.Warnings {
		fmt.Fprintf(&b, tr("- %s\n"), tr(w))
	}
	b.WriteString(tr("\nLAMPIRAN REFERENSI (TANPA OCR / VERIFIKASI KEASLIAN)\n"))
	for _, a := range c.Attachments {
		fmt.Fprintf(&b, tr("- %s | %s | %d byte | SHA-256: %s\n"), a.ID, a.Filename, a.Size, a.SHA256)
		fmt.Fprintf(&b, tr("  Hubungan pilihan pengguna — pesan: %s | transfer: %s\n"), value(strings.Join(a.SourceIDs, ", ")), value(strings.Join(a.TransferIDs, ", ")))
	}
	b.WriteString(tr("\nDisusun dengan review manusia. Pemeriksaan heuristik bukan jaminan kebenaran.\n"))
	return b.String()
}
func ExportZIP(c domain.Case, attachments []domain.AttachmentData) ([]byte, error) {
	return ExportZIPInLanguage(c, attachments, "id")
}

func ExportZIPInLanguage(c domain.Case, attachments []domain.AttachmentData, language string) ([]byte, error) {
	var b bytes.Buffer
	z := zip.NewWriter(&b)
	write := func(name string, data []byte) error {
		w, err := z.Create(name)
		if err != nil {
			return err
		}
		_, err = w.Write(data)
		return err
	}
	if err := write("kronologi.txt", []byte(ReportInLanguage(c, language))); err != nil {
		return nil, err
	}
	pdf, err := ReportPDF(c, language)
	if err != nil {
		return nil, err
	}
	if err = write("laporan.pdf", pdf); err != nil {
		return nil, err
	}
	if err := write("percakapan-asli.txt", []byte(c.RawChat)); err != nil {
		return nil, err
	}
	snapshot, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return nil, err
	}
	if err = write("kasus-direview.json", snapshot); err != nil {
		return nil, err
	}
	var csvBytes bytes.Buffer
	cw := csv.NewWriter(&csvBytes)
	headers := []string{"id", "bank_pengirim", "rekening_pengirim", "bank_penerima", "rekening_penerima", "nominal_rupiah", "tanggal_transfer", "status", "catatan"}
	if language == "en" {
		headers = []string{"id", "sending_bank", "sending_account", "receiving_bank", "receiving_account", "amount_rupiah", "transfer_date", "status", "note"}
	}
	if err := cw.Write(headers); err != nil {
		return nil, err
	}
	for _, t := range c.Transfers {
		amount := ""
		if t.Amount != nil {
			amount = fmt.Sprint(*t.Amount)
		}
		row := []string{t.ID, t.FromBank, t.FromAccount, t.ToBank, t.ToAccount, amount, t.OccurredAt, t.Status, t.Note}
		for i, v := range row {
			trimmed := strings.TrimLeft(v, " \t\r\n")
			if strings.HasPrefix(trimmed, "=") || strings.HasPrefix(trimmed, "+") || strings.HasPrefix(trimmed, "-") || strings.HasPrefix(trimmed, "@") {
				row[i] = "'" + v
			}
		}
		if err = cw.Write(row); err != nil {
			return nil, err
		}
	}
	cw.Flush()
	if err = cw.Error(); err != nil {
		return nil, err
	}
	if err = write("transfer.csv", csvBytes.Bytes()); err != nil {
		return nil, err
	}
	for _, a := range attachments {
		if err = write("lampiran/"+a.Meta.ID+"-"+a.Meta.Filename, a.Bytes); err != nil {
			return nil, err
		}
	}
	if err = z.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
