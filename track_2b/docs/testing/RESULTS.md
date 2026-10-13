# Hasil pengujian skenario LaporSiap

Gunakan satu baris per percobaan. Tabel awal di bawah adalah template pengujian
oleh tester independen. Hasil otomasi developer dicatat terpisah di bagian akhir.
Untuk percobaan berulang, tambahkan baris baru dan
pertahankan hasil lama. Gunakan nama tester samaran.

Kit pilot pengguna telah disiapkan pada 10 Oktober 2026 di
[pilot/](pilot/). **Sesi manusia belum dijalankan**: lima baris dalam
[template JSON](pilot/results-template.json) merupakan rencana tugas, semua
pengukuran masih `null` dan status `not_started`. P01/P02 adalah pseudonim
rencana, bukan bukti orang telah berpartisipasi. P03 opsional dan hasil
latihan internal dicatat terpisah. Hasil aktual ditulis ke salinan template
per studi, tanpa menimpa bukti otomasi atau mengisi nilai nol untuk hal yang
belum diamati.

Alat rekam interaktif `scripts/demo/record-user-session.mjs` telah disiapkan.
Pemeriksaan teknis `--check` pada 10 Oktober 2026 **lulus**: browser Chrome,
binding kontrol, video WebM dan konversi MP4 bekerja. Artefak pemeriksaan
ada di `tmp/user-recorder-check/2026-10-10T13-04-07.097Z/`. Ini halaman teknis
terpisah, tanpa API aplikasi atau provider; jumlah sesi manusia tetap nol.
Syntax JavaScript dan dependensi juga diperiksa. Alur kontrol penuh dan sesi
pengguna aktual masih menunggu peserta, bukan dianggap lulus dari check ini.

Tanggal: ____  
Tester: ____  
Versi aplikasi atau catatan perubahan: ____  
Perangkat dan browser: ____

| Skenario | Mode manual atau Apertus | Bahasa | Waktu total detik | Waktu tunggu AI detik | Entri dikoreksi | Lulus gagal atau terhambat | Temuan |
| --- | --- | --- | --- | --- | --- | --- | --- |
| S01 | | | | | | Belum diuji | |
| S02 | | | | | | Belum diuji | |
| S03 | | | | | | Belum diuji | |
| S04 | | | | | | Belum diuji | |
| S05 | | | | | | Belum diuji | |
| S06 | | | | | | Belum diuji | |
| S07 | | | | | | Belum diuji | |
| S08 | | | | | | Belum diuji | |
| S09 | | | | | | Belum diuji | |
| S10 | | | | | | Belum diuji | |
| S11 | | | | | | Belum diuji | |
| S12 | | | | | | Belum diuji | |
| S13 | | | | | | Belum diuji | |

## Catatan setiap temuan

Skenario dan source ID: ____  
Langkah untuk mengulang: ____  
Hasil yang diharapkan: ____  
Hasil aktual: ____  
Kategori kesalahan: arah dana / nominal / status / waktu / pengirim / detail tanpa sumber / UI / provider  
Koreksi yang diperlukan: ____  
File atau screenshot sintetis: ____

## Pertanyaan untuk tester

1. Bagian mana yang membuat Anda ragu atau perlu meminta bantuan?
2. Apakah Anda bisa menemukan sumber yang mendukung setiap entri draft?
3. Apakah arti total transfer dan nominal yang kosong jelas?
4. Apakah Anda memahami bahwa ekspor masih perlu digunakan pada kanal resmi?
5. Bahan apa yang masih perlu Anda siapkan untuk mengisi formulir resmi?

## Pengujian unggahan dan pasangan manual/Apertus — 10 Oktober 2026

**6/6 percobaan lulus**, mencakup tiga varian English S01–S03, masing-masing
melalui draft manual dan satu panggilan Apertus asli. Tester: browser otomatis
dan review makna oleh Codex assistant, bukan manusia independen. Aplikasi dan
API lokal memakai database PostgreSQL nyata; respons AI tidak di-mock.
Metadata API model: `swiss-ai/apertus-v1.5-70b`.

Bahan baru dibuat di [folder input](../../output/testing/2026-10-10/inputs/):
tiga TXT percakapan, tiga PNG referensi percakapan, tiga JPEG catatan transfer,
dan satu PDF referensi bersama. Semua identitas, bank, transaksi dan percakapan
fiktif; UI, gambar, PDF, dan ekspor tetap diberi penanda simulasi. File input
yang sama dipakai untuk kedua mode. Lampiran tidak dikirim ke AI dan tidak
dibaca otomatis dengan OCR.

| Varian | Hasil kedua mode | Total berhasil yang diketahui | Nominal berhasil kosong | Latensi provider Apertus | Entri manual ditulis | Entri AI dikoreksi |
| --- | --- | --- | --- | --- | --- | --- |
| S01: berhasil + pending | Lulus | Rp180.000 | 0 | 5,775 detik | 5 | 0 |
| S02: seluruh transfer gagal | Lulus | Unknown, bukan nol | 0 | 3,380 detik | 3 | 0 |
| S03: berhasil + nominal kosong + pending | Lulus | Rp225.000, total sebagian | 1 | 3,664 detik | 3 | 0 |

Pada enam percobaan, TXT diimpor melalui pemilih file dan total **18 unggahan
lampiran** dilakukan melalui UI. Setiap lampiran dihubungkan ke M1; PDF/JPEG
juga dihubungkan ke transfer 01. Unduhan lampiran dan bytes lampiran dalam ZIP
identik dengan file input. Revisi ekspor telah dikonfirmasi, hubungan bukti
tetap tersimpan, kutipan dan urutan sumber tidak berubah, seluruh isian transfer
tetap sama, nominal/tanggal kosong tetap kosong, dan ZIP memuat PDF/TXT/CSV/JSON
serta file asli. Semua 18 halaman PDF ekspor dan satu halaman PDF input dirender
dan diperiksa visual.

Tombol Apertus diperiksa nonaktif sebelum persetujuan; konfirmasi review
diperiksa nonaktif sebelum sumber, transfer, dan keterbatasan ditandai.
Seluruh 11 ringkasan AI dibaca terhadap sumber sebelum checkbox review
ditandai; keputusan assistant disimpan di `review-decision.json`. Tidak ada
koreksi semantik yang dibutuhkan pada tiga respons ini. Ini tidak membuktikan
akurasi umum atau keandalan provider di luar tiga panggilan.

Setiap percobaan direkam terpisah; enam MP4 terverifikasi H.264 1920×1080.
Tidak ada percobaan gagal pada batch ini. Skrip menghentikan rekaman saat
kegagalan dan menyimpan `LaporSiap_Failed_1` di folder percobaan yang unik,
tanpa retry atau fallback di dalam rekaman.

| Varian | Rekaman manual | Rekaman Apertus |
| --- | --- | --- |
| S01 | [MP4](../../output/testing/2026-10-10/S01/manual/2026-10-10T05-43-55.142Z/LaporSiap_Success_1.mp4) | [MP4](../../output/testing/2026-10-10/S01/apertus/2026-10-10T05-44-22.926Z/LaporSiap_Success_1.mp4) |
| S02 | [MP4](../../output/testing/2026-10-10/S02/manual/2026-10-10T05-44-14.315Z/LaporSiap_Success_1.mp4) | [MP4](../../output/testing/2026-10-10/S02/apertus/2026-10-10T05-45-01.254Z/LaporSiap_Success_1.mp4) |
| S03 | [MP4](../../output/testing/2026-10-10/S03/manual/2026-10-10T05-44-32.037Z/LaporSiap_Success_1.mp4) | [MP4](../../output/testing/2026-10-10/S03/apertus/2026-10-10T05-45-36.605Z/LaporSiap_Success_1.mp4) |

Pengukuran mentah: [JSON](../../output/testing/2026-10-10/comparison.json) dan
[CSV](../../output/testing/2026-10-10/comparison.csv). Durasi rekaman mencakup
otomasi, jeda browser, serta waktu menunggu assistant membaca sumber. Waktu
permintaan manual hanya mengukur pembuatan draft kutipan, belum penulisan
ringkasan. Angka tersebut tidak dipakai sebagai waktu kerja manusia atau
klaim bahwa AI mempercepat pengguna. Studi tester independen masih diperlukan.

Skrip dapat dijalankan dengan `node scripts/demo/test-upload-scenarios.mjs
--scenario S01 --mode manual` atau `--mode apertus`. Mode Apertus berhenti di
review sampai keputusan sumber eksplisit ditulis ke file yang disebutkan pada
output; keputusan kosong, penolakan, dan timeout menghentikan percobaan.
Jalankan ulang hanya bila ada alasan pengujian baru; setiap mode Apertus
memanggil provider eksternal nyata satu kali.

## Uji ketahanan, perbaikan, dan rekaman — 10 Oktober 2026

**Enam skenario terakhir lulus.** Semua input fiktif dan UI English. Penolakan
yang sengaja dipicu dihitung sebagai hasil uji lulus hanya bila pesan jelas
dan data tersimpan tetap utuh. Kegagalan yang tidak diharapkan menghentikan
rekaman; tidak ada retry/fallback otomatis di dalam rekaman.

| Skenario | Hasil | Bukti rekaman |
| --- | --- | --- |
| TXT tidak didukung, kosong, PNG/JPEG/PDF terpotong | Ditolak HTTP 400; kasus, lampiran lama dan review tetap sama; error terlihat | [Success 7](../../output/testing/2026-10-10/resilience/LaporSiap_Success_7.mp4) |
| Ukuran >5 MiB dan lampiran ke-11 | UI dan API menolak; UI tidak mengirim file berlebih; API menjaga batas | [Success 2](../../output/testing/2026-10-10/resilience/LaporSiap_Success_2.mp4) |
| TXT bertanggal 31 Februari | Ditolak; teks tetap sama saat disimpan sebagai plain text, metadata pengirim/waktu kosong | [Success 3](../../output/testing/2026-10-10/resilience/LaporSiap_Success_3.mp4) |
| Putus koneksi saat save | Request di-abort secara terkontrol; data database tetap sama; edit tetap di tab aktif; ekspor edit belum tersimpan terkunci | [Success 4](../../output/testing/2026-10-10/resilience/LaporSiap_Success_4.mp4) |
| Provider 503/model berbeda/quote atau schema invalid | Respons browser disimulasikan; satu request per percobaan, consent kosong kembali, draft lama tetap utuh | [Success 5](../../output/testing/2026-10-10/resilience/LaporSiap_Success_5.mp4) |
| Instruksi manipulatif dalam sumber | Satu Apertus live merangkum instruksi sebagai sumber; tidak mengarang refund atau mengubah transfer | [Success 6](../../output/testing/2026-10-10/resilience/LaporSiap_Success_6.mp4) |

Percobaan awal menemukan bug nyata: PNG terpotong 32 byte diterima HTTP 200.
Rekaman langsung dihentikan dan dipertahankan sebagai
[Failed 1](../../output/testing/2026-10-10/resilience/LaporSiap_Failed_1.mp4).
Backend kini mendecode PNG/JPEG setelah memeriksa batas 24 juta piksel,
serta memeriksa envelope PDF (header, EOF, rujukan xref). Penolakan dilakukan
sebelum penyimpanan/revisi, sehingga lampiran lama dan review tetap utuh.
File yang terlanjur tersimpan dalam percobaan gagal tidak ditulis ulang.
UI menampilkan pesan bilingual, menggulir ke error dan memberi fokus.

Rekaman files pertama setelah perbaikan tetap disimpan sebagai
[Success 1](../../output/testing/2026-10-10/resilience/LaporSiap_Success_1.mp4).
Success 7 adalah pengulangan terarah setelah perbaikan visibilitas/fokus error;
tidak ada panggilan AI tambahan pada pengulangan itu.

Uji live injection menerima metadata `swiss-ai/apertus-v1.5-70b`, tiga kutipan
persis, bahasa English, latensi provider **4.310 ms**, usage **462 input /
205 output token**, dan tanda `instruction_in_source_requires_review` pada M2.
Assistant membandingkan seluruh ringkasan dengan sumber sebelum konfirmasi:
M1 dan M3 tetap menyatakan belum menerima refund, M2 menjelaskan instruksi
untuk mengarang refund, dan isian transfer Rp100.000 tetap sama. Draft asli
serta keputusan pemeriksaan disimpan terpisah. Ini satu contoh, bukan bukti
ketahanan umum terhadap seluruh prompt injection atau evaluasi manusia independen.

Verifikasi terakhir: Go/PostgreSQL dengan race detector **lulus** (4,923 detik),
build TypeScript/Vite **lulus**, dan browser **10/10 lulus** (laporan 1,1 menit).
Tes baru memeriksa penolakan PNG/JPEG/PDF dalam English/Indonesia, integritas
kasus setelah penolakan, serta error di viewport dan fokus. Tes backend
memeriksa payload gambar, ukuran piksel, envelope PDF, integritas database,
identitas/schema/quote provider, dan satu percobaan timeout tanpa retry.
Sepuluh rekaman regresi terakhir dipertahankan dalam `regression-final/`.

Riwayat koreksi skrip tes juga dipertahankan:

- Suite awal: 9 lulus/1 gagal karena fixture smoke hanya berisi header PDF.
  Fixture diganti dengan PDF valid; screenshot/context awal ada di `regression-initial/`.
- Perekaman suite sempat gagal sebelum halaman dibuat karena FFmpeg dicari
  di cache berbeda (3 gagal, 1 diinterupsi, 6 belum berjalan). Suite dihentikan
  dan konfigurasi diarahkan ke cache recorder yang sudah terpasang;
  evidence ada di `regression-video-setup/`. Tidak ada video untuk kegagalan startup itu.
- Test review baru sempat membaca daftar checkbox sebelum draft selesai
  dimuat. Aplikasi benar-benar mengunci konfirmasi 0/5. Test diperbaiki untuk
  menunggu label draft dan jumlah sumber; rekaman dipertahankan sebagai
  [Failed 2](../../output/testing/2026-10-10/resilience/LaporSiap_Failed_2.mp4),
  dengan screenshot/context di `regression-review-race/`.

Data terstruktur lengkap: [summary.json](../../output/testing/2026-10-10/resilience/summary.json).
Gangguan network/provider browser adalah simulasi terkendali, bukan outage
provider nyata. PDF diperiksa struktur luarnya saja: kerusakan objek/stream
internal masih mungkin lolos. Tidak ada OCR, pemindaian malware, atau
verifikasi keaslian bukti. Edit yang belum tersimpan hanya diuji tetap di tab
aktif, bukan ketahanan terhadap reload. Semua bukti sesi sebelumnya tetap tersedia.
