SIMULASI LAPORSIAP — 10 OKTOBER 2026

REKAMAN
LaporSiap_Success_1.mp4 — simulasi lengkap yang berhasil.
LaporSiap_Failed_1.mp4 — percobaan pertama, berhenti pada langkah lampiran.
Keduanya rekaman layar 1920 x 1080 / 30 fps, tanpa voice-over.
Versi WebM adalah rekaman asli dari browser.

HASIL
Simulasi memakai browser otomatis dan API/PostgreSQL lokal yang sesungguhnya.
Data kasus, percakapan, identitas, rekening, dan PDF seluruhnya rekaan.
Alur: muat contoh, simpan bahan, unggah PDF referensi, susun draft manual,
sunting 5 entri, cocokkan sumber, periksa transfer dan data kosong,
konfirmasi review, unduh ZIP dan TXT, lalu tampilkan riwayat.
Semua 8 langkah pada percobaan sukses lulus.

Rp180.000 berstatus berhasil; Rp420.000 berstatus pending.
Total transfer berhasil yang dilaporkan pengguna adalah Rp180.000.
Tanggal transfer tetap belum diketahui dan tidak diambil dari waktu chat.
Paket ZIP diperiksa integritas dan isinya: 5 sumber, 5 entri draft manual,
revisi 5 yang telah dikonfirmasi, dua catatan transfer, dan PDF sintetis.

Rekaman ini menggunakan jalur manual. Tidak ada panggilan AI live ke Apertus,
respons AI simulasi, atau pengiriman laporan ke IASC.
Teks petunjuk dan penanda kursor ditambahkan oleh perekam, bukan fitur aplikasi.

MENGAPA PERCOBAAN PERTAMA BERHENTI?
Skrip perekam mencari tombol dengan nomor "01 Bahan kasus". Sesudah bahan
disimpan, nomor langkah berubah menjadi ikon centang; tombolnya masih ada.
Skrip tidak menemukan tombol, menunggu 12 detik, kemudian menghentikan
rekaman segera setelah timeout. Pencarian tombol diperbaiki agar berdasarkan
nama "Bahan kasus". Tidak ada error dari API atau validasi kasus pada kejadian ini.
Percobaan baru direkam terpisah; rekaman gagal tetap disimpan utuh.

FILE HASIL SIMULASI
LaporSiap-simulasi.zip — paket hasil ekspor aplikasi.
LaporSiap-kronologi.txt — kronologi yang telah disunting dan direview.
percakapan-sintetis.txt — pesan asli yang dapat diimpor kembali.
bukti-sintetis.pdf — lampiran referensi rekaan, bukan bukti bank.
LaporSiap_Success_1.json / LaporSiap_Failed_1.json — langkah, waktu, dan status.
01...05 PNG — screenshot tahap simulasi.
preview-simulasi.jpg — cuplikan dari rekaman sukses.

IKUTI SENDIRI
1. Buka http://127.0.0.1:5175 dan pilih bahasa Indonesia.
2. Klik Coba contoh, lalu Simpan & lanjut.
3. Kembali ke Bahan kasus dan unggah bukti-sintetis.pdf bila ingin menampilkan lampiran.
4. Pilih Susun draft, lalu Gunakan draft manual.
5. Cocokkan dan sunting tiap entri dengan pesan sumber, lalu Simpan koreksi.
6. Centang review 5 sumber, catatan transfer, dan data belum lengkap.
7. Klik Konfirmasi & lanjut, lalu unduh ZIP dan TXT.

SOURCE PEREKAM
scripts/demo/record-simulation.mjs
Persiapan dari root proyek (Chrome lokal, Node, Python, FFmpeg tersedia):
  cd laporsiap_FE
  PLAYWRIGHT_BROWSERS_PATH=../tmp/playwright pnpm exec playwright install ffmpeg
  cd ..
  node scripts/demo/record-simulation.mjs

Browser/GPU dan akses layanan lokal bisa membutuhkan izin di lingkungan sandbox.
Skrip menghentikan rekaman jika ada error aplikasi, API, unduhan, atau pengecekan.
Nomor berkas bertambah otomatis sehingga rekaman gagal tidak ditimpa.
