# Catatan teknis MVP

## Batas kepercayaan

`Input` hanya menerima judul, status sintetis, teks, pilihan parser, serta
transfer. Parser memberi ID `M1…Mn` dan metadata asli. AI menerima hanya
ID/teks pesan, tanpa token kasus, isian transfer, atau bytes lampiran. Teks
pesan sendiri mungkin mengandung data pribadi; UI menjelaskannya dan meminta
persetujuan pada setiap permintaan AI.

Preview menampilkan ID/teks pesan yang sama dengan payload sumber provider.
Instruksi aplikasi ikut dikirim, tanpa isian transfer, metadata pengirim/waktu,
atau lampiran terpisah. Persetujuan UI direset ketika permintaan draft dimulai,
termasuk percobaan yang gagal, serta setelah perubahan bahasa/revisi. Batas 12 KB
di frontend diukur sebagai byte UTF-8 agar sama dengan validasi Go.

Prompt di `internal/infrastructure/apertus/prompt.json` memakai versi
`bilingual-v6-time-preserved` dan memilih instruksi berdasarkan `Accept-Language`.
AI hanya mengembalikan `source_id` dan `description` per pesan; versi aktif
sekarang `bilingual-v6-time-preserved` mempertegas preservasi waktu. Kutipan dipasang
oleh domain dari teks sumber lokal setelah validasi jumlah/urutan/ID; AI tidak
menyalin ulang evidence. Field evidence dalam respons provider justru ditolak.
Dalam mode redacted, adapter memasang teks masked dan use case memvalidasinya
lagi sebelum mengembalikan kutipan asli lokal. Tidak ada normalisasi atau
penerimaan kutipan salah dari hasil lama. Kutipan tersimpan tetap berbahasa asli.
Provider memakai endpoint chat completions, temperatur nol,
timeout 40 detik, dan tidak melakukan retry/fallback tersembunyi. Redirect
ditolak agar credential tidak diteruskan ke endpoint baru. Endpoint HTTPS
dapat dikonfigurasi; HTTP hanya untuk localhost/127.0.0.1.

Respons harus menyebut tepat `swiss-ai/apertus-v1.5-70b`, satu choice lengkap,
serta JSON dengan hanya `events`. Setiap event provider memiliki tepat
`source_id` dan `description`, berurutan dengan ID yang dikenal dan description
nonkosong ≤35 kata. Kontrak draft lengkap setelah binding masih memvalidasi
`source_id`, `description`, dan kutipan utuh persis. Edit pengguna juga tetap
dilarang mengubah ID atau kutipan.
Pemeriksaan identitas mengikuti metadata response pada
[schema Public AI](https://platform.publicai.co/api/~schemas); ini bukan audit
bobot fisik server. Respons tidak dapat membuat atau mengubah transfer.

Guard Go memberi tanda review bilingual untuk arah dana berlawanan, instruksi
AI dalam sumber, perubahan pending/gagal menjadi selesai, permintaan/janji yang
ditulis sebagai kejadian, dan nominal IDR/Rp yang tidak ada dalam kutipan. Guard
istilah dari riset juga dipertahankan. Regex ini tidak membuktikan entailment,
dapat memberi false positive/negative, dan tidak menangkap semua bentuk bahasa.
Regresi mencakup respons aktual batch recovery dan contoh sintetis positif/
negatif; belum merupakan evaluasi manusia independen atau bukti generalisasi.

Draft yang diterima menyimpan bahasa saat dibuat, waktu, latensi permintaan,
versi prompt, dan usage token jika provider menyediakannya. Metadata bahasa
adalah bahasa yang diminta, bukan hasil audit bahasa respons; edit pengguna bisa
mengubah bahasa isi. Pemeriksaan model memakai metadata API, bukan bukti bobot.
Mismatch mengembalikan diagnostik ID model terbatas, tanpa raw response atau
credential. Draft ditolak dan tidak ada fallback model otomatis.

## Penyimpanan dan revisi

Antarmuka memakai katalog English/Indonesia dengan default English. Preferensi
bahasa tersimpan di localStorage; token kasus tetap di sessionStorage.
Pergantian bahasa memperbarui label, pesan sistem, atribut bahasa dokumen,
format angka/tanggal UI, dan judul tab tanpa mengubah revisi atau isi kasus.
API tetap memakai nilai status dan kontrak input yang sama. Pesan validasi
API diterjemahkan saat ditampilkan oleh frontend.

Permintaan draft dan unduhan membawa `Accept-Language`. Unduhan memilih label
TXT/PDF dan header CSV di ZIP; draft baru memilih instruksi output AI.
Format API tanpa header tetap kompatibel dengan ekspor Indonesia. Ekspor
menerjemahkan label milik aplikasi sebelum memasukkan isian pengguna;
quote, description, rekening, catatan, dan snapshot JSON tetap asli.

PostgreSQL memiliki `cases`, `case_history`, `attachments`, dan
`schema_migrations`. Snapshot kasus berupa JSONB; file berupa BYTEA dengan
nama aman, MIME, ukuran, dan SHA-256. Snapshot mencakup sumber, transfer,
draft, daftar file, revisi, dan konfirmasi review. Riwayat mempertahankan
draft awal dan koreksi pada revisi berbeda.

Setiap mutasi memakai compare-and-swap revisi dalam transaksi. Mutasi snapshot,
riwayat, dan file bersifat atomik. Konflik mengembalikan HTTP 409. Jika pengguna
mengubah sumber saat AI berjalan, respons lama tidak dapat menimpa revisi baru.

Konfirmasi review membuat revisi baru yang ditandai diperiksa. Seluruh mutasi
selanjutnya menghapus konfirmasi. Ekspor memakai snapshot repeatable-read,
memeriksa review revisi tepat, lalu mengambil file dari snapshot database yang
sama. CSV mempertahankan nilai kosong dan melindungi sel berawalan formula.
ZIP tidak berisi token kasus.

Hubungan lampiran ke pesan/transfer dipilih pengguna dan disimpan dalam snapshot
JSONB. ID harus ada pada kasus yang sama, unik, dan sesuai revisi. Perubahan
hubungan membatalkan review. ID pesan bersifat positional; perubahan raw chat/
parser menghapus draft dan semua hubungan pesan agar tidak bergeser ke pesan
lain. Penghapusan transfer membersihkan hubungannya. Tidak ada OCR atau penilaian
keaslian berdasarkan hubungan ini.

PDF A4 dibuat di backend dengan fpdf dan DM Sans statis lokal berlisensi OFL.
ZIP menyertakan PDF di samping TXT/CSV/JSON dan lampiran asli. Font tidak mencakup
seluruh Unicode: karakter yang tidak didukung ditulis sebagai `[U+XXXX]` dengan
catatan, sedangkan sumber TXT/JSON tetap utuh. PDF diuji dengan sumber panjang
dan dirender untuk pemeriksaan tata letak; tidak ada layanan konversi eksternal.

## API

| Method | Route di bawah `/api` | Fungsi |
| --- | --- | --- |
| GET | `/health` | Koneksi PostgreSQL |
| GET | `/policy` | Batas retensi/AI dan label provider |
| POST | `/cases` | Input → kasus + token akses sekali pada respons |
| GET / PUT | `/cases/{id}` | Baca / koreksi input dengan revisi |
| DELETE | `/cases/{id}` | Konfirmasi revision + confirm_delete; hapus kasus/riwayat/file |
| POST | `/cases/{id}/access/rotate` | Ganti token; revoke yang lama tanpa membatalkan review |
| POST | `/cases/{id}/ai-preview` | Payload aktual setelah masking, hash consent, tanpa panggilan provider |
| POST | `/cases/{id}/draft` | Mode manual atau Apertus + external_consent |
| PUT | `/cases/{id}/draft` | Koreksi description; ID/kutipan tetap |
| POST | `/cases/{id}/review` | Konfirmasi sumber, transfer, dan keterbatasan |
| POST | `/cases/{id}/attachments` | Satu file multipart, X-Case-Revision |
| GET / DELETE | `/cases/{id}/attachments/{file}` | Unduh / hapus milik kasus |
| PUT | `/cases/{id}/attachments/{file}` | Hubungan source_ids/transfer_ids + revision |
| GET | `/cases/{id}/export?format=zip` | ZIP, PDF, atau TXT default, setelah review |
| GET | `/cases/{id}/history` | 50 snapshot revisi terakhir |

Semua route kasus setelah create memerlukan Bearer token kasus. Akses salah
menghasilkan 404 yang sama dengan kasus tidak ditemukan. Respons memakai
`Cache-Control: no-store`, nosniff, dan no-referrer. Credential/teks pribadi
tidak dicetak ke log aplikasi. JSON dibatasi 256 KB dan menolak field tambahan.

## Verifikasi dan batas interpretasi

Unit test memeriksa total/status/unknown, tanggal/parser TXT, kontrak AI,
heuristik yang pernah gagal, formula CSV, isi ZIP, dan nama file. Provider
mock memeriksa model mismatch, kutipan palsu, injection schema, respons
terpotong, HTTP 504, header/path, serta tidak bocornya transfer/metadata ke AI.

Uji PostgreSQL memakai database terpisah. Ia mencakup akses baca/tulis kasus,
riwayat dan file dengan token kasus lain, revisi lama, semua jenis perubahan
setelah review, export file, kegagalan provider, dan race respons AI vs input.
Uji UI memakai Chrome dengan demo sintetis: input, consent, manual draft,
koreksi, file, review, ZIP, reload sesi, lock ekspor setelah perubahan, serta
viewport desktop/mobile. Uji UI tidak mengirim permintaan AI live.

Regresi browser juga memeriksa persetujuan ulang setelah respons AI sukses
dan gagal memakai respons simulasi, serta gerbang byte untuk teks multibyte.
Keempat test browser lulus pada verifikasi ulang sesi awal.

Setelah fitur bahasa dan perbaikan UI, 7/7 test browser lulus. Pemeriksaan
mencakup preferensi default/persisten, validasi API dalam bahasa terpilih,
review dan unduhan TXT dua bahasa, preservasi input/koreksi saat bahasa
diganti, serta lebar viewport 320/390/768/1440 px tanpa overflow halaman.
Uji Go tambahan memeriksa preservasi sumber/unknown pada ekspor bilingual
dan kesesuaian delapan fixture testing dengan parser serta total aktual.

Penguatan 10 Oktober: **9/9 test browser lulus**, termasuk preview/consent bahasa,
metadata bahasa draft, pesan guard spesifik, hubungan bukti yang bertahan setelah
reload, gerbang PDF dua bahasa, pembatalan review, dan pembersihan hubungan usang.
Go dengan PostgreSQL dan race detector lulus setelah penambahan fitur dan
diagnostik model. Build TypeScript/Vite lulus. Uji browser ini memakai manual
draft atau respons mock; keberhasilannya tidak membuktikan availability live AI.

Satu probe live Go melalui backend memakai tiga pesan sintetis N03. Respons
memiliki metadata Apertus yang diminta, tiga source/quote valid, serta
ringkasan transfer gagal, permintaan mencoba ulang, dan belum adanya transfer
berhasil. Review ini oleh assistant, bukan evaluasi manusia independen.
Tidak ada retry atau batch riset baru. Hasil ada di
`evidence/apertus-integration-20261009.json`.

Bukti 10 Oktober `evidence/apertus-english-20261010.json` menyimpan satu draft
English asli yang diterima, lima kutipan persis, usage token, latensi, dan
estimasi tarif inferensi dari katalog API. Pemeriksaan makna dilakukan oleh
assistant. Rekaman keempat berhenti karena assertion skrip yang terlalu sempit
pada kata "demands"; draft tetap valid terhadap schema/kutipan/model. Semua
rekaman gagal dipertahankan, tanpa fallback atau retry dalam rekaman. Sampel
ini tidak membuktikan akurasi umum, manfaat pengguna, atau uptime provider.

Rekaman live berikutnya `LaporSiap_Success_1` menyelesaikan delapan langkah
dengan satu panggilan Apertus asli, review, PDF/TXT/ZIP, dan riwayat revisi.
Metadata draft: 545 input / 366 output token, latensi 3.535 ms. JSON laporan
rekaman dan draft asli disimpan bersama MP4/WebM di
`output/simulation/2026-10-10/apertus/english/`. Isi paket ekspor dan PDF live
divalidasi; empat rekaman gagal sebelumnya tetap dipertahankan. Ini tetap
demonstrasi sintetis dan assistant review, bukan studi kegunaan independen.

Keterbatasan: belum ada uji korban nyata, OCR, pemulihan token sesi, akun,
streaming, antrean multi-user, pengukuran biaya operasional/waktu review,
audit keamanan independen, atau deployment model mandiri. Penyimpanan disk
belum dienkripsi oleh aplikasi. Static frontend kini dapat disajikan Go pada
origin API yang sama melalui FRONTEND_DIR; filesystem root membatasi akses di
luar hasil build, file tersembunyi dan directory listing ditolak. Semua font
disajikan lokal. Deployment/publication publik tetap tahap terpisah.

## Lifecycle, kapasitas dan privasi — 10 Oktober 2026

Migrasi 002 menambahkan expires_at (30 hari), access_expires_at (7 hari) dan
meter ai_attempts. Kasus yang sudah ada dibackfill dari created_at; riwayat
dan file tidak dihapus migrasi. Startup/jadwal hourly menghapus hanya kasus
dengan retensi berakhir; access expiry sendiri tidak mempercepat penghapusan.
DELETE berkunci token/revisi dan cascades menghapus history/files. Token rotation
memperbarui hash/expiry, tanpa content revision atau membatalkan review. Meter AI
kasus dihapus relasinya dan dipertahankan anonim hingga 48 jam untuk batas global.

Advisory lock PostgreSQL menyerialisasi pemeriksaan kapasitas+write dan reservasi
AI. Batas 256 MiB menghitung JSON kasus/riwayat tersimpan serta bytes lampiran;
bukan batas filesystem/WAL/backup. Kuota AI bergulir 24 jam, persisten lintas
restart; lease aktif maksimum dua selama 45 detik dengan timeout request 40 detik.
Percobaan gagal tetap dihitung. Tidak ada retry/fallback atau antrean AI otomatis.
Rate request per IP berlaku pada proses lokal dan hilang setelah restart.

Preview adalah fungsi deterministic dari case ID/revisi/bahasa dan pengaturan
masking. Hash berubah jika payload berubah; consent eksplisit dengan hash lama
ditolak sebelum reservasi/panggilan AI. Legacy consent tanpa privacy settings
menggunakan masking default. Hanya ID/teks pesan dikirim. Placeholder tidak
memuat nilai asli. File, transfer, pengirim/waktu metadata, judul, dan token kasus
tetap lokal. Source quote ter-mask harus lolos contract sebelum quote asli
dipulihkan dari sumber lokal. Deskripsi tetap hasil model dan memerlukan review.
Masking regex bukan anonimisasi sempurna; teks bebas/alamat/nama lain perlu
istilah tambahan dan preview. Kutipan/ekspor asli tetap dapat memuat data pribadi.

Prompt v4 mempertahankan tindakan/negasi, pelaku di dalam pesan, nominal utama
dan konteks klaim. Heuristik tambahan menandai pengiriman uang yang hilang,
peran perekrut yang hilang, serta nominal rupiah hilang. Deteksi kosakata terbatas
dapat menghasilkan false positive dan tidak membuktikan entailment. Koreksi
tetap keputusan manusia. Dua panggilan v4 pada input baru timeout, sehingga
peningkatan kualitas model belum dapat dinilai.


## Bukti prompt dan provider — 11 Oktober 2026

`PromptSnapshot()` mengembalikan salinan prompt embedded. `PromptSHA256()` dan
`PromptVersion` tersedia pada `/api/policy`; draft valid menyimpan keduanya.
Evaluator CLI maupun recorder mencocokkan hash backend dengan snapshot lokal
sebelum panggilan AI. Ini mengidentifikasi isi prompt, bukan identitas bobot model.
Evaluator menyimpan diagnostik model/rule yang terbatas, tidak respons provider
mentah, lalu membaca kembali kasus untuk memeriksa sumber, transfer dan state
setelah gagal. `state_verified` hanya true jika semua pemeriksaan selesai lulus.

Pada dua panggilan v4 sebelumnya lewat UI, respons Inggris beridentitas Apertus ditolak
karena kutipan berbeda; respons Indonesia beridentitas
`aisingapore/Qwen-SEA-LION-v4-32B-IT` ditolak karena model tidak cocok. Tidak ada
retry atau draft tersimpan. Penyebab routing internal provider dan detail karakter
kutipan yang berbeda belum diketahui; tidak dilakukan normalisasi kutipan agar
respons dapat diterima. Katalog tetap mencantumkan Apertus, tetapi ini tidak
menjamin model tersebut dipakai dalam sebuah respons. Tidak ada output v4 yang
bisa dinilai secara semantik pada pemeriksaan ini.


## Clean Architecture dan Redux — 11 Oktober 2026

Domain backend tidak mengimpor HTTP, pgx, fpdf, atau layer luar. Model, parser,
validasi lampiran/tautan, masking dan kontrak kronologi berada di `internal/domain`.
`internal/application.Service` menangani create/update/draft/review/attachment/
lifecycle/export melalui `CaseRepository`, `Generator`, dan `ReportRenderer`.
Pembuatan agregat, kenaikan revisi dan invalidasi/konfirmasi review berada di use
case; repository menerima agregat berikutnya dan expected revision untuk CAS.
Transaksi, capacity lock, lease AI, expiry/auth predicates dan snapshot repeatable
read tetap berada di adapter PostgreSQL. Renderer TXT/CSV/PDF/ZIP dan provider
Apertus terpisah. Handler hanya menangani autentikasi/decode HTTP, batas multipart,
pemetaan error, serta header respons. `cmd/api` menghubungkan adapter dengan use
case; command evaluasi/maintenance juga memakai package sesuai tanggung jawabnya.

Seluruh regresi lama dipindahkan ke `internal/integration`, dengan nama kompatibel
pada helper khusus `_test.go`; tidak ada facade monolitik pada kode produksi.
Unit use case baru memakai fake port tanpa HTTP/DB/provider. Test arsitektur Go
memeriksa import production sehingga core tidak bergantung ke adapter.

Frontend memiliki `domain`, `application`, `infrastructure`, `presentation`, dan
composition root `app/bootstrap.ts`. `CaseWorkflow` bergantung ke interface untuk
repository kasus, session vault dan file saver. Use case tidak mengimpor React,
Redux, atau akses browser langsung. Adapter HTTP menerima vault lewat constructor
untuk mendapatkan token saat mengirim permintaan. Respons bertoken diserap vault
sebelum data kasus dikembalikan ke presentation.

Redux Toolkit 2.13.0 dan React Redux 9.3.0 menyimpan state workspace dan locale.
Controller React menggunakan selector/action bertipe. Respons kasus diterapkan
sebagai satu transisi reducer; state privasi/review direset ketika kasus/revisi
berubah. Rotasi akses memperbarui vault dan access epoch, membatalkan consent dan
preview. Tidak ada token, File/Blob, fungsi callback atau AbortController pada
state/action Redux. File diunggah/diekspor lewat use case dan adapter browser.
Session vault dapat bekerja dalam memory bila browser menolak sessionStorage.
DevTools dan logger tidak diaktifkan, dan state kasus tidak dipersistensikan.

`pnpm check:architecture` memeriksa import TypeScript dan akses browser langsung
pada core. Unit test frontend menjalankan modul TypeScript produksi untuk
verifikasi pemisahan credential, kegagalan delete, dan reset privasi/review. Unit
Node berada di `tests/unit` dan dikecualikan dari discovery Playwright.

## Kontrak v5 dan pemeriksaan provider — 11 Oktober 2026

`sh scripts/check-provider.sh` memeriksa konfigurasi, GET katalog model, hash
prompt backend lokal, dan snapshot kuota/capacity PostgreSQL. Pemeriksaan tidak
memanggil completion, tidak memigrasi schema, tidak membuat kasus dan tidak
menampilkan key/URL privat/token. Key yang terisi dengan konfigurasi tidak valid
juga menghentikan startup API. Tidak ada parameter routing privat yang ditebak;
request tetap menyebut model eksplisit sesuai dokumentasi resmi Public AI.
Katalog tidak membuktikan routing completion atau identitas bobot.

Dua recheck v5 pada sumber development v4 diterima lewat UI (EN/ID HTTP 200),
model metadata Apertus sesuai, empat kutipan persis termasuk newline akhir,
serta sumber/transfer tetap utuh. Latensi provider 2051/2847 ms; usage total
1325 input dan 196 output token. Tidak ada retry atau fallback. Rekaman Success_1
berarti kontrak diterima, bukan semua makna sudah sempurna. Review assistant
menemukan EN M2 menghilangkan next week; ID mempertahankan pekan depan.

Guard relatif waktu kini membandingkan kelompok padanan tomorrow/besok, next
week/pekan depan/minggu depan, next month/bulan depan, today/hari ini dan
yesterday/kemarin, lalu menampilkan relative_time_requires_review jika berbeda.
Tidak menghitung tanggal, memperbaiki ringkasan otomatis, atau membuktikan
entailment. Guard dihitung saat membuat/mengedit draft; draft historis tidak
ditulis ulang dan tetap memerlukan pemeriksaan semua entri. Output live asli
dan review assistant disimpan terpisah. Dua contoh bukan holdout independen
atau bukti kestabilan provider; penyebab mismatch sebelumnya masih belum
diketahui. Bukti: output/evaluation/provider-v5-2026-10-11/verification.json.

## CI dan batas verifikasinya

Workflow lokal `.github/workflows/ci.yml` membagi backend (Go/race/domain/use
case/architecture/PostgreSQL), frontend (23-file architecture/Redux unit/build),
dan browser (13 regresi). Dua job yang memakai database mendapat instance
service PostgreSQL terpisah. Tidak membaca .env lokal atau memakai key Apertus.
AI test memakai server lokal/mock atau respons browser terkontrol. Readiness
check dibatasi 45 detik dan hanya membaca health/HTML lokal. Actions dipatok
commit SHA yang dicek dari repositori resmi; Node 24.15.0, pnpm 11.1.3 dan Go
dari go.mod konsisten dengan build. Frozen lockfile wajib.

Actionlint memvalidasi sintaks/workflow lokal; ini tidak menggantikan eksekusi
runner GitHub/Ubuntu. Semua 13 browser lokal lulus (23,4s); health check lokal
lulus setelah izin localhost. Percobaan health dalam sandbox tidak dapat
terhubung dan timeout tanpa menulis kasus atau memanggil AI. Tidak ada push
atau GitHub CI run pada pekerjaan ini. Instruksi runtime/action mengikuti
[setup-node resmi](https://github.com/actions/setup-node) dan
[Playwright CI](https://playwright.dev/docs/ci-intro).

## Perbaikan waktu v6 dan persiapan GPU — 11 Oktober 2026

Prompt v6 memberi instruksi dan contoh waktu bersama janji/pelaku/nominal serta
negasi. Domain ValidateGeneratedTiming membandingkan lima kelompok frasa waktu
yang sudah didokumentasikan dan menolak kehilangan/perubahan/penambahan sebelum
draft AI tersimpan. Pemeriksaan dipakai pada binding respons provider dan use
case restoration terhadap sumber masked. Tidak mengubah ringkasan otomatis,
menghitung tanggal, menormalisasi kutipan, retry, atau menerima model lain.
Kontrak edit manusia tetap menyimpan flag waktu untuk koreksi/review; tidak
menulis ulang snapshot historis. Kata/frasa lain dan hubungan waktu dengan
kejadian berbeda belum dapat diverifikasi sepenuhnya dengan regex ini.

Dua live v6 EN/ID diterima pada model metadata Apertus dan mempertahankan next
week/pekan depan pada janji komisi. Review assistant empat entri mempertahankan
kriteria yang dinyatakan: pengiriman 264000, refund belum diterima, perekrut menjanjikan
63000, pekan depan dan komisi belum dibayar. Kutipan/sumber/transfer tetap utuh.
Tidak ada klaim holdout, akurasi umum atau kestabilan provider dari dua recheck.
Output v5 yang menghilangkan waktu dan seluruh kegagalan lama tetap disimpan.

GPU: deploy/gpu-lock.json membekukan digest manifest runtime resmi dan revisi
checkpoint 70B yang benar-benar dibaca dari registry/metadata publik. Compose
menambah dtype BF16/max_num_seqs 2 dan mount alat read-only; origin 8878 sesuai
tunnel GPU default, dengan project dan database terpisah. Profil direncanakan
4xH100 80GB, RAM 256GiB, disk kosong 250GiB dan driver 580+ untuk image CUDA 13.
Angka ini adalah profil konservatif, bukan hasil ukur kapasitas model berjalan.
preflight/serve/init-env/collect ada di scripts/gpu. Launcher/init default dry-run;
probe/collector di Mac mencatat GPU tidak tersedia dan tetap deployment false.
Collector mengekspor field container terpilih tanpa Env/key/rawinspect, memeriksa
cache pinned, health/katalog model dan telemetry; tidak memanggil completion atau
menganggap dirinya sebagai bukti lengkap. Pod Runpod tidak mendukung Compose;
panduan memisahkan rute VM penuh dan inference-only melalui SSH.

Rencana, tarif, sumber, prosedur pelaksanaan dan bukti review/export/biaya yang
masih diperlukan: deploy/GPU-PROOF.txt. Tidak ada GPU disewa, weights/image 9GB
diunduh, layanan remote dibuat, atau akun/model terms diterima otomatis. Source
readiness/pin/struktur Compose dan empat unit GPU lulus secara lokal; kompatibilitas
flags runtime/GPU sebenarnya harus diperiksa pada server sebelum serving.
