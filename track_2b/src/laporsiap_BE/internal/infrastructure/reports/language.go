package reports

import (
	"fmt"
	"regexp"
	"strings"
)

// Translate only application-owned labels. Evidence, descriptions, and other
// user-authored fields remain verbatim, including words that match these keys.
var reportEnglish = map[string]string{
	"Asal draft: %s | Model (metadata API): %s | Bahasa saat dibuat: %s\n": "Draft origin: %s | Model (API metadata): %s | Language when generated: %s\n",
	"  Hubungan pilihan pengguna — pesan: %s | transfer: %s\n":             "  User-selected links — messages: %s | transfers: %s\n",
	"Belum diketahui": "Unknown",
	"LAPORSIAP — BAHAN BANTU LAPORAN\n%s\nRevisi: %d | Diperiksa: %s\n": "LAPORSIAP — SUPPORTING REPORT MATERIALS\n%s\nRevision: %d | Reviewed: %s\n",
	"DATA SIMULASI — BUKAN KASUS ATAU TRANSAKSI NYATA\n":                "SIMULATED DATA — NOT A REAL CASE OR TRANSACTION\n",
	"\nDokumen ini bukan tiket IASC, bukti laporan diterima, atau verifikasi keaslian transaksi. Tidak ada pengiriman otomatis. Gunakan kanal resmi https://iasc.ojk.go.id/.\n\nKRONOLOGI YANG DITINJAU PENGGUNA\n":                                                                                "\nThis document is not an IASC ticket, proof of report acceptance, or verification of transaction authenticity. There is no automatic submission. Use the official channel https://iasc.ojk.go.id/.\n\nCHRONOLOGY REVIEWED BY THE USER\n",
	"\n%d. [%s] %s\nPengirim: %s | Waktu pesan: %s (bukan waktu transfer)\nSumber asli: %s\n":                                                                                                                                                                                                      "\n%d. [%s] %s\nSender: %s | Message time: %s (not transfer time)\nOriginal source: %s\n",
	"\nRINCIAN TRANSFER DARI INPUT PENGGUNA\nNominal transfer berhasil yang dilaporkan pengguna (nominal diketahui): %s rupiah\nTransfer berhasil dengan nominal belum diketahui: %d\nTotal ini bukan kerugian bersih atau jumlah yang diverifikasi bank. Pending/gagal/unknown tidak dihitung.\n": "\nTRANSFER DETAILS FROM USER INPUT\nSuccessful transfer amount reported by the user (known amounts): %s rupiah\nSuccessful transfers with unknown amounts: %d\nThis total is not net losses or a bank-verified amount. Pending/failed/unconfirmed transfers are excluded.\n",
	"\nPasangan rekening: %s\n":                                                         "\nAccount pair: %s\n",
	"- %s | %s rupiah | Tanggal transfer: %s | %s | Catatan: %s\n":                      "- %s | %s rupiah | Transfer date: %s | %s | Note: %s\n",
	"\nDATA YANG MASIH PERLU DILENGKAPI / DIPERIKSA\n":                                  "\nDETAILS STILL TO COMPLETE / REVIEW\n",
	"\nLAMPIRAN REFERENSI (TANPA OCR / VERIFIKASI KEASLIAN)\n":                          "\nREFERENCE ATTACHMENTS (NO OCR / AUTHENTICITY VERIFICATION)\n",
	"\nDisusun dengan review manusia. Pemeriksaan heuristik bukan jaminan kebenaran.\n": "\nPrepared with human review. Heuristic checks do not guarantee accuracy.\n",
	"Berhasil (dilaporkan pengguna)":                                                    "Successful (reported by the user)", "Dalam proses": "Pending", "Gagal": "Failed", "Belum dipastikan": "Unconfirmed",
	"Total sebagian: ada transfer berhasil dengan nominal belum diketahui.": "Partial total: some successful transfers have unknown amounts.",
	"Belum ada lampiran bukti. Lampiran dapat ditambahkan tanpa OCR.":       "No evidence attachments yet. Attachments can be added without OCR.",
}

var missingTransferFields = regexp.MustCompile(`^Transfer (\d+): (.+) belum lengkap\.$`)
var duplicateTransfer = regexp.MustCompile(`^Transfer (\d+) memiliki rincian yang sama dengan catatan lain; periksa duplikasi\.$`)

func reportText(text, language string) string {
	if language != "en" {
		return text
	}
	if translated, ok := reportEnglish[text]; ok {
		return translated
	}
	if match := missingTransferFields.FindStringSubmatch(text); match != nil {
		fields := strings.Split(match[2], ", ")
		labels := map[string]string{"nominal": "amount", "tanggal transfer": "transfer date", "rekening pengirim": "sending account", "rekening penerima": "receiving account", "status": "status"}
		for i, field := range fields {
			if label, ok := labels[field]; ok {
				fields[i] = label
			}
		}
		return fmt.Sprintf("Transfer %s: missing %s.", match[1], strings.Join(fields, ", "))
	}
	if match := duplicateTransfer.FindStringSubmatch(text); match != nil {
		return fmt.Sprintf("Transfer %s has the same details as another record; check for duplication.", match[1])
	}
	return text
}
