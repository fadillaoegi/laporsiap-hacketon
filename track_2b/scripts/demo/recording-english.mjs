export const synthetic = {
  title: 'Online task offer and deposit requests',
  rawChat: '[08/10/2026, 09:10] Nadia: Hi Raka, we offer tasks that involve liking posts, with a promised commission of IDR 50,000 per day.\n[08/10/2026, 09:14] Nadia: To continue the tasks, please pay a deposit of IDR 180,000 to ****4821.\n[08/10/2026, 09:28] Raka: I have just sent IDR 180,000. I have not recorded the transfer time yet.\n[08/10/2026, 09:35] Nadia: Now add a deposit of IDR 420,000 so the balance can be withdrawn.\n[08/10/2026, 09:40] Raka: The additional deposit is still pending. I have not seen the transaction on my bank statement.',
  transfers: [
    { fromBank: 'Bank Aruna', fromAccount: '****7392', toBank: 'Bank Lestari', toAccount: '****4821', amount: '180000', status: 'successful', note: 'Transfer date not recorded; do not use the message timestamp.' },
    { fromBank: 'Bank Aruna', fromAccount: '****7392', toBank: 'Bank Lestari', toAccount: '****4821', amount: '420000', status: 'pending', note: 'Not yet confirmed as successful.' },
  ],
  descriptions: [
    'Nadia offers tasks involving liking posts, with a promised commission of IDR 50,000 per day.',
    'Nadia requests a deposit of IDR 180,000 to ****4821 to continue the tasks.',
    'Raka states that they have just sent IDR 180,000 and have not recorded the transfer time.',
    'Nadia requests an additional deposit of IDR 420,000 so the balance can be withdrawn.',
    'Raka states that the additional deposit is still pending and is not yet visible on their bank statement.',
  ],
  pdfLines: [
    [22, 'LaporSiap - SIMULATED DATA'],
    [12, 'Fictional reference attachment. Not a bank receipt.'],
    [12, 'Bank Aruna / ****7392'],
    [12, 'Destination: Bank Lestari / ****4821'],
    [12, 'User record: IDR 180,000 marked successful.'],
    [12, 'Additional record: IDR 420,000 still pending.'],
    [12, 'Transfer date unknown; message time is not transfer time.'],
    [12, 'All names, accounts and transactions are fictional.'],
  ],
}

export const englishText = {
  '01 — Buka aplikasi dan muat contoh': '01 — Open the app and import the English fictional case',
  '02 — Periksa status dan simpan bahan': '02 — Check transfer status and save materials',
  '03 — Tambahkan lampiran sintetis': '03 — Upload the fictional reference PDF',
  '04 — Susun draft manual': '04 — Create a manual draft',
  '05 — Cocokkan dan sunting kelima entri': '05 — Compare sources and edit all five entries',
  '06 — Review sumber, transfer, dan data kosong': '06 — Review sources, transfers and unknown details',
  '07 — Unduh ZIP dan kronologi TXT': '07 — Download ZIP and chronology TXT',
  '08 — Periksa riwayat dan hasil akhir': '08 — Check revision history and the final result',
  'SIMULASI LAPORSIAP • DATA SINTETIS • DRAFT MANUAL': 'LAPORSIAP WALKTHROUGH • SYNTHETIC DATA • MANUAL DRAFT',
  'Mulai dari contoh rekaan yang tersedia di aplikasi.': 'Import an English fictional conversation and enter the fictional transfer records.',
  '01 / BAHAN KASUS': '01 / CASE MATERIALS',
  'Rp180.000 berhasil; Rp420.000 masih pending. Tanggal transfer tetap kosong.': 'IDR 180,000 is successful; IDR 420,000 is pending. Transfer dates remain blank.',
  '01 / LAMPIRAN REFERENSI': '01 / REFERENCE ATTACHMENT',
  'Unggah PDF rekaan. Lampiran disimpan sebagai referensi, tanpa OCR.': 'Upload a fictional PDF. It is stored as a reference attachment, without OCR.',
  '02 / SUSUN DRAFT MANUAL': '02 / CREATE A MANUAL DRAFT',
  'Pada rekaman ini, draft dimulai dari kutipan pesan asli. Apertus tidak dipanggil.': 'This recording starts from original message quotes. No live Apertus call is made.',
  'Baca pesan asli di kanan, lalu rapikan kalimat tanpa menambah fakta.': 'Read the original message on the right. Edit the wording without adding facts.',
  '03 / SIMPAN KOREKSI': '03 / SAVE EDITS',
  'Koreksi disimpan sebagai revisi. Selanjutnya periksa semua sumber kembali.': 'Save the edits as a revision, then review all sources again.',
  'Tandai entri setelah membandingkan draft dengan pesan sumbernya.': 'Mark each entry after comparing the draft with its original source message.',
  '03 / PERIKSA CATATAN TRANSFER': '03 / REVIEW TRANSFER RECORDS',
  'Total Rp180.000 berasal dari status berhasil. Rp420.000 pending tidak masuk total.': 'The total is IDR 180,000 from successful transfers. The pending IDR 420,000 is excluded.',
  '03 / DATA BELUM LENGKAP': '03 / UNKNOWN DETAILS',
  'Tanggal transfer belum diketahui. Ekspor adalah bahan bantu untuk kanal resmi.': 'Transfer dates are unknown. The export provides supporting materials for official channels.',
  '04 / EKSPOR BAHAN': '04 / EXPORT MATERIALS',
  'Unduh paket ZIP dan kronologi TXT dari revisi yang telah dikonfirmasi.': 'Download the ZIP package and chronology TXT from the confirmed revision.',
  'SIMULASI SELESAI': 'WALKTHROUGH COMPLETE',
  'Input → lampiran → draft manual → koreksi → review → ZIP dan TXT. Semua berhasil.': 'Input → attachment → manual draft → edits → review → ZIP and TXT. All steps passed.',
  'Ini data simulasi': 'This is simulated data', 'Percakapan': 'Conversation',
  'Catatan transfer': 'Transfer records', 'Simpan & lanjut': 'Save & continue',
  'Susun dengan Apertus': 'Draft with Apertus', 'Pilih file bukti': 'Choose an evidence file',
  'Susun draft dengan Apertus': 'Create draft with Apertus', 'Gunakan draft manual': 'Use a manual draft',
  'DRAFT MANUAL': 'MANUAL DRAFT', 'Simpan koreksi': 'Save edits',
  'Koreksi tersimpan.': 'Edits saved.',
  'Saya telah mencocokkan entri ini dengan sumber': 'I have compared this entry with its source',
  'Periksa catatan transfer': 'Review transfer records', 'Konfirmasi & lanjut': 'Confirm & continue',
  'Unduh paket ZIP': 'Download ZIP package', 'Unduh kronologi TXT': 'Download chronology TXT',
  'Lihat riwayat revisi →': 'View revision history →',
}

export function recordingText(text) {
  if (englishText[text]) return englishText[text]
  if (/^03 \/ PERIKSA SUMBER M\d+$/.test(text)) return text.replace('PERIKSA SUMBER', 'CHECK SOURCE')
  if (/^03 \/ REVIEW M\d+ DARI 5$/.test(text)) return text.replace('DARI 5', 'OF 5')
  throw new Error(`Missing English recording text: ${text}`)
}
