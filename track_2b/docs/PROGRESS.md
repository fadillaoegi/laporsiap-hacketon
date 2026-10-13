# Progress implementasi — 9 Oktober 2026

## Selesai pada sesi awal

- Struktur `laporsiap_BE` dan `laporsiap_FE`, repo Git lokal, env diabaikan Git.
- PostgreSQL 18-alpine, volume khusus `laporsiap_laporsiap_pgdata`, database
  aplikasi dan `laporsiap_test`, koneksi/migrasi terverifikasi.
- Go API untuk input, parser terbatas, transfer, sumber, draft, review,
  lampiran, riwayat, serta ekspor TXT/ZIP.
- Token per kasus, kontrol akses seluruh endpoint kasus, optimistic revision,
  review invalidation, serta export snapshot konsisten.
- Adapter Apertus memakai credential yang sudah diizinkan dalam handoff,
  schema/quote/model guard, persetujuan pengiriman teks, dan jalur manual.
- React/TypeScript UI empat langkah, sumber/draft berdampingan, data kosong,
  tanda review, demo sintetis, kanal resmi dari awal, dan tampilan mobile.
- Font lokal dengan lisensi OFL, README, catatan teknis, screenshot demo.

## Hasil pemeriksaan

- `go test -race ./...` dengan `TEST_DATABASE_URL`: lulus, termasuk PostgreSQL.
- `pnpm build`: lulus TypeScript dan build Vite.
- Playwright Chrome: alur manual lengkap, upload referensi PDF sintetis,
  koreksi draft, konfirmasi, unduh ZIP, reload, ekspor terkunci setelah edit,
  dan mobile tanpa horizontal overflow lulus; tidak ada page error.
- Satu probe live **Go** melalui API aplikasi: metadata
  `swiss-ai/apertus-v1.5-70b`, tiga entri/quote cocok; isi kasus N03 mempertahankan
  transfer gagal dan belum adanya transaksi berhasil. Tidak ada review/ekspor
  otomatis; tidak mengulang batch dan tidak mengklaim akurasi 100%.

Screenshot sintetis di `screenshots/`: halaman input, review desktop, ekspor,
dan review mobile. Paket browser demo hanya di `tmp/` yang diabaikan Git.
Uji browser menyimpan kasus sintetis pada database aplikasi; uji Go memakai
database uji dan membersihkan kasus yang dibuatnya. Probe live menyimpan
satu kasus sintetis tanpa konfirmasi review.

## Cara melanjutkan

Layanan sesi ini: frontend `http://127.0.0.1:5175`, API `127.0.0.1:8787`,
database `127.0.0.1:55432`. Jika sesi terminal sudah berhenti, ikuti README.
Konfigurasi Apertus tetap pada `../.apertus-api.env`; jangan mencetak/commit key.

Prioritas berikutnya: perluas parser hanya berdasarkan format TXT nyata yang
disepakati, uji kegunaan dengan tester yang sesuai, lalu ukur waktu review,
kesalahan semantik, biaya dan ketersediaan provider. Rancang lifecycle akun,
retensi, dan penghapusan sebelum penggunaan data pribadi/publik. Siapkan
rekaman demo dan materi submission setelah validasi manfaat; cocokkan
aturan/lisensi event sebelum publikasi. OCR dan deploy Apertus mandiri belum
dilakukan. Belum ada izin deployment publik, submission, atau pesan keluar.

## Verifikasi ulang terhadap handoff — sesi lanjutan 9 Oktober 2026

Implementasi sudah tersedia ketika sesi dimulai. Handoff, rencana MVP, dan
hasil validasi riset dibaca untuk mencocokkan batas data dan alur aplikasi.
Prompt aplikasi cocok dengan system prompt, instruksi, dan batas token fixture
v2 yang dibekukan. PostgreSQL khusus project tetap aktif; API dan frontend
dijalankan kembali.

- Persetujuan pengiriman AI sekarang habis pada setiap percobaan, termasuk
  ketika provider gagal. Pengguna harus mencentang ulang sebelum mencoba lagi.
- Batas input AI di UI sekarang memakai ukuran byte UTF-8, sama dengan Go,
  sehingga teks multibyte tidak melewati gerbang 12 KB secara keliru.
- `sh scripts/test-backend.sh`: lulus dengan PostgreSQL dan race detector.
- `pnpm build`: lulus setelah perubahan.
- `pnpm test:ui`: **4/4 lulus**, meliputi alur manual sampai ZIP dan pembatalan
  review setelah edit, persetujuan setelah respons sukses/gagal, serta teks
  multibyte yang terlalu besar untuk AI tetapi tetap bisa disusun manual.
- Health API: `status=ok`, `database=postgresql`. Screenshot desktop/mobile
  diperbarui oleh smoke test; screenshot review mobile diperiksa secara visual.

Respons AI pada uji UI baru disimulasikan di browser. Tidak ada panggilan
provider live tambahan; bukti live sesi awal tetap terpisah. Uji browser
menambahkan empat kasus sintetis ke database aplikasi, sedangkan uji Go
memakai database uji tersendiri.

## Paket testing dua bahasa dan keterbacaan UI

- `docs/testing/SCENARIOS.md`: 13 skenario manual dengan langkah dan hasil
  yang diharapkan. Delapan input sintetis tersedia sebagai TXT dan JSON;
  `RESULTS.md` disediakan untuk mencatat percobaan tester.
- English menjadi bahasa default. Selector English/Indonesia di header
  menyimpan preferensi di browser. Label, petunjuk, validasi, status, dan
  label ekspor mengikuti bahasa. Sumber/draft/isian pengguna tetap asli.
- DM Sans untuk isi dan Manrope untuk heading tetap disajikan lokal. Teks
  utama/input diperbesar menjadi 16 px, warna teks dibuat lebih gelap, tombol
  lebih besar, dan source/draft dipisahkan dengan latar yang jelas.
- Progress review menunjukkan jumlah sumber diperiksa; tanda selesai pada
  navigasi berasal dari keadaan kasus, bukan posisi halaman yang dikunjungi.
- Go dengan PostgreSQL/race detector lulus. Build frontend lulus. Browser
  **7/7 lulus**, termasuk dua bahasa dan viewport 320/390/768/1440 px.
- Screenshot English input desktop, review desktop, dan review mobile
  tersedia pada `screenshots/05…08`. UI desktop/mobile diperiksa visual.

Paket skenario baru belum merupakan hasil uji kegunaan manusia. Fixture
divalidasi terhadap parser/total dengan Go; isi AI memerlukan percobaan dan
review tersendiri. Tidak ada panggilan provider live tambahan.

## Penguatan fitur — 10 Oktober 2026

- Draft Apertus baru mengikuti English/Indonesia yang dipilih; kutipan sumber
  tetap asli. Bahasa saat dibuat, waktu, versi prompt, token jika tersedia, dan
  latensi dicatat pada draft yang diterima. Bahasa UI dapat diganti tanpa
  menerjemahkan draft lama. Contoh sintetis bawaan mengikuti bahasa saat dimuat.
- Preview memperlihatkan ID/teks sumber sebelum persetujuan pengiriman AI.
  Consent direset setelah perubahan bahasa/revisi dan pada setiap percobaan.
- Guard bilingual menandai perubahan arah dana, status transfer, permintaan/
  janji menjadi fakta, dan nominal rupiah di luar kutipan. Pesan review UI kini
  menjelaskan jenis masalah. Ini heuristik yang tetap memerlukan review manusia.
- Lampiran dapat dihubungkan secara manual ke pesan dan transfer. Token kasus,
  validasi ID, revisi, dan pembatalan review berlaku pada endpoint baru. Perubahan
  chat/parser membersihkan hubungan pesan; transfer yang dihapus dibersihkan
  dari hubungan. Isi/keaslian lampiran tidak diperiksa dan tidak ada OCR.
- Ekspor PDF bilingual tersedia setelah review dan disertakan di ZIP. Font
  disematkan lokal, sumber panjang dipaginasi, karakter di luar cakupan font
  ditampilkan dengan kode Unicode. TXT/JSON menjaga sumber lengkap verbatim.
- Lisensi lokal ditambahkan: Apache 2.0 untuk kode asli, CC BY 4.0 untuk materi
  nonkode asli, CDLA Permissive 2.0 untuk input sintetis asli. NOTICE menjelaskan
  cakupan dan pengecualian pihak ketiga. Repo belum dipublikasikan/submitted.

Verifikasi: **9/9 test UI lulus**, build TypeScript/Vite lulus, dan Go dengan
PostgreSQL/race detector lulus. Tes baru mencakup preview/consent bilingual,
14 contoh guard, identitas model dan diagnostik yang dibatasi, isolasi token
hubungan bukti, konflik revisi, pembersihan hubungan usang, PDF dua bahasa, serta
sumber sangat panjang/Unicode. PDF sintetis English dan Indonesia masing-masing
tiga halaman sudah dirender dan diperiksa secara visual. Uji browser ini manual
atau memakai AI mock; bukan bukti ketersediaan provider live.

Percobaan rekaman live English menggunakan data fiktif dan satu panggilan AI
per rekaman. Dua percobaan awal dihentikan saat gagal dan dipertahankan sebagai
`LaporSiap_Failed_1` dan `LaporSiap_Failed_2` (MP4/WebM/JSON) dalam
`output/simulation/2026-10-10/apertus/english/`. Percobaan kedua ditolak karena
metadata model tidak sama dengan Apertus yang diminta; tidak ada draft fallback.
Pemeriksaan katalog API 10 Oktober mengembalikan HTTP 200 dan masih mencantumkan
model tersebut; bukti aman di `evidence/provider-catalog-20261010.json`.

Percobaan ketiga berhenti pada validasi schema/kutipan (`LaporSiap_Failed_3`).
Contoh JSON prompt yang memakai backslash berlebih diperbaiki; instruksi juga
menegaskan preservasi whitespace dan larangan pagar Markdown. Guard tetap
meminta kutipan persis. Diagnostik validasi kini hanya menyebut aturan yang
gagal dan ID model, tanpa isi respons/sumber atau credential.

Percobaan keempat menerima respons dengan metadata `swiss-ai/apertus-v1.5-70b`,
lima kutipan persis, ringkasan English, dan usage 545 input / 366 output token;
latensi permintaan 9.060 ms. Rekaman berhenti pada assertion skrip yang belum
menerima paraphrase "demands" untuk permintaan deposit M4, lalu disimpan sebagai
`LaporSiap_Failed_4`. Assertion diperbaiki setelah ringkasan dan sumber diperiksa.
Draft asli keempat dipertahankan terpisah, dan bukti integrasi/assistant review
tersimpan di `evidence/apertus-english-20261010.json`. Estimasi biaya inferensi
satu respons ini USD 0,00151562 dari usage dan katalog tarif, bukan tagihan atau
biaya operasi total; tidak mencakup panggilan gagal, infrastruktur, atau review.

Percobaan kelima menyelesaikan **8/8 langkah simulasi live English**: input,
lampiran/hubungan bukti, satu respons Apertus asli, pemeriksaan lima sumber,
konfirmasi review, unduhan PDF/TXT/ZIP, dan riwayat lima revisi. Rekaman final
`output/simulation/2026-10-10/apertus/english/LaporSiap_Success_1.mp4` memakai
data fiktif, tanpa mock AI, fallback, atau penggantian ringkasan dengan kalimat
hardcoded. Respons ini memiliki metadata model yang tepat, lima kutipan utuh,
usage 545 input / 366 output token, dan latensi permintaan 3.535 ms. Metadata
adalah keterangan API, bukan audit bobot fisik provider.

ZIP aktual divalidasi: transfer berhasil Rp180.000, pending Rp420.000 tidak
masuk total, tanggal transfer tetap kosong, lampiran dan PDF tersedia. PDF
live tiga halaman dirender dan diperiksa visual. Empat rekaman gagal lengkap
beserta laporan JSON tetap dipertahankan. Keberhasilan satu contoh ini bukan
pengukuran akurasi umum atau manfaat oleh tester independen.

Belum dibuktikan: manfaat/waktu review oleh tester independen, akurasi semantik
di luar fixture, biaya operasional dan skalabilitas, serta deployment Apertus
mandiri. Bukti live 9 Oktober tetap historis dan tidak membuktikan availability
hari ini. Fitur MVP yang lulus tes tidak berarti seluruh persyaratan submission
sudah dipenuhi; Getting Started Guide resmi masih belum dapat diverifikasi.

## Label simulasi dan skenario natural — 10 Oktober 2026

Kata "demo" di label UI serta penanda ekspor TXT/PDF diganti menjadi
"Data simulasi" / "Simulated data". Contoh baru memakai tokoh fiktif Nadia dan
Raka, nama bank fiktif Bank Aruna/Bank Lestari, rekening tersamarkan `****7392`
dan `****4821`, serta sapaan yang lebih natural. Status `synthetic: true` tetap
disimpan, dan UI/lampiran/ekspor tetap menjelaskan bahwa data ini rekaan.
Data kasus, sumber, dan rekaman historis tidak ditulis ulang.

Verifikasi setelah perubahan: Go dengan PostgreSQL/race detector lulus,
build TypeScript/Vite lulus, dan **9/9 test browser lulus** (31,5 detik).
Assertion tambahan memeriksa label tanpa kata "demo", penanda simulasi dalam
dua bahasa, rekening tersamarkan, serta status sintetis pada kasus tersimpan.
PDF sampel English/Indonesia, masing-masing tiga halaman, diekstrak dan
dirender untuk pemeriksaan visual; penanda baru serta identitas fiktif cocok.

Dua percobaan rekaman terhenti sebelum aplikasi terbuka ketika frontend tidak
melayani port 5175. Rekaman dipertahankan sebagai `LaporSiap_Failed_1` dan
`LaporSiap_Failed_2` di subfolder `apertus/english/natural-case/`, tanpa panggilan
AI. Pemeriksaan langsung menunjukkan proses frontend lama sudah berhenti;
frontend dinyalakan kembali dan proxy `/api/health` mengembalikan PostgreSQL
sehat sebelum uji browser dan rekaman berikutnya. Tidak ada retry dalam video.

Rekaman ulang `natural-case/LaporSiap_Success_1.mp4` menyelesaikan **8/8 langkah**
dengan satu panggilan Apertus asli: input fiktif → lampiran/hubungan bukti → draft
English → pemeriksaan sumber dan transfer → review → PDF/TXT/ZIP → riwayat.
Metadata model cocok `swiss-ai/apertus-v1.5-70b`; lima kutipan utuh cocok persis,
usage 546 input / 357 output token, dan latensi permintaan 8.984 ms. Ringkasan
dibaca terhadap sumber oleh assistant, tanpa mengganti hasil AI dengan kalimat
hardcoded. Ini bukan evaluasi manusia independen atau kasus korban nyata.

Paket ZIP diperiksa terhadap revisi review 5: sumber asli, hubungan M3/transfer,
rekening tersamarkan, transfer berhasil Rp180.000, pending Rp420.000, dan tanggal
transfer kosong tetap benar. Isi PDF dalam ZIP sama dengan unduhan PDF terpisah;
bytes bisa berbeda karena metadata pembuatan PDF. PDF live tiga halaman dirender
dan diperiksa visual. Bukti sesi di `evidence/natural-case-20261010.json`.

## Uji unggahan mandiri dan perbandingan manual/Apertus — 10 Oktober 2026

Atas permintaan pengguna, dibuat sepuluh file unggahan fiktif English:
tiga TXT, tiga PNG, tiga JPEG, serta satu PDF referensi bersama. Varian S01
(berhasil/pending), S02 (seluruh transfer gagal), dan S03 (nominal berhasil
belum diketahui) diuji melalui manual dan Apertus dengan bahan yang sama.
**6/6 percobaan lulus**; enam impor TXT dan 18 unggahan lampiran dilakukan
melalui UI nyata, termasuk hubungan sumber/transfer serta unduhan file asli.
Tidak ada perubahan kode aplikasi dalam tahap ini.

Tiga respons live memiliki metadata model `swiss-ai/apertus-v1.5-70b`, bahasa
English, kutipan persis, serta latensi provider 5.775 / 3.380 / 3.664 ms.
Seluruh 11 ringkasan AI dibaca terhadap sumber oleh assistant sebelum review
dikonfirmasi; tidak ada koreksi semantik pada contoh ini. Input transfer tetap
utuh: total Rp180.000, Unknown (bukan nol), dan total sebagian Rp225.000
dengan satu transfer berhasil tanpa nominal. Tanggal transfer tetap kosong.

ZIP/TXT aktual dan bytes lampiran diperiksa; seluruh hubungan bukti/revisi
review tetap tersimpan. Semua 18 halaman PDF ekspor dan PDF input satu halaman
dirender untuk pemeriksaan visual. Enam MP4 H.264 1920×1080 beserta input,
draft asli, keputusan review, hash, dan metrik ada di
`output/testing/2026-10-10/`. Hasil rinci: `testing/RESULTS.md`,
`comparison.json`, dan `comparison.csv` di folder output tersebut.

Ini pengujian developer otomatis, bukan studi manusia independen. Durasi
video/permintaan manual tidak digunakan sebagai penghematan waktu pengguna.
Tiga keberhasilan tidak membuktikan akurasi umum, skalabilitas, atau deployment
model mandiri. Rekaman dan bukti batch terdahulu tetap dipertahankan.

## Uji ketahanan dan validasi lampiran — 10 Oktober 2026

Uji awal membuktikan PNG terpotong 32 byte diterima HTTP 200 karena backend
hanya memeriksa signature MIME. Rekaman langsung dihentikan dan disimpan
sebagai `resilience/LaporSiap_Failed_1.mp4`. Validasi sebelum penyimpanan kini
mendecode PNG/JPEG dengan batas 24 juta piksel, serta memeriksa header,
EOF dan rujukan xref dasar PDF. Pemeriksaan PDF belum validasi penuh objek/
stream, OCR, pemindaian malware atau verifikasi bukti. UI menjelaskan batas
gambar serta error dalam dua bahasa, menggulir ke pesan, dan memberi fokus.

Enam skenario terakhir lulus: file tidak didukung/kosong/terpotong, batas
5 MiB/10 lampiran, tanggal TXT tidak valid, save dengan request terputus,
provider gagal/model berbeda/schema invalid, dan satu source injection live.
Penolakan tidak mengubah sumber/draft/transfer/lampiran/review di database.
Network/provider UI adalah simulasi browser terkendali; Go menguji kontrak
HTTP provider dan timeout tanpa retry. Data yang belum disimpan pada uji
network tetap di tab aktif dan mengunci ekspor, tanpa janji persistensi reload.

Satu Apertus live (`swiss-ai/apertus-v1.5-70b`) menghasilkan tiga kutipan persis,
latensi 4.310 ms, usage 462 input / 205 output token. M2 menjelaskan instruksi
manipulatif sebagai data sumber, bukan menyatakan refund benar-benar terjadi;
M1/M3 tetap tidak menerima refund. Transfer Rp100.000 tetap isian pengguna.
Assistant memeriksa makna sebelum review; ini bukan studi independen atau
bukti umum ketahanan terhadap injection.

Go/PostgreSQL/race lulus (4,923 detik), build frontend lulus, dan suite UI final
10/10 lulus (laporan 1,1 menit). Fixture PDF lama diganti file valid; test review
menunggu draft/sumber selesai dimuat. Konfigurasi Playwright merekam tes dan
menggunakan cache FFmpeg recorder bila tersedia. Rekaman kegagalan sinkronisasi
test dipertahankan sebagai `LaporSiap_Failed_2`, beserta konteks kegagalan setup
FFmpeg dan fixture awal. Ada tujuh MP4 sukses, dua MP4 gagal, dan sepuluh WebM
regresi final dalam `output/testing/2026-10-10/resilience/`. Ringkasan, hash,
input negatif, draft AI asli dan keputusan sumber disimpan di folder itu;
hasil rinci di `testing/RESULTS.md`. Tidak ada publikasi/deployment/submission
atau penulisan ulang bukti batch sebelumnya.

## Persiapan pilot pengguna — 10 Oktober 2026

Kit untuk dua peserta independen dan satu peserta opsional tersedia di
`testing/pilot/`: brief Indonesia/English, panduan fasilitator, rencana urutan
tugas, dan template JSON kosong. P01 mencoba S01 manual lalu S03 Apertus;
P02 mencoba S01 Apertus lalu S03 manual. Peserta ketiga dapat memeriksa S02
secara eksploratif. Setiap peserta utama mendapat kasus berbeda pada kedua
mode; perbandingan tetap pilot kecil dengan kemungkinan pengaruh kemampuan
peserta. Pemilik boleh latihan internal pada S02, dicatat terpisah.

Bahan peserta menggunakan TXT/PNG/JPG fiktif yang telah ada. PDF referensi
dan JSON kunci jawaban dipisahkan karena mengungkap hasil yang diharapkan.
Waktu tugas, tunggu AI, koreksi, bantuan, kualitas akhir, pemahaman total,
serta rekaman gagal/sukses memiliki aturan pencatatan. Tidak ada sesi manusia
atau panggilan provider tambahan pada persiapan ini. Template tetap kosong;
belum ada klaim penghematan waktu, testimoni atau manfaat terukur.

## Alat perekaman sesi manusia — 10 Oktober 2026

Atas permintaan bantuan menjalankan tahap berikutnya, tersedia
`scripts/demo/record-user-session.mjs`. Defaultnya latihan internal S02 manual;
pasangan tugas P01/P02/P03 mengikuti rencana pilot. Chrome membuka aplikasi
dengan kontrol Start, Pause/Resume, finish, failure, dan stop. Input fiktif
disalin ke folder sesi unik; unduhan dan result.json disimpan bersama WebM
serta MP4. Tidak ada pengisian data, consent, review, retry, atau fallback
otomatis. Gangguan API/server/JavaScript saat tugas berjalan menghentikan
rekaman. Kualitas, koreksi dan pemahaman tetap perlu penilaian fasilitator.

Health lokal menunjukkan aplikasi/PostgreSQL aktif. Syntax/dependensi lulus;
check teknis browser/binding/video/transcode MP4 lulus, tanpa API aplikasi
atau provider. Sesi manusia dan alur kontrol lengkap belum dilakukan.

## Penguatan tujuh poin — 10 Oktober 2026

Permintaan pemilik dilaksanakan mulai dari lifecycle data, batas biaya/kapasitas,
lalu privasi sebelum evaluasi/deployment.

| Prioritas | Hasil konkret | Status bukti |
| --- | --- | --- |
| 1. Lifecycle | Hapus kasus+history+files, token rotation/revoke, access 7 hari, retention 30 hari, cleanup startup/hourly | Implementasi dan tes Go/PostgreSQL/UI lulus |
| 2. Biaya/kapasitas | Kuota persisten 3/kasus dan 20/global per rolling 24 jam, 2 aktif; storage/revision/case/request limits | Tes concurrency/restart/deletion bypass lulus |
| 3. Privasi | Masking default, custom terms, preview aktual, hash consent, original quote lokal, label provider sesuai endpoint | Boundary HTTP dan UI dua bahasa lulus; deteksi terbatas |
| 4. Pilot manusia | Recorder, kit peserta, counterbalanced tasks, analyzer pengamatan | Kit siap; peserta independen aktual **0** |
| 5. Evaluasi AI | 12 kasus bilingual baru, 26 ringkasan dibaca, 4 isu dicatat; prompt/guard diperkuat dan 2 input follow-up baru | 10/12 diterima awal; 2 follow-up timeout, kualitas live prompt baru belum terbukti |
| 6. Deployment/biaya | Compiled app+PostgreSQL terpisah; GPU Compose/guide; tarif/usage aktual; bounded read check | Deployment aplikasi lokal diuji; **model GPU belum dijalankan** |
| 7. Paket review | README/technical report/fixture/script, source/evidence ZIP dan manifest SHA256 | Lokal; tidak dipublikasikan/disubmit; panduan akhir lomba belum terverifikasi |

Panel Data dan akses Anda / Your data and access tersedia di semua tahap.
Token baru mencabut token lama tanpa mengubah content revision/review.
Deletion mencakup semua history/file; backup dan unduhan terpisah tidak ikut
terhapus. Retention sejak creation, bukan setiap edit. Expired access sendiri
tidak mempercepat deletion; token hilang tidak dipulihkan. Migrasi backfill
expiry mempertahankan evidence lama. Kapasitas menghitung payload logis
JSON/history/file, bukan ukuran disk/WAL. Gagal/ditolak tetap terhitung;
global AI meter tidak reset oleh penghapusan kasus.

Hash preview berubah bila revision/bahasa/terms/mode berubah, sehingga consent
harus diberikan ulang. Istilah personal direset saat berganti kasus. Quotes
asli dipulihkan hanya setelah respons terhadap masked sources lolos contract;
original source/export tetap dapat mengandung data pribadi.

Evaluasi di `output/evaluation/2026-10-10T14-36-32.287511000Z/` menyimpan
12 attempt tanpa retry/fallback. Sepuluh diterima; source-instruction English
ditolak karena metadata model berbeda, Indonesian gagal pada timeout/koneksi.
Assistant review menemukan empat isu pada tiga draft Indonesia: tindakan kirim
hilang/nominal pembayaran dikaitkan refund, pelaku janji perekrut berubah,
nominal transfer gagal hilang, dan konteks klaim sebelumnya hilang. Saran koreksi
dicatat terpisah; draft asli tidak diedit. Ini development review, bukan holdout
atau penilaian independen. Contract-valid tidak berarti makna akurat.

Prompt v4/guard menandai tindakan/nominal/peran yang hilang. Dua follow-up baru di
`output/evaluation/2026-10-10T15-29-06.825664000Z/` gagal pada timeout/koneksi
40 detik; peningkatan live belum dapat dinilai. Read-only verification memeriksa
14 kasus: raw chat/messages/transfers tidak berubah; empat gagal mempertahankan
revision 1 tanpa draft/review. Sepuluh accepted melaporkan 4578 input/1809 output
token, estimasi USD 0.00903624 dari katalog authenticated. Usage/biaya empat gagal
tidak diketahui. Ini bukan invoice atau biaya operasional keseluruhan.

Dependency check menemukan 28 advisory berbeda pada Go 1.26.3/x-text 0.29.0,
termasuk fungsi yang dipakai. Go minimum/Docker builder kini 1.26.9,
x/text 0.41.0 dan x/sync 0.22.0. Govulncheck v1.8.0 setelah perbaikan mencatat
**0 finding** (database 8 Oktober); audit produksi frontend juga 0 advisory.
JSON-mode exit 0 berarti scan selesai, sehingga finding dibaca eksplisit.
Ini pemeriksaan advisory, bukan audit independen atau image scan.

Backend/PostgreSQL/race final lulus (6,383 detik), build frontend lulus,
**13/13 browser** lulus (46,8 detik), **2/2 compiled deployment** lulus
(6,7 detik), dengan MP4/WebM dan ekspor dua bahasa. **20/20 authenticated reads**
lulus pada concurrency empat tanpa AI; kasus khususnya dihapus sesudah uji.
Ini cek lokal kecil, bukan throughput AI/kapasitas produksi. Container app
UID65532, read-only, cap-drop ALL, no-new-privileges, CSP, port127.0.0.1:8877.
Digest Docker base/DB dipin; image final dicatat di final-checks.json.

Suite intermediate dimulai saat restart backend oleh assistant sehingga create
mendapatkan proxy 502. Rekaman gagal dipertahankan sebagai
`output/hardening/2026-10-10/browser-intermediate/LaporSiap_Failed_1.mp4`.
Dua tes deployment skipped sempat salah diberi nama Failed oleh teardown;
catatan menjelaskan tidak ada alur deployment yang dijalankan. Hook diperbaiki,
tes deployment dipisahkan dari suite default. Hasil final tidak menimpa
intermediate/first pass. Kegagalan setup tidak dihitung sebagai keberhasilan.

Laporan review: `evaluation/REVIEW-REPORT.txt`. Hasil final, biaya, state check,
pilot kosong, security scan dan video di `output/hardening/`. Generator
`scripts/package-review.py` membuat ZIP source/evidence dan manifest hash tanpa
.env/credential/dependency cache/model weights. NOTICE mempertahankan hak output
model/dependensi. Belum ada public deployment, repository publication, submission,
human benefit study atau GPU proof. Guide lomba tertaut masih 404 saat diperiksa;
kesesuaian semua deliverable belum dapat dinyatakan lengkap.


## Pemeriksaan lanjutan — 11 Oktober 2026

Hash prompt embedded dikunci dan diverifikasi terhadap backend sebelum evaluasi.
Metadata draft menyimpan hash; evaluator menyimpan diagnostik aman dan membaca
kembali state meskipun provider gagal. Tes backend/race lulus; dua alur manual
hasil build dalam Inggris dan Indonesia lulus (6,8 detik). Container privat
terbaru sehat, masih di loopback 8877.

Dua recheck pada fixture development yang sama, satu panggilan per bahasa,
tersimpan di `output/evaluation/2026-10-10T17-28-37.268Z/` (UTC; 11 Oktober WIB).
- Inggris: HTTP 502, identitas Apertus, kutipan berbeda; ditolak.
- Indonesia: HTTP 502, identitas `aisingapore/Qwen-SEA-LION-v4-32B-IT`; ditolak.
- Dua video `LaporSiap_Failed_1.mp4`, masing-masing dalam folder bahasa.
- Sumber/transfer tetap utuh, revisi tetap 1, draft/review tetap kosong.
- Tidak ada retry, fallback aplikasi, output v4 yang diterima, atau skor kualitas.

Recorder awal gagal menemukan FFmpeg karena lokasi cache keliru; dua kasus lokal
tetap tersimpan dan nol provider call terjadi. Setelah lokasi cache diperbaiki,
recorder check kosong lulus sebelum panggilan live. Riwayat teknis ini berada di
`output/evaluation/2026-10-10T17-27-10.741Z/` dan tidak dihitung sebagai dua
kegagalan provider tambahan. Audit baca-saja seluruh 18 kasus memverifikasi
sumber/transfer dan state kegagalan.

`output/hardening/2026-10-11/live-recheck.json` merangkum bukti dan batasnya.
Prioritas selanjutnya: endpoint/operator yang menjamin identitas Apertus, uji
peserta independen (masih nol), inference GPU mandiri (belum dijalankan), dan
panduan resmi track. Link Getting Started Guide dari Devpost kembali 404;
halaman Resources saat diperiksa menyebut pendaftaran dukungan CSCS/Phoeniqs
sudah ditutup. Tidak ada klaim seluruh persyaratan lomba telah terverifikasi.


## Clean Architecture dan Redux — 11 Oktober 2026

Backend direfaktor ke domain, application/use case, adapter infrastructure dan
transport HTTP. Domain/use case tidak mengimpor HTTP, pgx, fpdf atau adapter.
Aturan create, kenaikan revisi, invalidasi/konfirmasi review, privacy/consent,
attachment/link, lifecycle dan export ditangani use case melalui port. PostgreSQL
mempertahankan transaksi, CAS expected revision, expiry, capacity dan lease quota.
Prompt, migration, font dan fixture dipindahkan bersama adapter/package pemiliknya.
Command API, maintenance/evaluation dan recorder memakai lokasi baru.

Frontend memiliki domain, application, infrastructure dan presentation, dengan
composition root untuk dependency injection. Redux Toolkit 2.13.0 / React Redux
9.3.0 mengelola workspace dan locale, termasuk form sementara. Session vault
menyimpan token sebelum respons kasus dikembalikan ke Redux; token/file/fungsi
callback tidak masuk state/action. Reducer dan typed selector/dispatch mengganti
seluruh useState aplikasi. Persistensi Redux/DevTools/logger tidak diaktifkan.

Pemeriksaan arah dependensi Go/TypeScript ditambahkan untuk menjaga layer. Unit
use case backend tidak memerlukan HTTP/DB/provider; tiga unit frontend memeriksa
pemisahan token, delete gagal, serta reset privasi/review. Regresi lama dipindahkan
ke internal/integration dengan helper nama kompatibel khusus test, tanpa facade
produksi. Unit Node dikecualikan dari discovery Playwright.

Validasi final: backend/race lulus (integration 6,207s); frontend unit 3/3;
architecture check 23 file TS; build 42 module; UI 13/13 dalam 22,9s; compiled
container Inggris/Indonesia 2/2 dalam 5,2s; audit dependensi produksi tanpa
advisori yang diketahui. Initial UI 13/13 dalam 58,1s juga disimpan. Kedua alur
container merekam video Success_1 dan menjalankan upload/review/export/rotation/
deletion. Zero provider generation call pada refaktor ini. Tidak ada perubahan
schema database atau klaim perbaikan provider live dari perubahan arsitektur.

Bukti: output/architecture/verification.json; struktur dan kontrak dijelaskan di
README.md dan TECHNICAL.md. Image privat final:
sha256:94d61ed081aeaa77f6cf85cc6c1b8a86751622f0b3c3421d6e18763de904d058.

## Integrasi Apertus v5 — 11 Oktober 2026

Kontrak AI diubah secara eksplisit: hanya source_id/description. Domain memasang
kutipan lokal setelah count/order/ID/schema/35 kata divalidasi. Evidence buatan
provider ditolak, bukan dinormalisasi. Validator draft lengkap, batas masking,
pinning model, consent/hash/revisi, quota dan manual fallback tetap dipertahankan.
Versi prompt bilingual-v5-source-bound, SHA256
b1a96426b541d40f84a21435975d81e4315fbede2a86efe2b3bec941421d66bf.

Provider ValidateConfiguration dipakai saat generasi dan startup ber-key;
CheckAvailability melakukan satu GET /models tanpa completion/retry/redirect.
Command providercheck membaca policy runtime dan kuota/capacity SQL; tidak
menampilkan env/key/token/URL privat, menulis atau memigrasi database. Request
model tetap eksplisit sesuai dokumentasi Public AI; penyebab mismatch model
sebelumnya tidak diketahui dan tidak diklaim telah diperbaiki oleh konfigurasi.

Preflight 16/20 attempt dan 175/200 kasus. Dua panggilan lewat UI diterima EN/ID
HTTP200, model Apertus, empat kutipan persis (termasuk newline akhir), sumber
dan transfer utuh; provider 2051/2847ms, usage 1325 input/196 output. Postflight
18/20 attempt, active lease 0, 179 kasus setelah regresi browser. Audit read-only
20 kasus evaluasi lulus tanpa generasi. Tidak ada retry/fallback. Rekaman:
output/evaluation/2026-10-10T18-47-16.862Z; keduanya LaporSiap_Success_1.mp4.

Review assistant terpisah: satu omission EN M2 (next week), tiga entri lainnya
mempertahankan kriteria yang dinyatakan. Output/revision AI asli tidak diganti.
Guard relative_time_requires_review ditambahkan untuk lima kelompok frasa
EN/ID; bukan verifikasi semantik/tanggal atau koreksi otomatis. Guard dihitung
saat generate/edit; snapshot historis tidak ditulis ulang. Dua sumber development
ini bukan holdout atau bukti provider stabil. Success_1 berarti kontrak diterima,
bukan tidak ada kekurangan makna. Seluruh hasil gagal terdahulu dipertahankan.

Backend/race final lulus (integration 6,101s), frontend build 42 module, unit 3/3,
architecture 23 file TS, browser terarah 2/2 (7,1s; mock AI/no provider calls).
Bukti final: output/evaluation/provider-v5-2026-10-11/verification.json.

Container privat final lulus 2/2 (5,1s), image
sha256:c9f89d70fb03d806facffac1d8029978af9a582b67608d662b295595f6579d8b.

CI lokal ditambahkan: tiga job backend/race/PostgreSQL, frontend architecture/
Redux/build, serta 13 UI regression pada DB service baru tanpa key/API live.
Action dipatok SHA resmi; Go dari go.mod, Node24.15.0/pnpm11.1.3/frozen lock.
Chrome dan helper FFmpeg diinstal di runner; readiness hanya GET lokal45s.
Seluruh 13 browser final lulus 23,4s. Readiness awal timeout karena sandbox
localhost; rerun berizin lulus tanpa perubahan kasus/provider. Workflow belum
dipush/dijalankan di GitHub; validasi lokal bukan bukti runner GitHub berhasil.

Actionlint v1.7.12 lulus; laporan output/ci/verification.json.

## Perbaikan waktu v6 dan persiapan GPU — 11 Oktober 2026

Domain kini menolak draft AI baru yang menghilangkan, mengubah atau menambah
lima kelompok frasa waktu eksplisit sebelum persist. Aturan dipakai pada binding
sumber dan restoration masked, tanpa koreksi otomatis, retry atau fallback.
Prompt mempertahankan waktu bersama janji asal. Versi bilingual-v6-time-preserved,
SHA256 96e58e67a5dafec4ac1c23d788619ed6acfa048a5e887872112e25a3baafb96e.

Backend/race lulus (integration 11,022s), mencakup omission aktual, terjemahan
EN/ID, perubahan waktu, beberapa frasa dan preservasi draft lama saat gagal.
Dua live EN/ID berhasil dan direkam sekali per bahasa: metadata Apertus,
next week/pekan depan tetap pada janji komisi 63000 yang belum dibayar;
pengiriman 264000 dan refund belum diterima tetap. Kutipan, sumber dan transfer
utuh. Review assistant empat entri tidak menemukan issue pada kriteria yang
dinyatakan, bukan review manusia independen atau bukti akurasi umum. Output v5
asli beserta issue dan seluruh kegagalan disimpan. Rekaman:
output/evaluation/2026-10-11T07-59-55.575Z (Success_1 dalam dua bahasa).

Audit read-only 22 kasus evaluasi lulus. Kuota 18->20/20, active lease 0; tidak
ada panggilan tambahan atau peningkatan batas. Postflight exit 1 berarti tidak
ada ruang dua panggilan baru, bukan kegagalan dua respons v6. Preflight awal
sebelum restart menolak hash backend lama tanpa completion; sesudah restart
hash cocok. Container compiled baru lulus 2/2 (10,1s) tanpa provider calls,
image sha256:0c2070b6742df2b9da395c73d77a03d1c3fd375d5134964d9552fa80403707de.
Ringkasan: output/evaluation/timing-v6-2026-10-11/verification.json.

Lock GPU kini berisi digest registry resmi
 d85533f251c8272218ae2e519510836376e5a335b44f2204f193018e37753bc6
serta revisi model 59e744e313e967811aefabde3732609b64201acc dari metadata publik.
Recipe CUDA 13 membutuhkan driver 580+; profil awal 4xH100 80GB, RAM>=256GiB,
ruang kosong>=250GiB, BF16/konteks8192/dua sequence belum diuji nyata. Panduan
GPU-PROOF.txt menyediakan jalur VM Linux dengan Compose dan rute Runpod
inference-only, karena Pod tidak mendukung Compose. Harga primer Lambda
$4,09/GPU/jam memberi contoh $16,36/jam atau $32,72 untuk dua jam compute;
bukan penawaran checkout atau tagihan total. Boot, loading, download, idle,
pajak, storage dan ketersediaan masih perlu dikonfirmasi.

Alat preflight/init-env/serve/collect GPU memakai default dry-run, pin resmi,
key env privat di luar argv, observasi field terpilih tanpa Env/raw inspect,
serta batas driver/RAM/VRAM/disk. Collector tidak memberi klaim proof lengkap.
Empat tes offline, compile Python, merge Compose GPU dan actionlint lulus.
Origin8878 sesuai tunnel; project/database terpisah. Probe Mac mencatat nol GPU.
Tidak ada GPU disewa, image/weights GPU diunduh, server remote dibuat, terms
diterima otomatis, atau akun cloud baru dipakai. Metadata dan persiapan bukan
bukti deployment. Laporan: output/gpu-preparation/2026-10-11/verification.json.

## Kepatuhan template Hack Apertus — 13 Oktober 2026

Syarat track 2B dibaca dari template resmi `HackApertus/project-template`
(`track_2b/README.md`): `make run` di Docker dari checkout bersih, env
`LLM_NAME`/`LLM_BASE_URL`/`LLM_API_KEY`, report PDF maks. 6 halaman, demo
maks. 2 menit, repo publik dengan folder `track_2b/`.

Ditambahkan `Makefile` dan `compose.run.yaml` (PostgreSQL + app, tanpa `.env`,
tanpa `pull_policy: never`). Backend membaca `LLM_*` lalu fallback `APERTUS_*`;
model diterima bila ID keluarga Apertus (case-insensitive), `http://` hanya
untuk loopback/IP privat/nama service satu label. Konfigurasi AI tidak valid
kini menonaktifkan AI dengan log, bukan menghentikan aplikasi. Suite backend
race + PostgreSQL lulus; `make run` stack diuji: health 200, create/manual
draft/ekspor terkunci sebelum review, label provider lokal, model non-Apertus
ditolak tanpa crash.

Restrukturisasi selesai: project root kini `track_2b/` (kode di `src/`), path
scripts/test/CI/Dockerfile/NOTICE disesuaikan. Media berat `output/` (zip, webm,
mp4, wav, png, jpg) diabaikan Git, bukti JSON/MD/PDF/TXT tetap; ukuran yang akan
di-commit 433 file / 7,2 MB tanpa `.env`. Verifikasi: backend race+PostgreSQL,
FE architecture/unit/build, test GPU offline, 13/13 browser pada database baru
(DB dev lama sudah mencapai MAX_CASES), dan `make run` dari salinan bersih
(build dari nol, health 200, create/manual draft/ekspor terkunci). Teks UI
provider kini netral, bukan "Public AI"/"70B" statis.
