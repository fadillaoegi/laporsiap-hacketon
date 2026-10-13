import { useEffect } from 'react'
import { useAppDispatch, useAppSelector } from './state/hooks'
import { localeSlice } from './state/localeSlice'
import type { Language } from '../domain/models'
import type { LanguagePreferences } from '../application/ports'
export type { Language } from '../domain/models'
import type { ReactNode } from 'react'

// Indonesian keys also match the API's validation messages. User-authored
// evidence, summaries, titles, and notes are never passed to this translator.
const english: Record<string, string> = {
  'Samarkan informasi pribadi sebelum dikirim ke AI': 'Mask personal information before sending to AI',
  'Mendeteksi nama pengirim, email, telepon Indonesia, serta nomor rekening/identitas berlabel. Deteksi dapat melewatkan data atau menyamarkan bagian yang tidak Anda maksud; periksa preview.': 'Detects sender names, emails, Indonesian phone numbers, and labelled account/identity numbers. Detection may miss data or mask unintended text; check the preview.',
  'Istilah tambahan untuk disamarkan': 'Additional terms to mask',
  'Satu nama atau istilah per baris; maksimal 50.': 'One name or term per line; up to 50.',
  'Menyiapkan preview pribadi…': 'Preparing private preview…',
  'Informasi pribadi disamarkan:': 'Personal information masked:', 'Teks asli akan dikirim:': 'Original text will be sent:', 'bagian disamarkan.': 'parts masked.',
  'Hanya ID dan teks preview di atas dikirim, bersama instruksi penyusunan. Sumber asli disimpan untuk review dan ekspor. File lampiran dan isian transfer tidak dikirim.': 'Only the source IDs and preview text above are sent, along with drafting instructions. Original sources remain available for review and export. Attachments and transfer fields are not sent.',
  'Pilih teks disamarkan atau teks asli untuk AI.': 'Choose masked or original text for AI.',
  'Maksimal 50 istilah pribadi, masing-masing 2–160 karakter.': 'Use up to 50 private terms, each 2–160 characters.',
  'Preview AI berubah. Periksa teks dan setujui kembali.': 'AI preview changed. Check the text and consent again.',
  'Data dan akses Anda': 'Your data and access', 'Retensi kasus:': 'Case retention:', 'hari sejak dibuat.': 'days from creation.',
  'Batas AI per kasus:': 'AI limit per case:', 'percobaan per 24 jam.': 'attempts per 24 hours.',
  'Kasus yang kedaluwarsa tidak dapat diakses dan akan dihapus dari database. Unduhan dan salinan cadangan terpisah tidak ikut dihapus.': 'Expired cases cannot be accessed and will be deleted from the database. Separate downloads and backups are not deleted.',
  'Kasus kedaluwarsa': 'Case expires', 'Akses sesi kedaluwarsa': 'Session access expires',
  'Perbarui akses untuk mencabut akses lama. Tab lain dengan akses lama perlu ditutup; sesi ini menerima akses baru. Retensi kasus tidak diperpanjang.': 'Renew access to revoke the old access. Close other tabs using the old access; this session receives new access. Case retention is not extended.',
  'Perbarui akses sesi': 'Renew session access', 'Hapus kasus ini': 'Delete this case',
  'Penghapusan permanen mencakup percakapan, draft, riwayat, dan lampiran dalam database. Unduhan Anda tetap tersimpan di luar aplikasi.': 'Permanent deletion includes the conversation, draft, history, and attachments in the database. Your downloads remain outside the app.',
  'Ketik DELETE untuk menghapus': 'Type DELETE to delete', 'Hapus kasus secara permanen': 'Permanently delete case',
  'Tutup tab akan mengakhiri akses sesi; data tersimpan sampai retensi berakhir. Anda dapat menghapus kasus saat akses masih aktif.': 'Closing the tab ends session access; data remains until retention expires. You can delete a case while access is active.',
  'Memperbarui akses…': 'Renewing access…', 'Menghapus kasus…': 'Deleting case…',
  'Akses lama dicabut. Akses sesi ini telah diperbarui.': 'Old access revoked. This session has new access.',
  'Kasus, riwayat, dan lampiran telah dihapus dari database.': 'Case, history, and attachments deleted from the database.',
  'Konfirmasi penghapusan kasus diperlukan.': 'Confirm case deletion first.',
  'Kapasitas penyimpanan atau revisi tercapai. Hapus kasus yang tidak diperlukan atau hubungi pengelola.': 'Storage or revision capacity reached. Delete unneeded cases or contact the operator.',
  'Batas panggilan AI harian tercapai. Draft manual tetap tersedia.': 'Daily AI request limit reached. Manual drafting remains available.',
  'Apertus sedang melayani permintaan lain. Coba nanti atau gunakan draft manual.': 'Apertus is serving other requests. Try later or use a manual draft.',
  'Terlalu banyak permintaan. Tunggu satu menit lalu coba lagi.': 'Too many requests. Wait one minute and try again.',
  'Sumber menyebut status belum selesai atau gagal. Periksa apakah draft mengubahnya menjadi berhasil.': 'The source describes an incomplete or failed transfer. Check whether the draft changes it to successful.',
  'Sumber berisi permintaan atau janji. Periksa apakah draft menyatakannya sebagai kejadian yang sudah terjadi.': 'The source contains a request or promise. Check whether the draft presents it as an event that already happened.',
  'Draft memuat nominal rupiah yang tidak ditemukan dalam sumber ini. Cocokkan setiap nominal.': 'The draft includes a rupiah amount not found in this source. Check each amount.',
  'Periksa teks sebelum dikirim': 'Review the text before sending',
  'Draft baru memakai bahasa pilihan Anda. Kutipan sumber tetap dalam bahasa aslinya.': 'New drafts use your selected language. Source quotes stay in their original language.',
  'Periksa waktu dalam pesan: draft mungkin menghilangkan atau mengubah keterangan seperti besok atau pekan depan.': 'Check the timing in the message: the draft may omit or change expressions such as tomorrow or next week.',
  'Lihat seluruh pesan yang akan dikirim': 'View all messages to be sent',
  'Hanya ID dan isi pesan di atas yang menjadi data kasus untuk AI. Instruksi penyusunan juga disertakan. Nama atau rekening yang tertulis di dalam pesan ikut terkirim.': 'Only the IDs and message text above are sent as case data to AI, together with drafting instructions. Names or account numbers written inside messages are also sent.',
  'Koreksi teks sebelum mengirim →': 'Edit the text before sending →',
  'Bahasa draft baru:': 'New draft language:', 'Bahasa saat draft dibuat:': 'Language when draft was created:', 'Belum tercatat': 'Not recorded',
  'Hubungkan lampiran': 'Link attachment', 'HUBUNGAN BUKTI': 'EVIDENCE LINKS',
  'Pilih bagian yang didukung lampiran menurut Anda. Hubungan ini tidak memverifikasi isi atau keaslian file.': 'Choose the items this attachment supports according to you. Links do not verify file contents or authenticity.',
  'Pesan sumber': 'Source messages', 'Simpan hubungan bukti': 'Save evidence links', 'Batal': 'Cancel',
  'Menyimpan hubungan bukti…': 'Saving evidence links…',
  'Hubungan bukti tersimpan. Review perlu dikonfirmasi ulang.': 'Evidence links saved. Review must be confirmed again.',
  'Hubungan bukti diperbarui': 'Evidence links updated', 'Hubungkan bukti': 'Link evidence',
  'Belum dihubungkan ke pesan': 'No linked messages', 'Belum dihubungkan ke transfer': 'No linked transfers',
  'Lampiran terkait:': 'Linked attachment:', 'Menyiapkan PDF…': 'Preparing PDF…', 'Unduh laporan PDF': 'Download PDF report',
  'Tautan bukti harus menunjuk sumber dan transfer dalam kasus ini, tanpa duplikasi.': 'Evidence links must refer to sources and transfers in this case, without duplicates.',
  'Bahasa draft harus Indonesia atau Inggris.': 'Draft language must be Indonesian or English.',
  'PDF gagal dibuat. TXT dan ZIP tetap tersedia.': 'PDF generation failed. TXT and ZIP are still available.',
  'Pilih format TXT, ZIP, atau PDF.': 'Choose TXT, ZIP, or PDF format.',
  'Pesan sumber diperiksa': 'Source messages reviewed', 'Review revisi ini sudah dikonfirmasi.': 'This revision has already been reviewed and confirmed.', 'Simpan koreksi, periksa setiap sumber, lalu konfirmasikan bahan.': 'Save your edits, check every source, then confirm the materials.',
  'Belum diketahui': 'Unknown', 'Berhasil': 'Successful', 'Dalam proses': 'Pending', 'Gagal': 'Failed', 'Belum dipastikan': 'Unconfirmed',
  'Bahan kasus': 'Case materials', 'Susun draft': 'Create draft', 'Periksa & koreksi': 'Review & edit', 'Ekspor bahan': 'Export materials',
  'Periksa arah dana: mengirim atau menerima.': 'Check the direction of funds: sent or received.',
  'kasus ini ke provider': 'from this case to provider',
  'Configured Apertus provider': 'Configured Apertus provider',
  'Local Apertus server': 'Local Apertus server',
  'Sumber menyebut pengiriman uang. Periksa apakah draft menghilangkannya atau mengaitkan nominal dengan pengembalian uang.': 'The source mentions sending money. Check whether the draft omits it or associates the amount with a refund.',
  'Sumber menyebut perekrut. Periksa apakah draft mengubah pihak yang meminta atau menjanjikan sesuatu.': 'The source mentions a recruiter. Check whether the draft changes who makes the request or promise.',
  'Sumber memuat nominal rupiah yang tidak disebut dalam draft. Periksa apakah nominal penting terlewat.': 'The source contains a rupiah amount omitted from the draft. Check whether an important amount is missing.',
  'Sumber memuat instruksi kepada AI. Periksa agar instruksi tidak dianggap fakta.': 'The source contains instructions to AI. Check that they are not treated as facts.',
  'Kasus dibuat': 'Case created', 'Input diperbarui': 'Input updated', 'Draft manual dibuat': 'Manual draft created', 'Draft Apertus dibuat': 'Apertus draft created',
  'Draft dikoreksi': 'Draft edited', 'Review dikonfirmasi': 'Review confirmed', 'Lampiran ditambahkan': 'Attachment added', 'Lampiran dihapus': 'Attachment removed',
  'Terjadi kesalahan.': 'An error occurred.', 'Memuat kasus…': 'Loading case…',
  'Judul dan percakapan wajib diisi.': 'Enter a title and conversation.',
  'Nominal harus rupiah bulat nonnegatif atau kosong.': 'Enter a nonnegative whole rupiah amount, or leave it blank.',
  'Bahan kasus tersimpan.': 'Case materials saved.',
  'Draft manual memuat kutipan pesan. Periksa dan sunting tiap entri.': 'The manual draft contains message quotes. Review and edit each entry.',
  'Draft Apertus tersimpan. Seluruh entri tetap perlu review makna.': 'Apertus draft saved. Review the meaning of every entry.',
  'Koreksi tersimpan. Periksa revisi terbaru sebelum konfirmasi.': 'Edits saved. Review the latest revision before confirming.',
  'Review revisi terbaru telah dikonfirmasi. Bahan siap diekspor.': 'The latest revision has been reviewed and confirmed. Materials are ready to export.',
  'Maksimal ukuran lampiran 5 MB.': 'Attachments must be no larger than 5 MB.',
  'Lampiran tersimpan. Konfirmasi review lama dibatalkan.': 'Attachment saved. The previous review confirmation has been cleared.',
  'PROTOTIPE LOKAL': 'LOCAL PROTOTYPE', 'Buka kanal resmi IASC': 'Open the official IASC channel', 'RUANG PERSIAPAN': 'PREPARATION SPACE',
  'Tahapan menyiapkan laporan': 'Report preparation steps', 'Teks, transfer, dan bukti': 'Text, transfers, and evidence', 'Apertus atau manual': 'Apertus or manual',
  'Cocokkan dengan sumber': 'Compare with sources', 'Untuk kanal resmi': 'For official channels',
  'Anda tetap memegang kendali.': 'You stay in control.',
  'Draft membantu merapikan bahan. Anda memeriksa isinya dan mengirim laporan melalui kanal resmi.': 'The draft helps organize your materials. You review the contents and submit a report through official channels.',
  'Kasus dalam sesi ini': 'Cases in this session', 'Pilih kasus': 'Select a case', '+ Buat kasus baru': '+ Create a new case',
  'Tanpa unggah KTP.': 'No identity card upload.', 'Tanpa pengiriman laporan otomatis.': 'No automatic report submission.',
  'Perlu segera melapor?': 'Need to report now?',
  'Anda dapat langsung mengakses kanal resmi, tanpa menunggu draft AI.': 'You can access official channels immediately, without waiting for an AI draft.',
  'Buka IASC': 'Open IASC', 'LANGKAH 0': 'STEP 0',
  'Mulai dari bahan yang Anda punya.': 'Start with the materials you have.', 'Rapikan cerita, tetap pegang sumber.': 'Organize your story with its sources.',
  'Periksa setiap bagian ceritanya.': 'Review every part of the story.', 'Bahan laporan Anda, siap dibawa.': 'Your report materials, ready to use.',
  'Percakapan yang belum rapi bisa menjadi titik awal. Data yang belum diketahui boleh tetap kosong.': 'An unorganized conversation can be a starting point. Leave unknown details blank.',
  'Apertus menyusun draft dari percakapan. Rincian transfer tetap berasal dari isian Anda.': 'Apertus drafts a chronology from the conversation. Transfer details come from your entries.',
  'Bandingkan draft dengan pesan asli. Koreksi isi yang berbeda, lalu konfirmasikan review.': 'Compare the draft with the original messages. Correct any differences, then confirm your review.',
  'Unduh bahan yang telah diperiksa untuk digunakan pada kanal resmi pilihan Anda.': 'Download your reviewed materials to use with the official channel of your choice.',
  'DATA SIMULASI — identitas, percakapan, dan transaksi ini rekaan.': 'SIMULATED DATA — identities, conversations and transactions are fictional.',
  'Memuat revisi terbaru…': 'Loading the latest revision…', 'Revisi terbaru dimuat.': 'Latest revision loaded.', 'Muat ulang kasus tersimpan': 'Reload saved case',
  'Ada perubahan yang belum disimpan. Simpan perubahan pada bahan kasus atau draft sebelum review dan ekspor.': 'You have unsaved changes. Save your case materials or draft before reviewing and exporting.',
  'Percakapan & cerita awal': 'Conversation & initial account', 'Tempel teks atau impor file TXT.': 'Paste text or import a TXT file.',
  'Contoh sintetis dimuat. Klik simpan untuk mencoba alur.': 'Synthetic example loaded. Save it to try the workflow.', 'Coba contoh': 'Try an example',
  'Judul kasus': 'Case title', 'Contoh: Tawaran pekerjaan daring dan top-up': 'Example: Online job offer and top-up requests', 'Format percakapan': 'Conversation format',
  'Teks biasa · satu pesan per baris': 'Plain text · one message per line', 'TXT · [DD/MM/YYYY, HH:MM] Pengirim: pesan': 'TXT · [DD/MM/YYYY, HH:MM] Sender: message',
  'Impor TXT': 'Import TXT', 'Membaca TXT…': 'Reading TXT…', 'TXT maksimal 32 KB.': 'TXT files must be no larger than 32 KB.', 'TXT harus menggunakan encoding UTF-8.': 'TXT files must use UTF-8 encoding.',
  'Percakapan': 'Conversation', 'Tempel pesan di sini.\nPengirim dan waktu tidak ditebak dari teks biasa.': 'Paste messages here.\nSender and time are not inferred from plain text.',
  'Teks biasa: satu pesan per baris, tanpa metadata pengirim/waktu. TXT berheader: format yang dipilih harus persis cocok; baris lanjutan menjadi bagian pesan. Maksimal 100 pesan / 32 KB.': 'Plain text: one message per line, without sender/time metadata. TXT with headers: the selected format must match exactly; continuation lines belong to the message. Maximum 100 messages / 32 KB.',
  'Ini data simulasi': 'This is simulated data', 'Catatan transfer': 'Transfer records', 'Isi berdasarkan catatan Anda, bukan hasil AI.': 'Enter details from your records, rather than AI output.',
  'catatan': 'records', 'Belum ada catatan transfer.': 'No transfer records yet.', 'Anda tetap dapat mulai dari percakapan.': 'You can still start with a conversation.',
  'Hapus catatan': 'Remove record', 'Bank pengirim': 'Sending bank', 'Rekening pengirim': 'Sending account', 'Bank penerima': 'Receiving bank', 'Rekening penerima': 'Receiving account',
  'Nominal (rupiah)': 'Amount (rupiah)', 'Kosong bila belum tahu': 'Leave blank if unknown', 'Status transfer': 'Transfer status', 'Tanggal transfer': 'Transfer date',
  'Boleh kosong. Waktu chat bukan tanggal transfer.': 'May be blank. Chat time is not the transfer date.', 'Catatan': 'Note', 'Koreksi, perbedaan keterangan, dll.': 'Corrections, conflicting accounts, etc.',
  '+ Tambah catatan transfer': '+ Add a transfer record', 'Menyimpan lampiran…': 'Saving attachment…', 'Mengunduh lampiran…': 'Downloading attachment…', 'Menghapus lampiran…': 'Removing attachment…',
  'Lampiran dihapus. Review perlu dikonfirmasi ulang.': 'Attachment removed. Review must be confirmed again.', 'Disimpan di PostgreSQL lokal.': 'Saved in local PostgreSQL.',
  'Menyimpan bahan kasus…': 'Saving case materials…', 'Simpan & lanjut': 'Save & continue', 'MULAI SECARA BERTAHAP': 'START STEP BY STEP',
  'Tidak harus lengkap sekarang.': 'You can fill in more later.',
  'Pertahankan pesan asli. Isi rincian transfer yang Anda ketahui, lalu tambahkan bukti bila tersedia.': 'Keep the original messages. Enter the transfer details you know, then add evidence if available.',
  'Teks percakapan atau ekspor TXT': 'Conversation text or TXT export', 'Nominal, rekening, dan status transfer': 'Amounts, accounts, and transfer status', 'Bukti transfer atau komunikasi': 'Transfer or communication evidence',
  'Jelas sebelum dikirim.': 'Know what will be sent.',
  'Menyimpan kasus belum mengirim teks ke AI. Pada langkah berikutnya, Anda memilih draft manual atau menyetujui pengiriman teks ke provider Apertus yang dikonfigurasi.': 'Saving a case does not send text to AI. In the next step, choose a manual draft or consent to sending text to the configured Apertus provider.',
  'Lampiran hanya sebagai referensi. Gambar dan PDF belum dibaca otomatis.': 'Attachments are for reference. Images and PDFs are not read automatically.',
  'Susun dengan Apertus': 'Draft with Apertus', 'Ringkasan per pesan, dengan ID dan kutipan asli.': 'A summary for each message, with its ID and original quote.',
  'AI membantu merapikan bahasa. Hasilnya bisa keliru dalam arah transfer atau makna janji. Seluruh draft akan Anda periksa pada langkah berikutnya.': 'AI helps organize the wording. It can misstate transfer direction or the meaning of a promise. Review the entire draft in the next step.',
  'pesan sumber': 'source messages', 'catatan transfer': 'transfer records', 'lampiran referensi': 'reference attachments',
  'Saya setuju mengirim': 'I consent to sending the', 'teks pesan': 'message text', 'kasus ini ke provider eksternal': 'of this case to the external provider', 'untuk disusun oleh Apertus.': 'for drafting with Apertus.',
  'Hanya ID sumber dan isi pesan dikirim. Isi pesan dapat memuat data pribadi; periksa sebelum menyetujui. Isian transfer dan file lampiran tidak dikirim. Pengiriman ini bukan laporan resmi.': 'Only source IDs and message contents are sent. Messages may contain personal data; check them before consenting. Transfer fields and attachment files are not sent. This is not an official report submission.',
  'Percakapan ini melebihi batas draft AI (20 pesan / 12 KB). Draft manual tetap tersedia.': 'This conversation exceeds the AI draft limit (20 messages / 12 KB). Manual drafting is still available.',
  'Apertus sedang menyusun draft… (maksimal 40 detik)': 'Apertus is drafting… (up to 40 seconds)', 'Susun ulang dengan Apertus': 'Redraft with Apertus', 'Susun draft dengan Apertus': 'Create draft with Apertus',
  'Susun secara manual': 'Draft manually', 'Mulai dari kutipan pesan asli. Anda dapat menyunting sendiri tanpa mengirim teks ke provider.': 'Start with the original message quotes. Edit them yourself without sending text to the provider.',
  'Menyiapkan draft manual…': 'Preparing manual draft…', 'Gunakan draft manual': 'Use a manual draft', 'Lanjutkan review draft yang tersimpan →': 'Continue reviewing the saved draft →',
  'Kanal resmi tetap tersedia.': 'Official channels remain available.',
  'Jika AI gagal atau merespons dengan model berbeda, Anda tetap bisa meninjau sumber dan membuat bahan laporan secara manual.': 'If AI fails or returns a different model, you can still review the sources and prepare report materials manually.',
  'Menuju IASC ↗': 'Go to IASC ↗', 'DRAFT MANUAL': 'MANUAL DRAFT', 'DRAFT APERTUS': 'APERTUS DRAFT',
  'entri · cocokkan makna, nominal dalam pesan, dan arah dana': 'entries · check meaning, amounts in messages, and direction of funds',
  'Menyimpan koreksi draft…': 'Saving draft edits…', 'Simpan koreksi': 'Save edits', 'DRAFT YANG DAPAT DIKOREKSI': 'EDITABLE DRAFT', 'PESAN SUMBER ASLI': 'ORIGINAL SOURCE MESSAGE',
  'Entri': 'Entry', 'Kronologi': 'Chronology', 'Saya telah mencocokkan entri ini dengan sumber': 'I have compared this entry with its source',
  'Pengirim belum diketahui': 'Sender unknown', 'Waktu pesan belum diketahui': 'Message time unknown', 'Kutipan utuh ·': 'Full quote ·', '· waktu pesan bukan waktu transfer': '· message time is not transfer time',
  'Periksa catatan transfer': 'Review transfer records', 'INPUT PENGGUNA': 'USER INPUT',
  'Saya telah memeriksa rekening, nominal, tanggal, dan status transfer.': 'I have checked the accounts, amounts, dates, and transfer status.',
  'Lengkapi atau tandai yang belum diketahui': 'Fill in or mark unknown details',
  'Pemeriksaan heuristik dapat melewatkan kesalahan atau menandai parafrasa yang benar. Tanda review bukan penentuan kebenaran.': 'Heuristic checks may miss errors or flag valid paraphrases. Review flags do not establish truth.',
  'Saya telah memeriksa tanda review dan data belum lengkap. Saya memahami ekspor ini hanya bahan bantu, bukan laporan resmi atau verifikasi bank.': 'I have checked the review flags and incomplete details. I understand this export is supporting material, not an official report or bank verification.',
  'entri diperiksa': 'entries reviewed', 'Menyimpan konfirmasi review…': 'Saving review confirmation…', 'Konfirmasi & lanjut': 'Confirm & continue',
  'Perubahan berarti review ulang.': 'Changes require a new review.',
  'Koreksi draft, input kasus, atau lampiran akan membatalkan konfirmasi sebelumnya. Ekspor hanya tersedia untuk revisi yang telah dikonfirmasi.': 'Edits to the draft, case input, or attachments clear the previous confirmation. Export is available only for a confirmed revision.',
  'Koreksi bahan kasus →': 'Edit case materials →', 'REVIEW TERSIMPAN': 'REVIEW SAVED', 'REVIEW DIPERLUKAN': 'REVIEW REQUIRED',
  'Bahan yang Anda periksa siap diunduh.': 'Your reviewed materials are ready to download.', 'Konfirmasikan revisi terbaru dahulu.': 'Confirm the latest revision first.',
  'Ekspor memuat kronologi, rincian transfer per pasangan rekening, serta daftar bukti. Nilai kosong tetap dicatat sebagai belum diketahui.': 'The export includes the chronology, transfers grouped by account pair, and an evidence list. Blank values remain unknown.',
  'Revisi': 'Revision', '· dikonfirmasi': '· confirmed', 'Menyiapkan paket ZIP…': 'Preparing ZIP package…', 'Unduh paket ZIP': 'Download ZIP package',
  'Menyiapkan kronologi…': 'Preparing chronology…', 'Unduh kronologi TXT': 'Download chronology TXT', 'Kembali ke review →': 'Back to review →',
  'Kronologi & sumber': 'Chronology & sources', 'Teks yang telah diperiksa + kutipan sumber': 'Reviewed text + source quotes', 'Rincian transfer': 'Transfer details',
  'CSV berdasarkan input Anda; status terpisah': 'CSV from your entries; separate statuses', 'Lampiran & snapshot kasus': 'Attachments & case snapshot', 'lampiran referensi + JSON revisi review': 'reference attachments + reviewed revision JSON',
  'LANGKAH SELANJUTNYA': 'NEXT STEP', 'Gunakan bahan ini pada kanal resmi.': 'Use these materials with official channels.',
  'LaporSiap tidak mengirim laporan, memblokir rekening, atau menjamin pengembalian dana. Kami tidak berafiliasi dengan OJK/IASC.': 'LaporSiap does not submit reports, block accounts, or guarantee refunds. We are not affiliated with OJK/IASC.',
  'Riwayat perubahan': 'Revision history', 'Draft awal dan koreksi disimpan dalam revisi terpisah.': 'The original draft and edits are saved as separate revisions.',
  'Memuat riwayat…': 'Loading history…', 'Lihat riwayat revisi →': 'View revision history →',
  'LaporSiap · Bahan bantu, dengan sumber dan review manusia.': 'LaporSiap · Supporting materials, with sources and human review.', 'Development lokal · Data tersimpan di perangkat ini': 'Local development · Data stored on this device',
  'CATATAN TRANSFER': 'TRANSFER RECORDS', 'Nominal transfer berhasil': 'Successful transfer amount', 'yang dilaporkan pengguna': 'reported by the user',
  'Hanya nominal yang diketahui. Bukan kerugian bersih atau nilai yang diverifikasi bank.': 'Known amounts only. Not net losses or bank-verified amounts.',
  'Total sebagian:': 'Partial total:', 'transfer berhasil belum memiliki nominal.': 'successful transfers have unknown amounts.',
  'Transfer dalam proses, gagal, dan belum dipastikan tidak masuk total.': 'Pending, failed, and unconfirmed transfers are excluded from the total.',
  'Tidak ada transfer yang diisi pengguna.': 'No transfers entered by the user.', 'Pasangan rekening': 'Account pair', 'Nominal & tanggal': 'Amount & date',
  'Bank belum diketahui': 'Bank unknown', 'Rekening belum diketahui': 'Account unknown', 'Tanggal transfer belum diketahui': 'Transfer date unknown',
  'Lampiran bukti': 'Evidence attachments', 'Referensi gambar atau PDF. Belum termasuk OCR.': 'Reference images or PDFs. OCR is not included.',
  'Pilih file bukti': 'Choose an evidence file', 'PNG, JPEG, PDF · maks. 5 MB/file · 10 lampiran': 'PNG, JPEG, PDF · max. 5 MB/file · 10 attachments',
  'Gambar maksimal 24 juta piksel.': 'Images must be no larger than 24 million pixels.',
  'Gambar rusak atau tidak lengkap. Pilih file PNG/JPEG yang dapat dibuka.': 'The image is damaged or incomplete. Choose a PNG/JPEG file that opens correctly.',
  'Dimensi gambar terlalu besar. Gunakan gambar maksimal 24 juta piksel.': 'The image dimensions are too large. Use an image no larger than 24 million pixels.',
  'PDF tidak lengkap atau struktur dasarnya tidak valid. Pilih PDF yang dapat dibuka.': 'The PDF is incomplete or its basic structure is invalid. Choose a PDF that opens correctly.',
  'KB · referensi': 'KB · reference', 'Unduh': 'Download', 'Hapus': 'Remove', 'Simpan bahan kasus dahulu untuk menambahkan lampiran.': 'Save your case materials before adding attachments.',
  // API validation and provider errors.
  'Backend belum terhubung. Periksa server lokal; kanal resmi tetap dapat diakses.': 'The backend is unavailable. Check the local server; official channels remain accessible.',
  'Unduhan gagal.': 'Download failed.', 'Database belum tersedia.': 'The database is unavailable.',
  'Input JSON tidak valid, memiliki field tambahan, atau terlalu besar.': 'The JSON input is invalid, contains extra fields, or is too large.', 'Input harus satu objek JSON.': 'Input must be a single JSON object.',
  'Host lokal tidak diizinkan.': 'This local host is not allowed.', 'Origin tidak diizinkan.': 'This origin is not allowed.',
  'Kasus tidak ditemukan atau akses tidak valid.': 'Case not found or access is invalid.', 'Kasus telah berubah. Muat ulang sebelum melanjutkan.': 'The case has changed. Reload it before continuing.',
  'Penyimpanan gagal. Coba kembali; jangan ekspor sebelum review tersimpan.': 'Saving failed. Try again; export only after the review has been saved.',
  'Persetujuan pengiriman teks ke provider eksternal diperlukan.': 'Consent is required to send text to the external provider.', 'AI belum tersedia. Gunakan draft manual.': 'AI is unavailable. Use a manual draft.',
  'Pilih Apertus atau manual.': 'Choose Apertus or manual.', 'Buat draft terlebih dahulu.': 'Create a draft first.',
  'Periksa sumber, transfer, dan data belum lengkap sebelum mengonfirmasi review.': 'Check the sources, transfers, and incomplete details before confirming your review.',
  'Unggah sebagai multipart/form-data.': 'Upload using multipart/form-data.', 'Satu file bukti diperlukan.': 'One evidence file is required.', 'Lampiran harus 1 byte–5 MB.': 'Attachments must be between 1 byte and 5 MB.',
  'Unggah satu file pada setiap permintaan.': 'Upload one file per request.', 'Bukti harus berupa PNG, JPEG, atau PDF. OCR belum tersedia.': 'Evidence must be a PNG, JPEG, or PDF. OCR is unavailable.',
  'Maksimal 10 lampiran per kasus.': 'Maximum 10 attachments per case.', 'Ekspor belum tersedia. Konfirmasikan review revisi terbaru.': 'Export is unavailable. Confirm your review of the latest revision.',
  'Isi percakapan UTF-8 wajib diisi (maksimal 32 KB).': 'Enter a UTF-8 conversation (maximum 32 KB).',
  'Tanggal/jam header TXT tidak valid. Pilih teks biasa untuk mempertahankan sumber tanpa metadata.': 'The TXT header date/time is invalid. Select plain text to preserve sources without metadata.',
  'Format TXT belum didukung. Gunakan [DD/MM/YYYY, HH:MM] Pengirim: pesan, atau pilih teks biasa.': 'This TXT format is unsupported. Use [DD/MM/YYYY, HH:MM] Sender: message, or select plain text.',
  'Header TXT ambigu. Pilih teks biasa agar metadata tidak ditebak.': 'The TXT header is ambiguous. Select plain text to avoid inferring metadata.',
  'Pilih format teks yang didukung.': 'Choose a supported text format.', 'Gunakan 1–100 pesan per kasus.': 'Use 1–100 messages per case.', 'Setiap pesan harus berisi teks (maksimal 6 KB).': 'Each message must contain text (maximum 6 KB).',
  'Judul wajib diisi, maksimal 160 karakter.': 'Enter a title (maximum 160 bytes).', 'Maksimal 100 catatan transfer.': 'Maximum 100 transfer records.', 'ID transfer harus unik dan tidak kosong.': 'Transfer IDs must be unique and nonempty.',
  'Status transfer tidak valid.': 'Invalid transfer status.', 'Nominal harus berupa rupiah bulat nonnegatif, atau kosong.': 'Amounts must be nonnegative whole rupiah values, or blank.',
  'Total transfer terlalu besar.': 'The transfer total is too large.', 'Detail transfer terlalu panjang.': 'Transfer details are too long.', 'Total melebihi batas angka yang dapat ditampilkan secara presisi.': 'The total exceeds the limit for precise number display.',
  'Total sebagian: ada transfer berhasil dengan nominal belum diketahui.': 'Partial total: some successful transfers have unknown amounts.',
  'Belum ada lampiran bukti. Lampiran dapat ditambahkan tanpa OCR.': 'No evidence attachments yet. Attachments can be added without OCR.',
  'Setiap pesan sumber wajib memiliki satu entri kronologi.': 'Each source message must have one chronology entry.', 'Sumber dan kutipan tidak dapat diubah melalui draft; ringkasan wajib diisi.': 'Sources and quotes cannot be changed through the draft; summaries are required.',
  'Apertus belum dikonfigurasi. Anda tetap dapat menyusun draft manual.': 'Apertus is not configured. You can still create a manual draft.',
  'Konfigurasi model harus Apertus 1.5 70B; fallback model lain tidak diizinkan.': 'The model must be Apertus 1.5 70B; fallback to another model is not allowed.',
  'Endpoint Apertus tidak valid.': 'Invalid Apertus endpoint.', 'Draft AI dibatasi 20 pesan / 12 KB untuk menjaga kutipan lengkap. Gunakan draft manual untuk kasus lebih panjang.': 'AI drafts are limited to 20 messages / 12 KB to preserve full quotes. Use a manual draft for longer cases.',
  'Tidak dapat menyiapkan permintaan AI.': 'Unable to prepare the AI request.', 'Apertus tidak merespons dalam batas waktu atau koneksi gagal. Coba lagi atau gunakan draft manual.': 'Apertus timed out or the connection failed. Try again or use a manual draft.',
  'Respons AI terlalu besar atau gagal dibaca.': 'The AI response is too large or could not be read.', 'Respons provider tidak valid.': 'Invalid provider response.',
  'Identitas model respons bukan Apertus yang diminta. Draft ditolak; tidak ada fallback otomatis.': 'The response model is not the requested Apertus model. The draft was rejected; there is no automatic fallback.',
  'Respons AI belum lengkap. Gunakan draft manual atau coba kembali.': 'The AI response is incomplete. Use a manual draft or try again.',
  'Draft AI tidak memenuhi schema atau kutipan sumber. Draft ditolak; data asli tetap tersedia.': 'The AI draft failed schema or source quote validation. The draft was rejected; original data remains available.',
}

export function translate(text: string, language: Language): string {
  if (language === 'id') return text
  if (english[text]) return english[text]
  let match = text.match(/^Revisi (\d+) · Perlu review$/)
  if (match) return `Revision ${match[1]} · Review required`
  match = text.match(/^Ada istilah di luar kutipan: (.*)\. Periksa maknanya\.$/)
  if (match) return `Term outside the quote: ${match[1]}. Check its meaning.`
  match = text.match(/^Transfer (\d+): (.+) belum lengkap\.$/)
  if (match) {
    const fields: Record<string, string> = { nominal: 'amount', 'tanggal transfer': 'transfer date', 'rekening pengirim': 'sending account', 'rekening penerima': 'receiving account', status: 'status' }
    return `Transfer ${match[1]}: missing ${match[2].split(', ').map(field => fields[field] || field).join(', ')}.`
  }
  match = text.match(/^Transfer (\d+) memiliki rincian yang sama dengan catatan lain; periksa duplikasi\.$/)
  if (match) return `Transfer ${match[1]} has the same details as another record; check for duplication.`
  match = text.match(/^Tanggal transfer (\d+) harus YYYY-MM-DD atau kosong\.$/)
  if (match) return `Transfer ${match[1]} date must be YYYY-MM-DD or blank.`
  match = text.match(/^Provider AI mengembalikan HTTP (\d+)\. Sumber tetap tersimpan; draft manual tersedia\.$/)
  if (match) return `The AI provider returned HTTP ${match[1]}. Sources remain saved; manual drafting is available.`
  match = text.match(/^Permintaan gagal \((\d+)\)\.$/)
  if (match) return `Request failed (${match[1]}).`
  return text
}

export function LanguageProvider({ children, preferences }: { children: ReactNode; preferences: LanguagePreferences }) {
 const language = useAppSelector(state => state.locale.language)
 useEffect(() => { preferences.write(language); document.documentElement.lang = language; document.title = language === 'en' ? 'LaporSiap — Prepare report materials' : 'LaporSiap — Siapkan bahan laporan' }, [language, preferences])
 return <>{children}</>
}
export function useLanguage() {
 const language = useAppSelector(state => state.locale.language)
 const dispatch = useAppDispatch()
 return { language, setLanguage: (language: Language) => dispatch(localeSlice.actions.setLanguage(language)), tr: (text: string) => translate(text, language) }
}
