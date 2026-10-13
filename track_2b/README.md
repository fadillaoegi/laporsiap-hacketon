# LaporSiap

MVP lokal untuk menyiapkan bahan laporan dari percakapan tentang pekerjaan
daring/top-up. Pengguna menyimpan sumber, mengisi transfer, menyusun draft
dengan **Apertus 1.5 70B** atau secara manual, memeriksa tiap entri, lalu
mengunduh bahan untuk kanal resmi.

Backend **Go 1.26.9 + pgx**, frontend **React + Redux Toolkit + Vite + TypeScript**, database
**PostgreSQL 18 Docker**. Tidak menggunakan SQLite.

## Run (judges / clean checkout)

```sh
make run          # builds and starts PostgreSQL + LaporSiap in Docker
# open http://localhost:8787  — stop with Ctrl+C, then `make stop`
```

No `.env` is needed. Without an LLM endpoint the full workflow works with the
manual draft path. To enable Apertus drafts, export the Hack Apertus variables
before `make run`:

```sh
export LLM_NAME=swiss-ai/apertus-v1.5-70b   # any Apertus 8B/70B identifier
export LLM_BASE_URL=https://api.publicai.co/v1   # or an on-prem vLLM, e.g. http://10.0.0.5:9000/v1
export LLM_API_KEY=...
```

Plain `http://` is accepted only for loopback, private-network addresses and
single-label container hostnames; public endpoints must use HTTPS. Invalid AI
configuration disables AI (logged) instead of stopping the app.

## Arsitektur

Backend dan frontend memakai Clean Architecture: domain dan use case bergantung
ke kontrak internal; HTTP, database, browser storage, dan provider menjadi adapter
yang disuntikkan pada composition root. Pemeriksaan otomatis menjaga batas impor.

```text
src/laporsiap_BE/internal/
  domain/                  model, parser, validasi, privasi, aturan kronologi
  application/             use case kasus dan interface repository/provider/export
  infrastructure/          PostgreSQL, Apertus, renderer laporan, konfigurasi
  transport/httpapi/       route, DTO HTTP, autentikasi, rate limit, frontend build
  integration/             regresi integrasi dengan adapter nyata
  architecture/            pemeriksaan arah dependensi
src/laporsiap_FE/src/
  domain/                  model dan aturan input
  application/             CaseWorkflow dan interface adapter
  infrastructure/          HTTP, session vault, preferensi bahasa, unduhan
  presentation/            UI, controller, Redux workspace/locale
  app/bootstrap.ts         composition root dan dependency injection
```

Redux Toolkit mengelola kasus, draft, consent/preview, review, navigasi, hubungan
lampiran, konfirmasi penghapusan, serta bahasa. Token akses berada di session vault,
tidak masuk ke state/action Redux. Case text hanya berada di memory state; tidak ada
persistensi Redux atau logger/DevTools. Session token tetap di `sessionStorage`,
sedangkan preferensi bahasa tetap di `localStorage`.

Pemeriksaan pengembangan:

```sh
sh scripts/test-backend.sh
pnpm --dir src/laporsiap_FE check:architecture
pnpm --dir src/laporsiap_FE test:unit
pnpm --dir src/laporsiap_FE build
pnpm --dir src/laporsiap_FE test:ui
```

## Jalankan lokal

Konfigurasi `.env` lokal sudah disiapkan pada sesi implementasi, dengan
credential database acak dan mode `600`. Jangan menyalin credential ke Git.
Untuk checkout baru, gunakan `.env.example` sebagai acuan; isi password acak
dan `DATABASE_URL` yang cocok. Image PostgreSQL harus tersedia karena Compose
memakai `pull_policy: never`.

```sh
cd track_2b
docker compose up -d --wait
sh scripts/dev-backend.sh
```

Di terminal lain:

```sh
cd src/laporsiap_FE
pnpm install --frozen-lockfile
pnpm dev
```

Buka **http://127.0.0.1:5175**. API berjalan di
`http://127.0.0.1:8787/api/health`; PostgreSQL di loopback `55432`.
Backend otomatis menerapkan migrasi yang disematkan pada startup, dalam
transaksi dengan advisory lock. Container dan volume khusus project ini.
PostgreSQL 18 memasang volume di `/var/lib/postgresql`, sesuai
[README image resmi](https://github.com/docker-library/docs/blob/master/postgres/README.md#pgdata).

Script backend membaca `.env` dan, bila tersedia,
`../../.apertus-api.env` (di luar repository). Variabelnya
`APERTUS_API_BASE_URL`, `APERTUS_API_KEY`, `APERTUS_MODEL`. Key hanya dibaca
backend; tidak disalin ke frontend atau repo. Tanpa konfigurasi provider,
jalur manual tetap tersedia. Endpoint mengikuti
[dokumentasi Public AI](https://platform.publicai.co/docs).

## Coba alur aplikasi

Antarmuka default **English**. Pilih **Indonesia** pada menu **Language /
Bahasa** di header untuk mengikuti label Indonesia di bawah. Pilihan bahasa
tersimpan di browser. Label/petunjuk, pesan validasi, dan label ekspor TXT/CSV/PDF
mengikuti pilihan. Draft AI baru diminta dalam bahasa terpilih; mengganti bahasa
UI tidak menerjemahkan draft yang sudah tersimpan. Kutipan sumber serta isian
pengguna tetap asli. Bahasa contoh sintetis mengikuti pilihan saat contoh dimuat.
Contoh memakai tokoh fiktif Nadia/Raka, nama bank fiktif, dan rekening tersamarkan.
Label UI/ekspor memakai "Data simulasi" / "Simulated data", sehingga contoh
tetap dikenali sebagai rekaan meskipun percakapannya dibuat lebih natural.

1. Klik **Coba contoh**. Semua contoh berlabel sintetis; tidak memuat korban nyata.
2. Klik **Simpan & lanjut**, kemudian **Gunakan draft manual**, atau
   periksa preview ID/teks sumber yang akan dikirim, setujui pengiriman ke
   provider yang tertera, dan pilih Apertus. Masking aktif secara default;
   nama pengirim yang terdeteksi, email, nomor telepon Indonesia, serta nomor
   identitas/rekening berlabel diganti placeholder. Tambahkan istilah sendiri
   untuk alamat/nama lain. Deteksi terbatas; baca seluruh preview sebelum consent.
3. Untuk lampiran, kembali ke **Bahan kasus** dan unggah PNG/JPEG/PDF. Pilih
   **Hubungkan bukti** untuk menautkan file ke pesan dan catatan transfer.
   Hubungan ini dipilih pengguna, tanpa verifikasi keaslian atau OCR.
4. Pada **Periksa & koreksi**, cocokkan draft dengan sumber. Simpan koreksi,
   lalu periksa ulang tiap entri, transfer, serta data belum lengkap.
5. Konfirmasi review dan unduh PDF, TXT, atau ZIP. ZIP mencakup PDF, kronologi,
   sumber mentah, transfer CSV, JSON revisi yang ditinjau, dan file lampiran asli.

Ekspor ditolak backend tanpa konfirmasi revisi terbaru. Edit input, draft,
penambahan/penghapusan lampiran, atau perubahan hubungan bukti membatalkan
konfirmasi. Perubahan percakapan/format menghapus draft dan hubungan pesan
yang sudah tidak sesuai sumber; hubungan ke transfer yang dihapus juga dibuang.
Persetujuan AI direset setelah perubahan revisi/bahasa dan setiap percobaan.

PDF dibuat lokal dengan font yang disematkan dan kutipan asli. Karakter yang
tidak didukung font ditampilkan sebagai kode seperti `[U+1F600]`, dengan catatan
di PDF. Sumber lengkap tetap tersedia verbatim di TXT/JSON dalam ZIP.

Paket pengujian manual ada di [docs/testing/SCENARIOS.md](docs/testing/SCENARIOS.md):
13 skenario, delapan TXT sintetis yang bisa diimpor, isian transfer, hasil yang
diharapkan, dan [lembar hasil](docs/testing/RESULTS.md). Mulai dari S01 untuk
alur normal, lalu S02/S03 untuk total gagal dan nominal kosong.

## Format input dan batas

- **Teks biasa:** satu pesan per baris nonkosong. Pengirim dan waktu kosong.
- **TXT berheader:** `[DD/MM/YYYY, HH:MM] Pengirim: pesan`, juga mendukung
  detik `HH:MM:SS`. Baris berikutnya tanpa header menjadi lanjutan pesan.
  Tanggal kalender diperiksa; timezone tidak ditebak. Ini format terbatas,
  bukan dukungan semua variasi ekspor WhatsApp. Format ambigu ditolak;
  pengguna dapat memilih teks biasa.
- Maksimal 100 pesan / 32 KB per kasus dan 6 KB per pesan. AI dibatasi
  20 pesan / 12 KB; kasus lebih panjang dapat disusun manual.
- Transfer maksimal 100 catatan. Nominal rupiah bulat atau kosong, tanggal
  `YYYY-MM-DD` atau kosong, status berhasil/pending/gagal/belum dipastikan.
- Lampiran maksimal 10 file, masing-masing 5 MB. PNG/JPEG diperiksa dengan
  decoding gambar, dengan batas 24 juta piksel. PDF diperiksa header,
  penutup, dan rujukan xref dasarnya; ini belum validasi penuh keterbacaan PDF.
  Lampiran tetap referensi, tanpa OCR, pemindaian malware, atau verifikasi keaslian.

Total hanya menjumlahkan nominal diketahui dari transfer berstatus berhasil.
Tanpa nominal berhasil yang diketahui, total tetap **belum diketahui**, bukan
nol. Total sebagian diberi keterangan. Nilai ini adalah nominal transfer
berhasil yang dilaporkan pengguna, bukan kerugian bersih atau verifikasi bank.
Waktu pesan tidak pernah diubah menjadi tanggal transfer.

## Verifikasi

Database `laporsiap_test` dan `TEST_DATABASE_URL` lokal sudah tersedia.
Untuk lingkungan baru, buat database uji tersendiri dan atur variabel tersebut.

```sh
sh scripts/test-backend.sh
cd src/laporsiap_FE
pnpm build
pnpm test:ui
```

Uji browser membutuhkan API/frontend aktif dan Google Chrome lokal. API/AI
diuji memakai Go; Playwright dipakai untuk interaksi UI. Tanpa
`TEST_DATABASE_URL`, `go test ./...` hanya menjalankan unit test dan melewati
uji PostgreSQL. `scripts/test-backend.sh` mewajibkan variabel database uji.

Playwright merekam uji browser. Bila FFmpeg belum terpasang, jalankan
`pnpm exec playwright install ffmpeg` dari `src/laporsiap_FE`. Jika direktori
`tmp/playwright` sudah digunakan skrip rekaman, konfigurasi memakai cache itu;
pasang dengan `PLAYWRIGHT_BROWSERS_PATH=../../tmp/playwright pnpm exec playwright
install ffmpeg` agar lokasi sama.

Uji negatif dengan rekaman: `node scripts/demo/test-resilience.mjs --case files`
dari root proyek. Pilihan lain: `size`, `dates`, `network`, `provider`,
`injection`. Network/provider adalah simulasi browser; `injection` mengirim
satu permintaan Apertus asli dan menunggu keputusan pemeriksaan sumber.
Hasil dan rekaman berada di `output/testing/2026-10-10/resilience/`.

Satu probe integrasi live opsional, dengan tiga pesan sintetis saja:

```sh
cd src/laporsiap_BE
go run ./cmd/smoke -live
```

Probe mengirim teks ke provider yang dikonfigurasi melalui API lokal. Ia tidak
mengonfirmasi review atau mengekspor otomatis. Jangan mengulangnya sebagai
batch akurasi. Detail teknis ada di [docs/TECHNICAL.md](docs/TECHNICAL.md).

Rekam simulasi English dengan panggilan Apertus asli:

```sh
node scripts/demo/record-simulation.mjs --apertus --english --natural
```

Satu permintaan AI dilakukan per rekaman. Jika gagal, rekaman langsung dihentikan
dan disimpan sebagai `LaporSiap_Failed_N.mp4` beserta WebM dan laporan JSON di
`output/simulation/YYYY-MM-DD/apertus/english/`. Tidak ada retry atau fallback
manual dalam rekaman itu. Rekaman manual memakai flag `--english` saja.
Flag `--natural` menyimpan hasil pada subfolder `natural-case/` agar rekaman
sebelumnya tetap tersedia.
Metadata token/latensi disimpan hanya untuk respons yang diterima; belum
merupakan pengukuran biaya atau manfaat pengguna.

## Akses dan penyimpanan

Setiap kasus mendapat token acak 192-bit. Database menyimpan hash SHA-256;
token dibawa lewat header Bearer, bukan URL. Semua endpoint kasus, riwayat,
lampiran, dan ekspor memeriksa pasangan kasus/token. Frontend menyimpan token
di **sessionStorage**, sehingga akses sesi hilang setelah sesi tab ditutup,
sementara data tetap ada di PostgreSQL. Akses berlaku tujuh hari; panel **Your
data and access / Data dan akses Anda** memperbarui token dan mencabut token
lama. Kasus disimpan maksimum 30 hari sejak dibuat dan dibersihkan pada startup
serta tiap jam; akses langsung ditolak sesudah kedaluwarsa. Penghapusan melalui
panel dengan mengetik DELETE menghapus kasus, semua riwayat dan lampiran dari
database. Unduhan dan backup terpisah tidak ikut dihapus. Belum ada akun atau
pemulihan token yang hilang/kedaluwarsa.

Batas default: 200 kasus, 200 revisi per kasus, dan 256 MiB payload logis kasus,
riwayat serta file (bukan ukuran disk PostgreSQL/WAL/backup). AI dibatasi tiga
percobaan per kasus dan 20 secara global dalam 24 jam bergulir, dengan dua
permintaan aktif. Hitungan AI persisten di PostgreSQL; gagal/ditolak tetap
terhitung. Penghapusan kasus tidak mereset batas global. Request API dibatasi
240 per IP per menit dalam memori proses. Semua batas dikonfigurasi pada startup
melalui `.env.example`; ini perlindungan prototipe, bukan kapasitas layanan publik.

Default hanya loopback; backend menolak bind publik, host/origin lain, dan
tidak membuka CORS. Ini prototipe lokal untuk demo sintetis. Penggunaan publik
dengan data pribadi memerlukan autentikasi akun, lifecycle akses, TLS, retensi,
backup, dan evaluasi keamanan lebih lanjut.

Tautan [IASC](https://iasc.ojk.go.id/) tersedia sejak halaman awal dan ketika AI
gagal. LaporSiap bukan kanal pengiriman laporan, bukan tiket IASC, dan tidak
berafiliasi dengan OJK/IASC. Tidak menjanjikan pemblokiran, pengembalian dana,
atau kepastian penipuan. Tidak perlu mengunggah KTP untuk menyusun draft.

Deployment aplikasi hasil build tersedia melalui [panduan pribadi](deploy/GUIDE.txt)
di loopback port 8877 dengan database terpisah, filesystem container read-only,
user non-root, serta CSP. Konfigurasi inference GPU disiapkan tetapi model 70B
belum dijalankan pada Mac ini. Belum ada deployment publik, submission, uji
korban nyata, OCR, atau pengukuran manfaat. Lisensi lokal tersedia di `LICENSE`,
`LICENSES/`, dan `NOTICE`: kode asli Apache 2.0, dokumentasi/design asli CC BY 4.0,
dan input sintetis asli CDLA Permissive 2.0. Komponen pihak ketiga mempertahankan
lisensinya; kasus pengguna dan keluaran model pihak ketiga tidak otomatis menjadi
dataset berlisensi tersebut. Kesesuaian seluruh aturan submission masih perlu
diverifikasi.

## Penguatan dan evaluasi lanjutan

Evaluasi 12 kasus baru bilingual mencatat **10 draft diterima, satu model mismatch
ditolak, satu timeout/koneksi gagal**. Sumber dan transfer asli tetap utuh pada
semua kasus. Assistant memeriksa 26 ringkasan yang diterima: ditemukan arah
transfer hilang, pelaku janji berubah, nominal hilang, dan konteks klaim sebelumnya
hilang. Prompt `bilingual-v4` dan penanda review tambahan disiapkan dari temuan
ini. Dua kasus baru untuk prompt v4 gagal pada batas waktu provider; perbaikan
live saat itu belum terbukti. Pemeriksaan ulang yang direkam pada 11 Oktober menolak
satu respons Apertus karena kutipan berubah dan satu respons beridentitas
`aisingapore/Qwen-SEA-LION-v4-32B-IT` karena bukan Apertus. Tidak ada draft
dari dua respons tersebut yang tersimpan; input tetap utuh. Hasil awal tidak
diganti, dan tidak ada retry/fallback.

Fixture fiktif dan pemeriksa kontrak bisa dijalankan tanpa API/biaya:

```sh
cd src/laporsiap_BE
go run ./cmd/evaluate
```

Flag `-live` menjalankan satu permintaan per kasus dan mengirim teks fiktif
ke provider. Gunakan hanya saat memang ingin evaluasi tambahan; data dan hash
fixture dibekukan sebelum panggilan. `scripts/check-evidence.sh` memeriksa keadaan
database kasus evaluasi serta katalog model tanpa panggilan generasi AI.

Pilot dua/tiga orang masih memerlukan peserta nyata; [panduan fasilitator](docs/testing/pilot/facilitator-guide.txt)
dan recorder sudah tersedia. `scripts/analyze-pilot.py` memeriksa hasil pengamatan;
template kosong menghasilkan nol peserta, tanpa klaim manfaat. Deployment GPU
dan kesesuaian aturan lomba yang belum dapat diakses juga tetap perlu bukti.

Pemeriksa kini membekukan isi prompt beserta SHA-256 dan mencocokkannya dengan
backend yang sedang berjalan sebelum memanggil AI. Setiap draft Apertus yang
valid menyimpan `prompt_version` dan `prompt_sha256`. Status provider dan
verifikasi data setelah kegagalan tercatat tanpa token akses atau respons mentah.
Ringkasan terkini: `output/hardening/2026-10-11/live-recheck.json`.

Recorder pemeriksaan dua kasus (fixture dibuat lewat API lokal; persetujuan dan
panggilan AI lewat UI):

```sh
node scripts/demo/evaluate-recorded.mjs --check
node scripts/demo/evaluate-recorded.mjs
```

Keduanya tidak memanggil AI. `--live` merekam dua permintaan baru dan memakai
kuota/biaya provider; jangan menjalankannya berulang untuk mengejar hasil lulus.
Setiap sesi mempunyai folder sendiri, termasuk `LaporSiap_Failed_1.mp4` jika gagal.

## Perbaikan integrasi Apertus v5

Kontrak `bilingual-v5-source-bound` meminta AI menulis ID pesan dan ringkasan.
Aplikasi memasang kutipan langsung dari sumber lokal; ID/urutan/jumlah/schema
dan identitas model tetap diperiksa. Kontrak kutipan utuh untuk draft tersimpan
dan edit tidak dilonggarkan. Dua recheck development Inggris/Indonesia berhasil
dan direkam, tetapi review assistant menemukan next week hilang pada ringkasan
Inggris. Guard review waktu ditambahkan dalam dua bahasa; manusia tetap perlu
memeriksa semua entri. Hasil lama dan output AI asli dipertahankan.

Periksa konfigurasi/katalog/prompt runtime/kuota tanpa membuat draft:

```sh
sh scripts/check-provider.sh
```

Pemeriksaan katalog tidak menjamin routing completion. Provider tetap menolak
model alternatif dan tidak retry/fallback. Hasil terbaru dan batas interpretasi:
`output/evaluation/provider-v5-2026-10-11/verification.json`.

## CI otomatis

`.github/workflows/ci.yml` menyiapkan tiga job GitHub Actions: backend/race dengan
PostgreSQL terpisah, frontend architecture/unit/build dengan lockfile tetap, dan
13 regresi browser bilingual dengan database baru serta API key AI kosong.
Workflow berjalan pada push/pull request atau manual dispatch setelah kode
tersimpan di GitHub. Workflow tidak menjalankan evaluator live, deployment, atau
submission. Action dipatok ke commit SHA; hasil browser/log fiktif disimpan
selama tujuh hari. Eksekusi GitHub belum dilakukan dari workspace ini.

## Perbaikan waktu v6 dan cara pembuktian GPU

Draft AI baru mempertahankan next week/pekan depan dalam prompt dan menolak
ketidaksesuaian frasa waktu yang dikenali sebelum disimpan. Dua uji live terbaru
Inggris/Indonesia berhasil dengan waktu, kutipan, sumber dan transfer utuh;
review assistant mempertahankan kriteria empat entri. Dua contoh development
bukan bukti provider stabil. Hasil asli v5/hasil gagal sebelumnya tetap disimpan.
Bukti: output/evaluation/timing-v6-2026-10-11/verification.json.

[GPU-PROOF.txt](deploy/GPU-PROOF.txt) berisi rute VM4GPU dengan Compose atau
Runpod inference-only viaSSH, estimasi biaya, pin model/runtime resmi, alat
preflight/telemetry dan bukti akhir yang diperlukan. Default alat persiapan
adalah dry-run. Mac tidak mempunyai GPU NVIDIA; belum ada deploymentGPU,
penyewaan atau tagihan baru. Referensi kesiapan: output/gpu-preparation/2026-10-11.
