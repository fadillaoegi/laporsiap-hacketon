# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: smoke.spec.ts >> synthetic case: manual review, attachment, export, and invalidation after input edit
- Location: tests/smoke.spec.ts:8:1

# Error details

```
Error: expect(locator).toBeVisible() failed

Locator: getByText('bukti-referensi.pdf', { exact: true })
Expected: visible
Timeout: 5000ms
Error: element(s) not found

Call log:
  - Expect "toBeVisible" getByText('bukti-referensi.pdf', { exact: true }) with timeout 5000ms
  - waiting for getByText('bukti-referensi.pdf', { exact: true })

```

```yaml
- banner:
  - text: LaporSiap PROTOTIPE LOKAL Language / Bahasa
  - combobox "Language / Bahasa":
    - option "English"
    - option "Indonesia" [selected]
- complementary:
  - text: RUANG PERSIAPAN
  - navigation "Tahapan menyiapkan laporan":
    - button "Bahan kasus Teks, transfer, dan bukti"
    - button "Susun draft Apertus atau manual"
    - button "03 Periksa & koreksi Cocokkan dengan sumber"
    - button "04 Ekspor bahan Untuk kanal resmi"
  - strong: Anda tetap memegang kendali.
  - paragraph: Draft membantu merapikan bahan. Anda memeriksa isinya dan mengirim laporan melalui kanal resmi.
  - text: Kasus dalam sesi ini
  - combobox "Kasus dalam sesi ini":
    - option "Pilih kasus" [disabled]
    - option "Tawaran tugas daring dan permintaan top-up" [selected]
  - button "+ Buat kasus baru"
  - text: Tanpa unggah KTP. Tanpa pengiriman laporan otomatis.
- main:
  - paragraph:
    - strong: Perlu segera melapor?
    - text: Anda dapat langsung mengakses kanal resmi, tanpa menunggu draft AI.
  - link "Buka IASC":
    - /url: https://iasc.ojk.go.id/
  - text: LANGKAH 01 / 04
  - heading "Mulai dari bahan yang Anda punya." [level=1]
  - paragraph: Percakapan yang belum rapi bisa menjadi titik awal. Data yang belum diketahui boleh tetap kosong.
  - text: Revisi 3 · Perlu review DATA SIMULASI — identitas, percakapan, dan transaksi ini rekaan.
  - alert:
    - text: PDF tidak lengkap atau struktur dasarnya tidak valid. Pilih PDF yang dapat dibuka.
    - button "Muat ulang kasus tersimpan"
  - heading "Percakapan & cerita awal" [level=2]
  - paragraph: Tempel teks atau impor file TXT.
  - group:
    - text: Judul kasus
    - textbox "Judul kasus":
      - /placeholder: "Contoh: Tawaran pekerjaan daring dan top-up"
      - text: Tawaran tugas daring dan permintaan top-up
    - text: Format percakapan
    - combobox "Format percakapan":
      - option "Teks biasa · satu pesan per baris"
      - 'option "TXT · [DD/MM/YYYY, HH:MM] Pengirim: pesan" [selected]'
    - text: Impor TXT
    - button "Impor TXT"
    - text: Percakapan
    - textbox "Percakapan":
      - /placeholder: "Tempel pesan di sini.\nPengirim dan waktu tidak ditebak dari teks biasa."
      - text: "[08/10/2026, 09:10] Nadia: Halo Raka, kami menawarkan tugas memberi like dengan janji komisi Rp50.000 per hari. [08/10/2026, 09:14] Nadia: Untuk lanjut tugas, silakan setor deposit Rp180.000 ke ****4821. [08/10/2026, 09:28] Raka: Saya baru mengirim Rp180.000. Waktu transfernya belum saya catat. [08/10/2026, 09:35] Nadia: Sekarang tambah deposit Rp420.000 agar saldo bisa ditarik. [08/10/2026, 09:40] Raka: Status setoran tambahan masih dalam proses. Saya belum melihat mutasi rekening."
    - text: "Teks biasa: satu pesan per baris, tanpa metadata pengirim/waktu. TXT berheader: format yang dipilih harus persis cocok; baris lanjutan menjadi bagian pesan. Maksimal 100 pesan / 32 KB."
    - checkbox "Ini data simulasi" [checked]
    - text: Ini data simulasi
  - text: Rp
  - heading "Catatan transfer" [level=2]
  - paragraph: Isi berdasarkan catatan Anda, bukan hasil AI.
  - text: 2 catatan
  - group:
    - strong: Transfer 01
    - button "Hapus catatan"
    - text: Bank pengirim
    - textbox "Bank pengirim":
      - /placeholder: Belum diketahui
      - text: Bank Aruna
    - text: Rekening pengirim
    - textbox "Rekening pengirim":
      - /placeholder: Belum diketahui
      - text: "****7392"
    - text: Bank penerima
    - textbox "Bank penerima":
      - /placeholder: Belum diketahui
      - text: Bank Lestari
    - text: Rekening penerima
    - textbox "Rekening penerima":
      - /placeholder: Belum diketahui
      - text: "****4821"
    - text: Nominal (rupiah)
    - spinbutton "Nominal (rupiah)": "180000"
    - text: Status transfer
    - combobox "Status transfer":
      - option "Berhasil" [selected]
      - option "Dalam proses"
      - option "Gagal"
      - option "Belum dipastikan"
    - text: Tanggal transfer
    - textbox "Tanggal transfer Boleh kosong. Waktu chat bukan tanggal transfer."
    - text: Boleh kosong. Waktu chat bukan tanggal transfer. Catatan
    - textbox "Catatan":
      - /placeholder: Koreksi, perbedaan keterangan, dll.
      - text: Tanggal transfer belum dicatat; jangan memakai waktu chat.
  - group:
    - strong: Transfer 02
    - button "Hapus catatan"
    - text: Bank pengirim
    - textbox "Bank pengirim":
      - /placeholder: Belum diketahui
      - text: Bank Aruna
    - text: Rekening pengirim
    - textbox "Rekening pengirim":
      - /placeholder: Belum diketahui
      - text: "****7392"
    - text: Bank penerima
    - textbox "Bank penerima":
      - /placeholder: Belum diketahui
      - text: Bank Lestari
    - text: Rekening penerima
    - textbox "Rekening penerima":
      - /placeholder: Belum diketahui
      - text: "****4821"
    - text: Nominal (rupiah)
    - spinbutton "Nominal (rupiah)": "420000"
    - text: Status transfer
    - combobox "Status transfer":
      - option "Berhasil"
      - option "Dalam proses" [selected]
      - option "Gagal"
      - option "Belum dipastikan"
    - text: Tanggal transfer
    - textbox "Tanggal transfer Boleh kosong. Waktu chat bukan tanggal transfer."
    - text: Boleh kosong. Waktu chat bukan tanggal transfer. Catatan
    - textbox "Catatan":
      - /placeholder: Koreksi, perbedaan keterangan, dll.
      - text: Belum dipastikan berhasil.
  - button "+ Tambah catatan transfer"
  - heading "Lampiran bukti" [level=2]
  - paragraph: Referensi gambar atau PDF. Belum termasuk OCR.
  - strong: Pilih file bukti
  - text: PNG, JPEG, PDF · maks. 5 MB/file · 10 lampiran
  - button "Pilih file bukti PNG, JPEG, PDF · maks. 5 MB/file · 10 lampiran"
  - paragraph: Gambar maksimal 24 juta piksel.
  - text: Disimpan di PostgreSQL lokal.
  - button "Simpan & lanjut"
  - complementary:
    - text: MULAI SECARA BERTAHAP
    - heading "Tidak harus lengkap sekarang." [level=3]
    - paragraph: Pertahankan pesan asli. Isi rincian transfer yang Anda ketahui, lalu tambahkan bukti bila tersedia.
    - list:
      - listitem: Teks percakapan atau ekspor TXT
      - listitem: Nominal, rekening, dan status transfer
      - listitem: Bukti transfer atau komunikasi
    - heading "Jelas sebelum dikirim." [level=3]
    - paragraph: Menyimpan kasus belum mengirim teks ke AI. Pada langkah berikutnya, Anda memilih draft manual atau menyetujui pengiriman teks ke Public AI.
    - paragraph: Lampiran hanya sebagai referensi. Gambar dan PDF belum dibaca otomatis.
  - text: LaporSiap · Bahan bantu, dengan sumber dan review manusia. Development lokal · Data tersimpan di perangkat ini
```

# Test source

```ts
  1  | import { test, expect } from '@playwright/test'
  2  | import path from 'node:path'
  3  | 
  4  | test.beforeEach(async ({ page }) => {
  5  |   await page.addInitScript(() => localStorage.setItem('laporsiap-language', 'id'))
  6  | })
  7  | 
  8  | test('synthetic case: manual review, attachment, export, and invalidation after input edit', async ({ page }) => {
  9  |   const errors: string[] = []
  10 |   page.on('pageerror', e => errors.push(e.message))
  11 |   await page.goto('/')
  12 |   const officialChannel = page.getByRole('link', { name: 'Buka IASC', exact: true })
  13 |   await expect(officialChannel).toBeVisible()
  14 |   await expect(officialChannel).toHaveAttribute('href', 'https://iasc.ojk.go.id/')
  15 |   await page.screenshot({ path: path.resolve('../docs/screenshots/01-input-desktop.png'), fullPage: true })
  16 |   await page.getByRole('button', { name: 'Coba contoh', exact: true }).click()
  17 |   await page.getByRole('button', { name: 'Simpan & lanjut' }).click()
  18 |   const ai = page.getByRole('button', { name: 'Susun draft dengan Apertus', exact: true })
  19 |   await expect(ai).toBeDisabled()
  20 |   await page.getByRole('checkbox', { name: /Saya setuju mengirim/ }).check()
  21 |   await expect(ai).toBeEnabled()
  22 |   // UI path remains usable without making a live provider request.
  23 |   await page.getByRole('button', { name: 'Gunakan draft manual' }).click()
  24 |   await expect(page.getByText('DRAFT MANUAL', { exact: true })).toBeVisible()
  25 |   await expect(page.getByText('Rp 180.000', { exact: true }).first()).toBeVisible()
  26 |   const first = page.getByRole('textbox', { name: 'Kronologi M1', exact: true })
  27 |   await first.fill('Pengirim menawarkan tugas memberi like dengan janji komisi Rp50.000 per hari.')
  28 |   await expect(page.getByRole('checkbox', { name: 'Saya telah mencocokkan entri ini dengan sumber' }).first()).toBeDisabled()
  29 |   await page.getByRole('button', { name: 'Simpan koreksi' }).click()
  30 |   await expect(page.getByText('Koreksi tersimpan.', { exact: false })).toBeVisible()
  31 |   await page.getByRole('button', { name: /Bahan kasus/ }).click()
  32 |   await page.locator('input[type=file][accept="image/png,image/jpeg,application/pdf"]').setInputFiles({ name: 'bukti-referensi.pdf', mimeType: 'application/pdf', buffer: Buffer.from('%PDF-1.4\nSynthetic demonstration attachment\n') })
> 33 |   await expect(page.getByText('bukti-referensi.pdf', { exact: true })).toBeVisible()
     |                                                                        ^ Error: expect(locator).toBeVisible() failed
  34 |   await page.getByRole('button', { name: /Periksa & koreksi/ }).click()
  35 |   const checks = page.getByRole('checkbox', { name: 'Saya telah mencocokkan entri ini dengan sumber' })
  36 |   await expect(checks).toHaveCount(5)
  37 |   for (const checkbox of await checks.all()) await checkbox.check()
  38 |   await page.getByRole('checkbox', { name: /Saya telah memeriksa rekening/ }).check()
  39 |   await page.getByRole('checkbox', { name: /Saya telah memeriksa tanda review/ }).check()
  40 |   await page.screenshot({ path: path.resolve('../docs/screenshots/02-review-desktop.png'), fullPage: true })
  41 |   await page.getByRole('button', { name: 'Konfirmasi & lanjut' }).click()
  42 |   await expect(page.getByRole('button', { name: 'Unduh paket ZIP' })).toBeEnabled()
  43 |   await page.screenshot({ path: path.resolve('../docs/screenshots/03-export-desktop.png'), fullPage: true })
  44 |   const event = page.waitForEvent('download')
  45 |   await page.getByRole('button', { name: 'Unduh paket ZIP' }).click()
  46 |   const file = await event
  47 |   expect(file.suggestedFilename()).toMatch(/laporsiap-revisi-\d+\.zip/)
  48 |   await file.saveAs(path.resolve('../tmp/browser-demo.zip'))
  49 |   await page.reload()
  50 |   await expect(page.getByRole('button', { name: 'Unduh paket ZIP' })).toBeEnabled()
  51 |   await page.getByRole('button', { name: /Bahan kasus/ }).click()
  52 |   await page.getByRole('textbox', { name: 'Judul kasus', exact: true }).fill('Simulasi setelah koreksi')
  53 |   await page.getByRole('button', { name: /Ekspor bahan/ }).click()
  54 |   await expect(page.getByRole('button', { name: 'Unduh paket ZIP' })).toBeDisabled()
  55 |   await page.getByRole('button', { name: /Bahan kasus/ }).click()
  56 |   await page.getByRole('button', { name: 'Simpan & lanjut' }).click()
  57 |   await page.getByRole('button', { name: /Ekspor bahan/ }).click()
  58 |   await expect(page.getByRole('button', { name: 'Unduh paket ZIP' })).toBeDisabled()
  59 |   await page.setViewportSize({ width: 390, height: 844 })
  60 |   await page.getByRole('button', { name: /Periksa & koreksi/ }).click()
  61 |   await page.screenshot({ path: path.resolve('../docs/screenshots/04-review-mobile.png'), fullPage: true })
  62 |   expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true)
  63 |   expect(errors).toEqual([])
  64 | })
  65 | 
```