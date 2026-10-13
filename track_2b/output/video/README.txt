LaporSiap — opening & closing footage

SIAP DIGUNAKAN
LaporSiap-opening.mp4: 9 detik, logo dirakit dari lapisan 3D,
kartu bukti bergerak, lalu nama aplikasi dan tagline muncul.
LaporSiap-closing.mp4: 8 detik, logo 3D dan pesan penutup demo.

Keduanya: 1920 × 1080, 16:9, 30 fps, H.264 MP4 / yuv420p,
audio AAC stereo 48 kHz. Teks bahasa Inggris agar selaras dengan
voice-over demo Inggris. Tanpa narasi tambahan.

PENGGUNAAN
Letakkan opening sebelum rekaman demo dan closing setelahnya.
Masing-masing sudah memiliki fade masuk/keluar dan audio sendiri.
Jika ada voice-over yang tumpang tindih, turunkan volume klip sekitar
6–10 dB sesuai rekaman suara Anda, atau gunakan versi silent.
Tidak perlu menambahkan soundtrack lain agar klip dapat dipakai.

FILE TAMBAHAN
clean/LaporSiap-opening-silent.mp4
clean/LaporSiap-closing-silent.mp4
audio/LaporSiap-opening.wav
audio/LaporSiap-closing.wav
opening-poster.png / closing-poster.png: preview dari hasil video.
contact-sheet.jpg: cuplikan beberapa tahap gerakan kedua video.

DESAIN DAN AUDIO
SceneKit + Metal: geometri 3D sungguhan dengan extrusion,
lighting, bevel, dan animasi deterministik; logo mengikuti SVG aplikasi.
Warna teal, mint, amber; tipografi Manrope dan DM Sans dari proyek
(lisensi OFL tersedia di laporsiap_FE/public/fonts).
Audio disintesis orisinal: pad D mayor, nada lembut, dan whoosh terfilter.
Tidak menggunakan musik atau sampel audio dari pihak ketiga.
Level audio ditargetkan sekitar -18 LUFS, dengan fade di kedua ujung.

REBUILD (macOS, Swift/SceneKit/Metal, FFmpeg, Python 3)
bash scripts/video/build-brand-motion.sh

Source:
scripts/video/render-brand-motion.swift
scripts/video/synthesize-sound.py
scripts/video/build-brand-motion.sh

Script menjaga render silent yang sudah ada. Untuk merender perubahan
motion, pindahkan file silent sebelumnya ke nama lain terlebih dahulu,
lalu jalankan rebuild. Renderer memerlukan akses GPU Metal macOS;
lingkungan sandbox tertentu memerlukan izin untuk akses tersebut.
