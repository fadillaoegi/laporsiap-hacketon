import { createRequire } from 'node:module'
import { fileURLToPath } from 'node:url'
import path from 'node:path'
import fs from 'node:fs/promises'
import { spawnSync } from 'node:child_process'
import { synthetic, recordingText } from './recording-english.mjs'

// Real browser interactions with the real local API. No mocked AI responses.
// Run from the repository root: node scripts/demo/record-simulation.mjs [--english] [--apertus] [--natural]
const require = createRequire(new URL('../../src/laporsiap_FE/package.json', import.meta.url))
const repo = fileURLToPath(new URL('../../', import.meta.url))
process.env.PLAYWRIGHT_BROWSERS_PATH ||= path.join(repo, 'tmp/playwright')
const { chromium, expect } = require('@playwright/test')
const runDate = new Intl.DateTimeFormat('en-CA', { timeZone: 'Asia/Jakarta' }).format(new Date())
const liveAI = process.argv.includes('--apertus')
const naturalCase = process.argv.includes('--natural')
const english = liveAI || process.argv.includes('--english')
const language = english ? 'en' : 'id'
const ui = text => { if (!english) return text; try { return recordingText(text) } catch (error) { if (liveAI) return text; throw error } }
const out = path.join(repo, 'output/simulation', runDate, ...(liveAI ? ['apertus'] : []), ...(english ? ['english'] : []), ...(naturalCase ? ['natural-case'] : []))
const raw = path.join(repo, 'tmp/simulation', runDate, ...(liveAI ? ['apertus'] : []), ...(english ? ['english'] : []), ...(naturalCase ? ['natural-case'] : []))
const englishFiles = { 'bukti-sintetis.pdf': 'transfer-reference.pdf', 'LaporSiap-simulasi.zip': 'LaporSiap-simulation.zip', 'LaporSiap-kronologi.txt': 'LaporSiap-chronology.txt', '01-contoh-sintetis.png': '01-synthetic-case.png', '02-draft-manual.png': '02-manual-draft.png', '03-review-transfer.png': '03-transfer-review.png', '04-ekspor-siap.png': '04-export-ready.png', '05-simulasi-selesai.png': '05-walkthrough-complete.png' }
const asset = name => path.join(out, english ? (englishFiles[name] || name) : name)
await fs.mkdir(out, { recursive: true })
await fs.mkdir(raw, { recursive: true })
const started = new Date().toISOString()
let browser, context, page, video, latestCase, problem
let stepName = 'Menyiapkan browser'
let liveAICalled = false
let generatedDraft
let providerDiagnostic
let apiErrorPending = false
let monitoring = false
let rejectFailure
const failure = new Promise((_, reject) => { rejectFailure = reject })
failure.catch(() => {})
const steps = []
function fail(reason) {
  if (!monitoring || problem) return
  problem = new Error(reason)
  rejectFailure(problem)
}
async function caption(title, message) {
  title = ui(title); message = ui(message)
  await page.evaluate(({ title, message }) => {
    let box = document.getElementById('simulation-caption')
    if (!box) {
      box = document.createElement('div')
      box.id = 'simulation-caption'
      box.style.cssText = 'position:fixed;bottom:18px;left:340px;right:28px;z-index:2147483646;padding:15px 24px;border:1px solid #6a9585;border-left:5px solid #f0ad4e;border-radius:12px;background:#092e2bf5;color:#f7fcf8;pointer-events:none;font-family:"DM Sans",sans-serif;box-shadow:0 8px 30px #0004;'
      document.body.append(box)
    }
    box.replaceChildren()
    const heading = document.createElement('div')
    heading.textContent = title
    heading.style.cssText = 'color:#f0ad4e;font-size:17px;font-weight:700;margin-bottom:4px;letter-spacing:.6px;'
    const body = document.createElement('div')
    body.textContent = message
    body.style.cssText = 'font-size:23px;line-height:1.35;'
    box.append(heading, body)
    if (!document.getElementById('simulation-cursor')) {
      const cursor = document.createElement('div')
      cursor.id = 'simulation-cursor'
      cursor.style.cssText = 'position:fixed;left:-30px;top:-30px;width:22px;height:22px;border:2px solid #b94f26;background:#f0ad4e66;border-radius:50%;pointer-events:none;z-index:2147483647;transform:translate(-50%,-50%);box-shadow:0 0 0 5px #f0ad4e22;'
      document.body.append(cursor)
      document.addEventListener('pointermove', e => { cursor.style.left = `${e.clientX}px`; cursor.style.top = `${e.clientY}px` })
    }
  }, { title, message })
}
const pause = (seconds = 2.4) => page.waitForTimeout(seconds * 1000)
async function focus(locator) {
  await locator.evaluate(el => el.scrollIntoView({ block: 'center', behavior: 'smooth' }))
  await pause(.55)
  const rect = await locator.boundingBox()
  if (rect) await page.mouse.move(rect.x + rect.width / 2, rect.y + rect.height / 2, { steps: 14 })
  await locator.evaluate(el => {
    const old = el.style.outline
    el.style.outline = '3px solid #f0ad4e'
    setTimeout(() => { el.style.outline = old }, 1200)
  })
  await pause(.4)
}
async function click(locator) { await focus(locator); await locator.click() }
async function step(name, fn) {
  name = ui(name)
  stepName = name
  console.log(`START ${name}`)
  const entry = { name, status: 'running', started_at: new Date().toISOString() }
  steps.push(entry)
  try {
    await Promise.race([fn(), failure])
    if (problem) throw problem
    entry.status = 'passed'; entry.finished_at = new Date().toISOString()
    console.log(`PASS ${name}`)
  } catch (error) {
    entry.status = 'failed'; entry.error = error.message; throw error
  }
}
function referencePDF() {
  const lines = english ? synthetic.pdfLines : [
    [22, 'LaporSiap - DATA REKAAN / SIMULASI'],
    [12, 'Lampiran referensi sintetis. Bukan bukti transaksi bank.'],
    [12, 'Bank Aruna / ****7392'],
    [12, 'Tujuan: Bank Lestari / ****4821'],
    [12, 'Catatan pengguna: Rp180.000 berstatus berhasil.'],
    [12, 'Catatan lain: Rp420.000 masih dalam proses.'],
    [12, 'Tanggal transfer belum diketahui; waktu chat bukan waktu transfer.'],
    [12, 'Seluruh identitas dan rekening di atas adalah fiktif.'],
  ]
  const commands = lines.map(([size, text], i) => `BT /F1 ${size} Tf 45 ${760 - i * 40} Td (${text}) Tj ET`).join('\n')
  const objects = [
    '<< /Type /Catalog /Pages 2 0 R >>',
    '<< /Type /Pages /Kids [3 0 R] /Count 1 >>',
    '<< /Type /Page /Parent 2 0 R /MediaBox [0 0 595 842] /Resources << /Font << /F1 4 0 R >> >> /Contents 5 0 R >>',
    '<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>',
    `<< /Length ${Buffer.byteLength(commands)} >>\nstream\n${commands}\nendstream`,
  ]
  let pdf = '%PDF-1.4\n'
  const offsets = [0]
  objects.forEach((object, i) => { offsets.push(Buffer.byteLength(pdf)); pdf += `${i + 1} 0 obj\n${object}\nendobj\n` })
  const xref = Buffer.byteLength(pdf)
  pdf += `xref\n0 ${objects.length + 1}\n0000000000 65535 f \n`
  pdf += offsets.slice(1).map(offset => `${String(offset).padStart(10, '0')} 00000 n \n`).join('')
  pdf += `trailer\n<< /Size ${objects.length + 1} /Root 1 0 R >>\nstartxref\n${xref}\n%%EOF\n`
  return Buffer.from(pdf)
}
let result = 'Failed'
try {
  browser = await chromium.launch({ channel: 'chrome', headless: true, slowMo: 180 })
  context = await browser.newContext({ viewport: { width: 1920, height: 1080 },
    recordVideo: { dir: raw, size: { width: 1920, height: 1080 } }, acceptDownloads: true })
  await context.addInitScript(lang => { try { localStorage.setItem('laporsiap-language', lang) } catch {} }, language)
  page = await context.newPage(); video = page.video()
  page.setDefaultTimeout(12000)
  monitoring = true
  page.on('pageerror', error => fail(`Kesalahan aplikasi: ${error.message}`))
  page.on('requestfailed', request => {
    if (new URL(request.url()).pathname.startsWith('/api/')) fail(`Permintaan API terputus: ${request.failure()?.errorText}`)
  })
  page.on('request', request => { if (request.method() === 'POST' && new URL(request.url()).pathname.endsWith('/draft') && request.postDataJSON()?.mode === 'apertus') liveAICalled = true })
  page.on('response', async response => {
    if (!new URL(response.url()).pathname.startsWith('/api/')) return
    if (!response.ok()) { apiErrorPending = true; const error = await response.json().catch(() => ({})); providerDiagnostic = error.provider_diagnostic; fail(`API memberi HTTP ${response.status()}: ${error.error || 'request failed'}`); return }
    try { const data = await response.json(); if (data.case) latestCase = data.case } catch {}
  })
  await step('01 — Buka aplikasi dan muat contoh', async () => {
    await page.goto('http://127.0.0.1:5175/', { waitUntil: 'networkidle' })
    await page.evaluate(() => {
      new MutationObserver(() => {
        const error = document.querySelector('.alert.error')
        if (error) window.dispatchEvent(new CustomEvent('simulation-app-error', { detail: error.textContent }))
      }).observe(document.body, { subtree: true, childList: true, characterData: true })
    })
    await page.exposeFunction('simulationError', message => { if (!apiErrorPending) fail(`Pesan error aplikasi: ${message}`) })
    await page.evaluate(() => window.addEventListener('simulation-app-error', e => window.simulationError(e.detail)))
    await caption(liveAI ? 'LAPORSIAP • SYNTHETIC DATA • LIVE APERTUS' : 'SIMULASI LAPORSIAP • DATA SINTETIS • DRAFT MANUAL', 'Mulai dari contoh rekaan yang tersedia di aplikasi.')
    await pause(3)
    if (english) {
      await expect(page.locator('html')).toHaveAttribute('lang', 'en')
      await expect(page.getByRole('combobox', { name: 'Language / Bahasa' })).toHaveValue('en')
      await page.getByRole('textbox', { name: 'Case title', exact: true }).fill(synthetic.title)
      await page.getByLabel('Conversation format').selectOption('whatsapp_bracket_dmy')
      const conversationFile = asset('synthetic-conversation.txt')
      await fs.writeFile(conversationFile, synthetic.rawChat + '\n')
      await focus(page.locator('label.file-button'))
      await page.locator('input[type=file][accept=".txt,text/plain"]').setInputFiles(conversationFile)
      await expect(page.getByRole('textbox', { name: 'Conversation', exact: true })).toHaveValue(synthetic.rawChat + '\n')
      await page.getByRole('checkbox', { name: 'This is simulated data', exact: true }).check()
      await caption('01 / BAHAN KASUS', 'Rp180.000 berhasil; Rp420.000 masih pending. Tanggal transfer tetap kosong.')
      for (const [index, transfer] of synthetic.transfers.entries()) {
        await click(page.getByRole('button', { name: '+ Add a transfer record', exact: true }))
        for (const [label, key] of [['Sending bank','fromBank'],['Sending account','fromAccount'],['Receiving bank','toBank'],['Receiving account','toAccount'],['Amount (rupiah)','amount'],['Note','note']]) {
          await page.getByLabel(label, { exact: true }).nth(index).fill(transfer[key])
        }
        // Nested select options are included in the label text; match its label prefix.
        await page.getByLabel('Transfer status').nth(index).selectOption(transfer.status)
      }
      await fs.writeFile(asset('synthetic-case.json'), JSON.stringify(synthetic, null, 2))
    } else {
      await click(page.getByRole('button', { name: 'Coba contoh', exact: true }))
    }
    await expect(page.getByRole('checkbox', { name: ui('Ini data simulasi') })).toBeChecked()
    // Input values, not element textContent, are the source of truth for textarea.
    await expect(page.getByRole('textbox', { name: ui('Percakapan'), exact: true })).toHaveValue(/Nadia/)
    if (naturalCase) {
      await expect(page.locator('.simulation-banner')).toContainText(english ? 'SIMULATED DATA' : 'DATA SIMULASI')
      await expect(page.locator('body')).not.toContainText(/\bdemo\b/i)
    }
    await pause(3)
    await page.screenshot({ path: asset('01-contoh-sintetis.png'), fullPage: true })
  })
  await step('02 — Periksa status dan simpan bahan', async () => {
    await caption('01 / BAHAN KASUS', 'Rp180.000 berhasil; Rp420.000 masih pending. Tanggal transfer tetap kosong.')
    await focus(page.getByRole('heading', { name: ui('Catatan transfer'), exact: true }))
    await pause(3)
    const amounts = page.locator('input[type="number"]')
    await expect(amounts.nth(0)).toHaveValue('180000'); await expect(amounts.nth(1)).toHaveValue('420000')
    await click(page.getByRole('button', { name: ui('Simpan & lanjut'), exact: true }))
    await expect(page.getByRole('heading', { name: ui('Susun dengan Apertus'), exact: true })).toBeVisible()
    await expect(page.locator('.summary-card .total')).toHaveText(english ? /180,000/ : /180.000/)
    await pause(2)
  })
  await step('03 — Tambahkan lampiran sintetis', async () => {
    await caption('01 / LAMPIRAN REFERENSI', 'Unggah PDF rekaan. Lampiran disimpan sebagai referensi, tanpa OCR.')
    await click(page.getByRole('navigation').getByRole('button', { name: english ? /Case materials/ : /Bahan kasus/ }))
    const reference = asset('bukti-sintetis.pdf')
    await fs.writeFile(reference, referencePDF())
    const upload = page.locator('input[type=file][accept="image/png,image/jpeg,application/pdf"]')
    await focus(page.getByText(ui('Pilih file bukti'), { exact: true }))
    await upload.setInputFiles(reference)
    await expect(page.getByText(english ? 'transfer-reference.pdf' : 'bukti-sintetis.pdf', { exact: true })).toBeVisible()
    if (liveAI) {
      await caption('01 / LINK EVIDENCE', 'Link the fictional attachment to source M3 and transfer 01. This is a user-selected reference, not verification.')
      await click(page.getByRole('button', { name: 'Link evidence', exact: true }))
      const editor = page.getByRole('region', { name: 'Link attachment' })
      await click(editor.getByRole('checkbox', { name: /^M3 / }))
      await click(editor.getByRole('checkbox', { name: /^Transfer 01 / }))
      await click(editor.getByRole('button', { name: 'Save evidence links', exact: true }))
      await expect(page.locator('.attachment-row')).toContainText('M3 · Transfer 01')
    }
    await pause(3)
    await click(page.getByRole('navigation').getByRole('button', { name: english ? /Create draft/ : /Susun draft/ }))
  })
  await step(liveAI ? '04 — Live Apertus draft' : '04 — Susun draft manual', async () => {
    if (liveAI) {
      await caption('02 / LIVE APERTUS', 'Review the source text and consent to one real Apertus request. No automatic retries.')
      await click(page.getByText('View all messages to be sent', { exact: true }))
      await pause(2)
      await click(page.getByRole('checkbox', { name: /I consent to sending/ }))
      await click(page.getByRole('button', { name: 'Create draft with Apertus', exact: true }))
      await expect(page.getByText('APERTUS DRAFT', { exact: true })).toBeVisible({ timeout: 55000 })
      generatedDraft = structuredClone(latestCase.draft)
      if (generatedDraft.model !== 'swiss-ai/apertus-v1.5-70b' || generatedDraft.language !== 'en' || generatedDraft.events.length !== 5) throw new Error('Live Apertus identity, language, or event count invalid')
      await fs.writeFile(asset('apertus-original-draft.json'), JSON.stringify({ synthetic: true, draft: generatedDraft, messages: latestCase.messages, note: 'One live application request. Model identity is API metadata. Automated scenario assertions and assistant review are not an independent user evaluation.' }, null, 2))
      await page.screenshot({ path: asset('02-apertus-draft.png'), fullPage: true })
      await pause(4)
    } else {
      await caption('02 / SUSUN DRAFT MANUAL', 'Pada rekaman ini, draft dimulai dari kutipan pesan asli. Apertus tidak dipanggil.')
      await expect(page.getByRole('button', { name: ui('Susun draft dengan Apertus'), exact: true })).toBeDisabled()
      await pause(3)
      await click(page.getByRole('button', { name: ui('Gunakan draft manual'), exact: true }))
      await expect(page.getByText(ui('DRAFT MANUAL'), { exact: true })).toBeVisible()
      if (english) await expect(page.getByRole('textbox', { name: 'Chronology M1', exact: true })).toHaveValue(/^Message quote: Hi Raka, we offer tasks/)
      await page.screenshot({ path: asset('02-draft-manual.png'), fullPage: true })
      await pause(3)
    }
  })
  const descriptions = liveAI ? generatedDraft.events.map(e => e.description) : english ? synthetic.descriptions : [
    'Nadia menawarkan tugas memberi like dengan janji komisi Rp50.000 per hari.',
    'Nadia meminta deposit Rp180.000 ke ****4821 untuk melanjutkan tugas.',
    'Raka menyatakan baru mengirim Rp180.000 dan belum mencatat waktu transfer.',
    'Nadia meminta tambahan deposit Rp420.000 agar saldo bisa ditarik.',
    'Raka menyatakan setoran tambahan masih dalam proses dan belum melihat mutasi rekening.',
  ]
  await step('05 — Cocokkan dan sunting kelima entri', async () => {
    if (liveAI) {
      const conditions = [/promis|offer/i, /ask|request|instruct|demand/i, /sent|send|transfer/i, /ask|request|instruct|demand/i, /pending|in progress/i]
      for (let i = 0; i < descriptions.length; i++) {
        const e = generatedDraft.events[i]
        if (e.source_id !== latestCase.messages[i].id || e.evidence !== latestCase.messages[i].text) throw new Error(`Source quote changed at M${i+1}`)
        if (!conditions[i].test(e.description) || /pengirim|menyatakan|meminta|menjanjikan/i.test(e.description)) throw new Error(`Scenario meaning or English language needs review at M${i+1}`)
        await caption(`03 / REVIEW LIVE APERTUS M${i+1}`, 'Compare the generated summary with the unchanged source. Review remains required.')
        await focus(page.getByRole('textbox', { name: `Chronology M${i+1}`, exact: true }))
        await pause(2)
      }
      return
    }
    for (let i = 0; i < descriptions.length; i++) {
      await caption(`03 / PERIKSA SUMBER M${i + 1}`, 'Baca pesan asli di kanan, lalu rapikan kalimat tanpa menambah fakta.')
      const input = page.getByRole('textbox', { name: `${english ? 'Chronology' : 'Kronologi'} M${i + 1}`, exact: true })
      await focus(input)
      await pause(1.8)
      await input.fill(descriptions[i])
      await pause(1.4)
    }
    await caption('03 / SIMPAN KOREKSI', 'Koreksi disimpan sebagai revisi. Selanjutnya periksa semua sumber kembali.')
    await click(page.getByRole('button', { name: ui('Simpan koreksi'), exact: true }))
    await expect(page.getByText(ui('Koreksi tersimpan.'), { exact: false })).toBeVisible()
    await pause(2)
  })
  await step('06 — Review sumber, transfer, dan data kosong', async () => {
    const checks = page.getByRole('checkbox', { name: ui('Saya telah mencocokkan entri ini dengan sumber'), exact: true })
    await expect(checks).toHaveCount(5)
    for (let i = 0; i < 5; i++) {
      const source = page.locator('.event-card').nth(i).locator('blockquote')
      const original = await source.textContent()
      const draft = await page.getByRole('textbox', { name: `${english ? 'Chronology' : 'Kronologi'} M${i + 1}`, exact: true }).inputValue()
      if (!original || draft !== descriptions[i]) throw new Error(`Sumber atau koreksi M${i + 1} tidak sesuai`)
      await caption(`03 / REVIEW M${i + 1} DARI 5`, 'Tandai entri setelah membandingkan draft dengan pesan sumbernya.')
      await click(checks.nth(i)); await pause(1.2)
    }
    await caption('03 / PERIKSA CATATAN TRANSFER', 'Total Rp180.000 berasal dari status berhasil. Rp420.000 pending tidak masuk total.')
    await focus(page.getByRole('heading', { name: ui('Periksa catatan transfer'), exact: true }))
    await pause(3)
    await click(page.getByRole('checkbox', { name: english ? /I have checked the accounts/ : /Saya telah memeriksa rekening/ }))
    await caption('03 / DATA BELUM LENGKAP', 'Tanggal transfer belum diketahui. Ekspor adalah bahan bantu untuk kanal resmi.')
    await click(page.getByRole('checkbox', { name: english ? /I have checked the review flags/ : /Saya telah memeriksa tanda review/ }))
    await pause(3)
    await page.screenshot({ path: asset('03-review-transfer.png'), fullPage: true })
    await click(page.getByRole('button', { name: ui('Konfirmasi & lanjut'), exact: true }))
    await expect(page.getByRole('button', { name: ui('Unduh paket ZIP'), exact: true })).toBeEnabled()
    await pause(3)
  })
  await step('07 — Unduh ZIP dan kronologi TXT', async () => {
    await caption('04 / EKSPOR BAHAN', 'Unduh paket ZIP dan kronologi TXT dari revisi yang telah dikonfirmasi.')
    await page.screenshot({ path: asset('04-ekspor-siap.png'), fullPage: true })
    for (const [button, filename] of [['Unduh paket ZIP', 'LaporSiap-simulasi.zip'], ['Unduh kronologi TXT', 'LaporSiap-kronologi.txt'], ...(liveAI ? [['Download PDF report', 'LaporSiap-report.pdf']] : [])]) {
      const downloadPromise = page.waitForEvent('download')
      downloadPromise.catch(() => {})
      await click(page.getByRole('button', { name: ui(button), exact: true }))
      const download = await downloadPromise
      await download.saveAs(asset(filename))
      const error = await download.failure(); if (error) throw new Error(`Unduhan gagal: ${error}`)
      await pause(2)
    }
    const check = spawnSync('python3', ['-c', `
import json, sys
from zipfile import ZipFile
with ZipFile(sys.argv[1]) as z:
    assert z.testzip() is None, 'ZIP corrupt'
    assert {'kronologi.txt','percakapan-asli.txt','kasus-direview.json','transfer.csv'}.issubset(z.namelist())
    c=json.loads(z.read('kasus-direview.json'))
    assert c['synthetic'] is True and c['reviewed_revision'] == c['revision']
    assert len(c['messages']) == 5 and len(c['draft']['events']) == 5
    assert c['draft']['origin'] == sys.argv[4]
    assert 'laporan.pdf' in z.namelist()
    assert [(t['amount'], t['status'], t['occurred_at']) for t in c['transfers']] == [(180000, 'successful', ''), (420000, 'pending', '')]
    assert any(n.endswith(sys.argv[2]) for n in z.namelist())
    if sys.argv[3] == 'en':
        assert c['title'].startswith('Online task offer')
        assert all(m['sender'] in ['Nadia', 'Raka'] for m in c['messages'])
        assert 'I have just sent IDR 180,000.' in c['raw_chat']
        assert 'LAPORSIAP — SUPPORTING REPORT MATERIALS' in z.read('kronologi.txt').decode()
        assert 'Kutipan pesan:' not in z.read('kronologi.txt').decode()
        assert 'SIMULATED DATA' in z.read('kronologi.txt').decode()
    print('ZIP validated: synthetic, 5 reviewed sources, correct transfer statuses, reference PDF')
`, asset('LaporSiap-simulasi.zip'), english ? 'transfer-reference.pdf' : 'bukti-sintetis.pdf', language, liveAI ? 'apertus' : 'manual'], { encoding: 'utf8' })
    if (check.status !== 0) throw new Error(`Validasi ZIP gagal: ${check.stderr.trim()}`)
    console.log(check.stdout.trim())
    const txt = await fs.readFile(asset('LaporSiap-kronologi.txt'), 'utf8')
    if (!txt.includes(english ? '180,000' : '180.000') || !txt.includes(english ? '420,000' : '420.000')) throw new Error('Rincian nominal pada TXT tidak ditemukan')
  })
  await step('08 — Periksa riwayat dan hasil akhir', async () => {
    await caption(liveAI ? 'LIVE APERTUS WALKTHROUGH COMPLETE' : 'SIMULASI SELESAI', liveAI ? 'Synthetic input → real Apertus draft → source review → ZIP, TXT and PDF. Model identity is API metadata.' : 'Input → lampiran → draft manual → koreksi → review → ZIP dan TXT. Semua berhasil.')
    await click(page.getByRole('button', { name: ui('Lihat riwayat revisi →'), exact: true }))
    await expect(page.locator('.history details')).toHaveCount(5)
    await pause(3)
    await page.screenshot({ path: asset('05-simulasi-selesai.png'), fullPage: true })
  })
  result = 'Success'
} catch (error) {
  problem ||= error
  console.error(`STOP at ${stepName}: ${problem.message}`)
  // No retries or fallback during a failed recording: preserve the failed attempt.
  monitoring = false
  if (context) await context.close().catch(() => {})
  context = undefined
} finally {
  monitoring = false
  if (context) await context.close().catch(() => {}) // Finalizes and stops recording immediately.
  if (video) {
    let number = 1
    while (await fs.stat(path.join(out, `LaporSiap_${result}_${number}.webm`)).then(() => true).catch(() => false)) number++
    const basename = `LaporSiap_${result}_${number}`
    await video.saveAs(path.join(out, `${basename}.webm`))
    const converted = spawnSync('ffmpeg', ['-hide_banner', '-loglevel', 'warning', '-n',
      '-i', path.join(out, `${basename}.webm`), '-vf', 'fps=30', '-c:v', 'libx264',
      '-preset', 'veryfast', '-crf', '18', '-pix_fmt', 'yuv420p', '-movflags', '+faststart',
      '-an', path.join(out, `${basename}.mp4`)], { encoding: 'utf8' })
    if (converted.status !== 0) console.error(`MP4 conversion failed; original WebM preserved: ${converted.stderr || converted.error?.message}`)
    const report = { app: 'LaporSiap', result, mode: liveAI ? 'apertus' : 'manual', language, synthetic_data_language: language, synthetic: true, scenario: naturalCase ? 'natural-fictional-case-v1' : 'fictional-case', live_ai_called: liveAICalled, draft_metadata: generatedDraft ? { model: generatedDraft.model, language: generatedDraft.language, generated_at: generatedDraft.generated_at, duration_ms: generatedDraft.duration_ms, usage: generatedDraft.usage, prompt_version: generatedDraft.prompt_version } : null,
      provider_diagnostic: providerDiagnostic || null,
      started_at: started, finished_at: new Date().toISOString(), recording: `${basename}.${converted.status === 0 ? 'mp4' : 'webm'}`,
      stopped_at: result === 'Failed' ? stepName : null, error: problem?.message || null,
      case_id: latestCase?.id || null, revision: latestCase?.revision || null, steps }
    await fs.writeFile(path.join(out, `${basename}.json`), JSON.stringify(report, null, 2))
    console.log(`RECORDING ${path.join(out, `${basename}.webm`)}`)
  }
  if (browser) await browser.close().catch(() => {})
}
process.exitCode = result === 'Success' ? 0 : 1
