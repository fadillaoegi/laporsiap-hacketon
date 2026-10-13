# Skrip demo LaporSiap dalam bahasa Inggris dan Indonesia

Demo ini menjelaskan latar belakang LaporSiap, alur empat tahap, fitur utama,
dan peran Apertus dalam menyusun draft kronologi yang perlu diperiksa manusia.
Gunakan contoh sintetis bawaan aplikasi. English merupakan bahasa antarmuka
yang digunakan pada arahan perekaman; versi Indonesia di bawah memiliki
makna yang sama untuk membantu memahami narasi.

File [VOICEOVER-EN.txt](VOICEOVER-EN.txt) dan [VOICEOVER-ID.txt](VOICEOVER-ID.txt)
berisi narasi bersih tanpa judul atau arahan layar, siap disalin ke alat
voice-over. Angka rupiah dan nama model ditulis sebagai kata agar mudah
dibacakan. Baca LaporSiap sebagai “Lapor Siap”; Public AI sebagai “Public
A I”, dan IASC sebagai “I A S C”.

Narasi English berjumlah sekitar **634 kata** dan Indonesia **586 kata**. Pada tempo 125–145 kata per menit, narasi English diperkirakan sekitar **4.4–5.1 menit**, sebelum tambahan jeda untuk demonstrasi.

## Alur perekaman

| Bagian | Fokus tampilan |
| --- | --- |
| 1 | Latar masalah dan contoh sintetis |
| 2 | Empat tahap dan akses langsung ke IASC |
| 3 | Percakapan, transfer, dan lampiran |
| 4 | Total berhasil dan status pending |
| 5 | Persetujuan serta draft Apertus dengan sumber |
| 6 | Pemeriksaan makna dan pilihan manual |
| 7 | Koreksi, progress, konfirmasi, dan riwayat |
| 8 | Ekspor, dua bahasa, dan batas prototipe |

Simpan lampiran sintetis sebelum membuat draft utama. Jika setelah merekam
bagian review Anda mengubah data untuk menunjukkan pembatalan konfirmasi,
periksa dan konfirmasikan ulang sebelum merekam unduhan. Pembukaan demo
memakai tombol **Try an example**, yang memuat lima pesan dan dua transfer.

Bagian Apertus menggunakan respons yang benar-benar diperoleh. Jika provider
gagal, gunakan kalimat pengganti berikut sambil memperlihatkan pesan error
serta tombol **Use a manual draft**; jangan memberi label APERTUS DRAFT pada
draft manual.

English: “The AI provider is unavailable in this attempt. The sources are still saved, and we can continue with a manual draft. The Apertus integration is designed to produce source-linked drafts when a valid response is available.”

Indonesia: “Provider AI belum tersedia pada percobaan ini. Sumber tetap tersimpan dan kita dapat melanjutkan dengan draft manual. Integrasi Apertus dirancang untuk menghasilkan draft yang terhubung ke sumber ketika respons valid tersedia.”

Jika menggunakan hasil integrasi yang disimpan di
[apertus-integration-20261009.json](../evidence/apertus-integration-20261009.json),
beri keterangan **Recorded integration result / Hasil integrasi tersimpan**.
Hasil tersebut memakai tiga pesan tentang transfer gagal, berbeda dari
contoh utama lima pesan; jangan menyambungkannya seolah respons untuk kasus
utama. Jika waktu tunggu dipotong saat editing, beri keterangan
**Generation time shortened / Waktu penyusunan dipersingkat**.

## Bagian 1 Latar belakang LaporSiap

Arahan layar: Tampilkan logo LaporSiap, lalu contoh percakapan sintetis. Sorot tawaran tugas, janji komisi, dan permintaan top-up. Jangan menampilkan data korban nyata.

### English voice over

A suspicious online job offer can leave someone with a confusing record of messages, promises, and top-up requests. Before filling in a report, they need to explain what happened and organize their transfer details and evidence.

LaporSiap was created around this preparation task. It is designed for the person preparing a report, or a family member helping them. In this demonstration, every message, account, and transaction is synthetic.

### Narasi Indonesia

Tawaran pekerjaan daring yang mencurigakan dapat meninggalkan rangkaian pesan, janji komisi, dan permintaan top-up yang membingungkan. Sebelum mengisi laporan, seseorang perlu menjelaskan apa yang terjadi serta merapikan rincian transfer dan bukti yang dimiliki.

LaporSiap dibuat untuk membantu tahap persiapan ini. Aplikasi ini ditujukan bagi orang yang menyiapkan laporan, maupun anggota keluarga yang membantunya. Dalam demonstrasi ini, seluruh pesan, rekening, dan transaksi adalah data sintetis.

## Bagian 2 Alur aplikasi dan kanal resmi

Arahan layar: Tampilkan halaman awal dalam English. Sorot empat tahap di sidebar dan tautan Open the official IASC channel. Jangan mengirim laporan ke portal resmi dalam demo.

### English voice over

The workflow has four steps: collect case materials, create a draft, review and edit, then export.

The official IASC link is available from the beginning. Someone who is ready to report can go directly to the official channel without waiting for an AI draft.

### Narasi Indonesia

Alur aplikasi terdiri dari empat tahap: mengumpulkan bahan kasus, menyusun draft, memeriksa dan mengoreksi, lalu mengekspor bahan.

Tautan resmi IASC tersedia sejak halaman awal. Pengguna yang sudah siap melapor dapat langsung menuju kanal resmi tanpa menunggu draft AI.

## Bagian 3 Input percakapan transfer dan lampiran

Arahan layar: Klik Try an example. Perlihatkan Conversation format, teks contoh, dua catatan transfer, dan tanggal yang kosong. Klik Save & continue; kembali ke Case materials untuk mengunggah lampiran sintetis PNG/JPEG/PDF.

### English voice over

First, we paste a conversation or import a TXT file. The app supports plain text and a specific header format containing the sender and message time. Missing metadata is left unknown.

We enter transfer details ourselves: the sending and receiving accounts, amount, date, and status. Unknown amounts and dates can remain blank.

We can also attach images or PDFs as references. The current version does not read them automatically with OCR.

### Narasi Indonesia

Pertama, kita menempelkan percakapan atau mengimpor file TXT. Aplikasi mendukung teks biasa dan format header tertentu yang memuat pengirim serta waktu pesan. Metadata yang tidak tersedia tetap ditandai belum diketahui.

Rincian transfer kita isi sendiri: rekening pengirim dan penerima, nominal, tanggal, serta status. Nominal dan tanggal yang belum diketahui boleh tetap kosong.

Kita juga dapat melampirkan gambar atau PDF sebagai referensi. Versi saat ini belum membacanya secara otomatis menggunakan OCR.

## Bagian 4 Perhitungan total transfer

Arahan layar: Sorot dua status pada contoh: Successful Rp180.000 dan Pending Rp420.000. Tunjukkan kartu Successful transfer amount reported by the user dengan total Rp180.000 atau IDR180,000 sesuai bahasa.

### English voice over

In this example, one hundred and eighty thousand rupiah is marked successful, while another four hundred and twenty thousand is pending.

The backend includes only successful transfers with known amounts, so the displayed total is one hundred and eighty thousand rupiah. This is the successful amount reported by the user, not a bank-verified loss.

Pending, failed, and unconfirmed transfers are excluded. Unknown values are not converted to zero.

### Narasi Indonesia

Dalam contoh ini, seratus delapan puluh ribu rupiah ditandai berhasil, sementara empat ratus dua puluh ribu rupiah lainnya masih dalam proses.

Backend hanya menjumlahkan transfer berhasil dengan nominal yang diketahui. Karena itu, total yang ditampilkan adalah seratus delapan puluh ribu rupiah. Nilai ini merupakan nominal transfer berhasil yang dilaporkan pengguna, bukan kerugian yang telah diverifikasi bank.

Transfer dalam proses, gagal, dan belum dipastikan tidak dihitung. Nilai yang belum diketahui tidak diubah menjadi nol.

## Bagian 5 Peran Apertus dalam penyusunan draft

Arahan layar: Buka Create draft. Tunjukkan consent Public AI dan tombol Create draft with Apertus. Jika merekam permintaan live, centang consent lalu kirim satu contoh sintetis. Setelah berhasil, tampilkan APERTUS DRAFT dan pasangan draft/source. Waktu tunggu dapat dipotong dengan keterangan pada video.

### English voice over

Apertus version one point five is a language model with seventy billion parameters. Here, it helps turn the message text into a concise chronology draft. Each entry is linked to a source ID and an exact original quote, so the user can check where it came from.

Before sending text to Public AI, the external provider, the app requests consent. It sends only source IDs and message text; that text may contain personal information. Transfer fields and attachment files are not sent.

Apertus does not create or change transfer records. The backend checks the response's reported model identity, structure, and source quotes before saving the draft.

### Narasi Indonesia

Apertus versi satu koma lima, dengan tujuh puluh miliar parameter, membantu mengubah teks pesan menjadi draft kronologi yang ringkas. Setiap entri terhubung ke ID sumber dan kutipan asli yang persis, sehingga pengguna dapat memeriksa asalnya.

Sebelum mengirim teks ke Public AI, sebagai provider eksternal, aplikasi meminta persetujuan. Yang dikirim hanya ID sumber dan teks pesan; teks tersebut dapat memuat informasi pribadi. Isian transfer dan file lampiran tidak dikirim.

Apertus tidak membuat atau mengubah catatan transfer. Backend memeriksa identitas model yang dilaporkan respons, struktur, dan kutipan sumber sebelum menyimpan draft.

## Bagian 6 Alasan review manusia tetap diperlukan

Arahan layar: Perbesar pasangan draft/source untuk pesan Saya baru mengirim Rp180.000. Tunjukkan cara memeriksa arah dana. Jika contoh aktual memiliki tanda heuristik, sorot tanda tersebut; jangan merekayasa output model agar terlihat sebagai hasil asli.

### English voice over

An exact quote does not guarantee an accurate summary. In early testing, one summary changed sending money into receiving money, even though its quote was correct.

That finding shaped the review process. The app displays the draft beside its source and flags some potential problems, including reversed transfer direction and instructions embedded in messages.

These checks can miss mistakes or flag valid wording. The user still reviews every entry. A manual draft is also available if AI fails or the user prefers not to send text.

### Narasi Indonesia

Kutipan yang persis belum menjamin ringkasannya akurat. Pada pengujian awal, satu ringkasan mengubah makna mengirim uang menjadi menerima uang, meskipun kutipannya benar.

Temuan itu membentuk proses review. Aplikasi menampilkan draft berdampingan dengan sumbernya dan menandai beberapa potensi masalah, termasuk arah transfer yang terbalik serta instruksi yang disisipkan dalam pesan.

Pemeriksaan ini dapat melewatkan kesalahan atau menandai kalimat yang sebenarnya benar. Pengguna tetap memeriksa setiap entri. Draft manual juga tersedia ketika AI gagal atau pengguna memilih tidak mengirim teks.

## Bagian 7 Koreksi konfirmasi dan riwayat revisi

Arahan layar: Sunting satu description sesuai sumber dan klik Save edits. Centang tiap sumber, review transfer, dan data belum lengkap. Sorot progress, lalu Confirm & continue. Di halaman ekspor, buka View revision history. Tunjukkan bahwa perubahan sesudah review mengunci ekspor kembali.

### English voice over

During review, we can edit the descriptions while keeping the original quotes intact. We check the transfer details, incomplete information, and each source message before confirming the revision.

The progress indicator shows how many sources have been reviewed. Revision history keeps the original draft and subsequent changes.

If we later change the case materials, draft, or attachments, the previous confirmation is cleared. Export becomes available again only after reviewing and confirming the updated revision.

### Narasi Indonesia

Saat review, kita dapat menyunting deskripsi sambil mempertahankan kutipan asli. Kita memeriksa rincian transfer, informasi yang belum lengkap, dan setiap pesan sumber sebelum mengonfirmasi revisi.

Indikator progress menunjukkan jumlah sumber yang sudah diperiksa. Riwayat revisi menyimpan draft awal serta perubahan berikutnya.

Jika kita kemudian mengubah bahan kasus, draft, atau lampiran, konfirmasi sebelumnya dibatalkan. Ekspor kembali tersedia setelah revisi terbaru diperiksa dan dikonfirmasi.

## Bagian 8 Ekspor dua bahasa dan batas prototipe

Arahan layar: Setelah review terbaru dikonfirmasi, unduh TXT dan ZIP. Tampilkan isi ZIP tanpa token atau credential. Ganti English ke Indonesia, tunjukkan label berubah dan sumber tetap asli. Akhiri pada halaman ekspor dengan tautan IASC.

### English voice over

After confirmation, we download a text chronology or a ZIP package containing the chronology, original conversation, transfer CSV, reviewed case snapshot, and attachments. Transfers are grouped by account pair in the chronology.

The interface supports English and Indonesian, with English as the default. Export labels follow the selected language, while the original evidence and user-written content remain unchanged.

The prototype uses a Go backend, a React and TypeScript interface, and PostgreSQL, running locally. User benefits still need testing.

LaporSiap prepares supporting materials. It does not submit reports, verify transactions, or guarantee refunds, and it is not affiliated with OJK or IASC. Its goal is to help people organize their account of events with sources they can review.

### Narasi Indonesia

Setelah konfirmasi, kita mengunduh kronologi teks atau paket ZIP yang berisi kronologi, percakapan asli, CSV transfer, snapshot kasus yang ditinjau, dan lampiran. Dalam kronologi, transfer dikelompokkan berdasarkan pasangan rekening.

Antarmuka mendukung bahasa Inggris dan Indonesia, dengan bahasa Inggris sebagai default. Label ekspor mengikuti bahasa yang dipilih, sementara bukti asli dan isi yang ditulis pengguna tetap dipertahankan.

Prototipe ini memakai backend Go, antarmuka React dan TypeScript, serta PostgreSQL, dan dijalankan secara lokal. Manfaatnya bagi pengguna masih perlu diuji.

LaporSiap menyiapkan bahan bantu. Aplikasi ini tidak mengirim laporan, memverifikasi transaksi, atau menjamin pengembalian dana, serta tidak berafiliasi dengan OJK maupun IASC. Tujuannya adalah membantu pengguna merapikan cerita kejadian dengan sumber yang dapat mereka periksa.

## Acuan isi demo

Fitur dan batas implementasi mengikuti [README](../../README.md),
[catatan teknis](../TECHNICAL.md), dan
[bukti integrasi Apertus](../evidence/apertus-integration-20261009.json).
Latar desain serta temuan arah transfer tercatat dalam
[handoff riset](../../../research/LAPORSIAP-HANDOFF.md).
IASC merupakan kanal resmi yang dapat dibuka melalui
[portal IASC](https://iasc.ojk.go.id/). Tautan acuan ini tidak perlu dibacakan
pada voice-over.
