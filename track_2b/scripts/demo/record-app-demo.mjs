// Human-paced screen recording of the LaporSiap app for the submission video.
// Real browser, real local API, one live Apertus request (no mocks, no retries).
// Smooth eased cursor, gentle scrolling and pauses; captured as high-quality
// Chrome screencast frames, then encoded to a constant 30 fps MP4 without audio.
// Usage (from track_2b): node scripts/demo/record-app-demo.mjs [--base http://127.0.0.1:8787]
import { createRequire } from 'node:module'
import { fileURLToPath } from 'node:url'
import path from 'node:path'
import fs from 'node:fs/promises'
import { spawnSync } from 'node:child_process'

const require = createRequire(new URL('../../src/laporsiap_FE/package.json', import.meta.url))
const repo = fileURLToPath(new URL('../../', import.meta.url))
const { chromium } = require('@playwright/test')
const base = process.argv.includes('--base') ? process.argv[process.argv.indexOf('--base') + 1] : 'http://127.0.0.1:8787'
const rehearsal = process.argv.includes('--manual') // no AI call; for rehearsing the motion
const stamp = new Date().toISOString().replace(/[:.]/g, '-')
const out = path.join(repo, 'output/demo', stamp + (rehearsal ? '-rehearsal' : ''))
const framesDir = path.join(repo, 'tmp/demo-frames', stamp)
await fs.mkdir(out, { recursive: true })
await fs.mkdir(framesDir, { recursive: true })

const VIEW = { width: 1600, height: 900 }, SCALE = 1.2 // → 1920×1080 frames
const frames = []
let t0 = 0
const marks = []
const mark = name => { marks.push({ scene: name, at_seconds: +((Date.now() / 1000) - t0).toFixed(2) }); console.log('MARK', name) }

// ---------- capture ----------
let cdp
async function capture(page) {
  if (cdp) await cdp.send('Page.stopScreencast').catch(() => {})
  cdp = await page.context().newCDPSession(page)
  cdp.on('Page.screencastFrame', async ({ data, metadata, sessionId }) => {
    const file = path.join(framesDir, `${String(frames.length).padStart(6, '0')}.jpg`)
    frames.push({ file, ts: metadata.timestamp })
    fs.writeFile(file, Buffer.from(data, 'base64')).catch(() => {})
    cdp.send('Page.screencastFrameAck', { sessionId }).catch(() => {})
  })
  await cdp.send('Page.startScreencast', { format: 'jpeg', quality: 92, maxWidth: 1920, maxHeight: 1080, everyNthFrame: 1 })
}

// ---------- human-like helpers ----------
const sleep = ms => new Promise(r => setTimeout(r, ms))
const ease = t => t < .5 ? 4 * t * t * t : 1 - Math.pow(-2 * t + 2, 3) / 2
let page, cursor = { x: 800, y: 450 }
async function moveTo(x, y) {
  const dx = x - cursor.x, dy = y - cursor.y, dist = Math.hypot(dx, dy)
  const ms = Math.min(900, 280 + dist * 0.55), steps = Math.max(12, Math.round(ms / 16))
  const bend = (Math.random() - .5) * Math.min(80, dist * .18) // slight natural arc
  for (let i = 1; i <= steps; i++) {
    const t = ease(i / steps), arc = Math.sin(Math.PI * t) * bend
    const nx = cursor.x + dx * t - (dy / (dist || 1)) * arc, ny = cursor.y + dy * t + (dx / (dist || 1)) * arc
    await page.mouse.move(nx, ny); await sleep(ms / steps)
  }
  cursor = { x, y }
}
async function scrollToEl(locator, where = .38, ms = 900) {
  await locator.evaluate(async (el, { where, ms }) => {
    const start = window.scrollY
    const target = Math.max(0, el.getBoundingClientRect().top + start - window.innerHeight * where)
    const t0 = performance.now()
    await new Promise(done => {
      const tick = now => {
        const t = Math.min(1, (now - t0) / ms), e = t < .5 ? 4 * t * t * t : 1 - Math.pow(-2 * t + 2, 3) / 2
        window.scrollTo(0, start + (target - start) * e)
        t < 1 ? requestAnimationFrame(tick) : done()
      }
      requestAnimationFrame(tick)
    })
  }, { where, ms })
  await sleep(250)
}
async function scrollBy(px, ms = 1200) {
  await page.evaluate(async ({ px, ms }) => {
    const start = window.scrollY, t0 = performance.now()
    await new Promise(done => {
      const tick = now => {
        const t = Math.min(1, (now - t0) / ms), e = t < .5 ? 4 * t * t * t : 1 - Math.pow(-2 * t + 2, 3) / 2
        window.scrollTo(0, start + px * e); t < 1 ? requestAnimationFrame(tick) : done()
      }
      requestAnimationFrame(tick)
    })
  }, { px, ms })
  await sleep(200)
}
async function hover(locator, linger = 700, scroll = true) {
  if (scroll) await scrollToEl(locator)
  const box = await locator.boundingBox()
  if (!box) throw new Error('Element not visible')
  await moveTo(box.x + box.width * (.42 + Math.random() * .16), box.y + box.height * (.4 + Math.random() * .2))
  await sleep(linger)
}
async function click(locator, linger = 260) {
  await hover(locator, linger)
  await page.mouse.down(); await sleep(90); await page.mouse.up()
  await sleep(350)
}

// Visible macOS-style cursor with a press ripple (the headless browser has no native cursor).
const cursorScript = () => {
  const install = () => {
    if (document.getElementById('demo-cursor')) return
    const c = document.createElement('div'); c.id = 'demo-cursor'
    c.innerHTML = '<svg width="26" height="34" viewBox="0 0 26 34"><path d="M2 2 L2 27 L8.5 21 L13 31.5 L17.5 29.5 L13 19.5 L22 19.5 Z" fill="#111" stroke="#fff" stroke-width="2" stroke-linejoin="round"/></svg>'
    c.style.cssText = 'position:fixed;left:0;top:0;z-index:2147483647;pointer-events:none;transform:translate(800px,450px);filter:drop-shadow(0 2px 3px rgba(0,0,0,.35));'
    const r = document.createElement('div'); r.id = 'demo-ripple'
    r.style.cssText = 'position:fixed;left:0;top:0;width:34px;height:34px;margin:-17px 0 0 -17px;border-radius:50%;background:rgba(122,31,43,.28);z-index:2147483646;pointer-events:none;opacity:0;transform:scale(.4);transition:transform .35s ease,opacity .45s ease;'
    document.body.append(r, c)
    addEventListener('mousemove', e => { c.style.transform = `translate(${e.clientX - 2}px,${e.clientY - 2}px)` }, true)
    addEventListener('mousedown', e => {
      r.style.transition = 'none'; r.style.left = e.clientX + 'px'; r.style.top = e.clientY + 'px'
      r.style.opacity = '1'; r.style.transform = 'scale(.4)'
      requestAnimationFrame(() => { r.style.transition = 'transform .35s ease,opacity .45s ease'; r.style.transform = 'scale(1.3)'; r.style.opacity = '0' })
    }, true)
  }
  document.readyState === 'loading' ? addEventListener('DOMContentLoaded', install) : install()
}

// ---------- recording ----------
const report = { app: 'LaporSiap', synthetic: true, base, started_at: new Date().toISOString(), live_ai_requests: 0 }
const browser = await chromium.launch({ channel: 'chrome', headless: true })
let ok = false
try {
  const context = await browser.newContext({ viewport: VIEW, deviceScaleFactor: SCALE, acceptDownloads: true })
  await context.addInitScript(() => { try { localStorage.setItem('laporsiap-language', 'en') } catch {} })
  await context.addInitScript(cursorScript)
  page = await context.newPage()
  page.setDefaultTimeout(15000)
  let latestCase
  page.on('request', r => { if (r.method() === 'POST' && r.url().endsWith('/draft') && r.postDataJSON()?.mode === 'apertus') report.live_ai_requests++ })
  page.on('response', async r => { if (r.url().includes('/api/')) { try { const d = await r.json(); if (d.case) latestCase = d.case } catch {} } })
  await page.goto(base + '/', { waitUntil: 'networkidle' })
  await page.mouse.move(cursor.x, cursor.y)
  await capture(page)
  t0 = Date.now() / 1000
  await sleep(1500)

  // 1 · Case materials
  mark('step1_input')
  await hover(page.getByRole('link', { name: /Open IASC/ }), 1400, false)
  await click(page.getByRole('button', { name: 'Try an example', exact: true }), 400)
  await sleep(900)
  await hover(page.getByRole('textbox', { name: 'Conversation', exact: true }), 900)
  await scrollBy(520, 1600); await sleep(900)
  await scrollBy(560, 1600); await sleep(1000)
  await click(page.getByRole('button', { name: 'Save & continue', exact: true }))
  await page.getByRole('button', { name: 'Create draft with Apertus', exact: true }).waitFor()
  await sleep(1000)

  // 2 · Privacy preview + live Apertus
  mark('step2_preview_and_apertus')
  await scrollBy(-2000, 700)
  await click(page.getByText('View all messages to be sent', { exact: true }), 400)
  await sleep(700)
  await scrollBy(320, 1200)
  const masked = page.locator('.preview-messages').first()
  if (await masked.count()) await hover(masked, 1600)
  await click(page.getByRole('checkbox', { name: /I consent to sending/ }), 500)
  await sleep(500)
  if (rehearsal) {
    await click(page.getByRole('button', { name: 'Use a manual draft', exact: true }), 500)
    await page.getByText('MANUAL DRAFT', { exact: true }).waitFor()
  } else {
    await click(page.getByRole('button', { name: 'Create draft with Apertus', exact: true }), 500)
    mark('apertus_requested')
    await page.getByText('APERTUS DRAFT', { exact: true }).waitFor({ timeout: 60000 })
  }
  mark('apertus_draft_ready')
  report.draft = latestCase?.draft ? { model: latestCase.draft.model, language: latestCase.draft.language, duration_ms: latestCase.draft.duration_ms, usage: latestCase.draft.usage, prompt_version: latestCase.draft.prompt_version, events: latestCase.draft.events.length } : null
  await sleep(1200)

  // 3 · Human review, entry by entry
  mark('step3_review')
  const firstDraft = page.getByRole('textbox', { name: 'Chronology M1', exact: true })
  await hover(firstDraft, 1200)
  const firstSource = page.locator('.event-card').first().locator('blockquote')
  await hover(firstSource, 1300, false)
  const flag = page.locator('.review-flag').first()
  if (await flag.count()) await hover(flag, 1500)
  const checks = page.getByRole('checkbox', { name: 'I have compared this entry with its source', exact: true })
  const n = await checks.count()
  for (let i = 0; i < n; i++) {
    const card = page.locator('.event-card').nth(i)
    if (i > 0) await hover(card.locator('blockquote'), 700)
    await click(checks.nth(i), 300)
  }
  await click(page.getByRole('checkbox', { name: /I have checked the accounts/ }), 500)
  await click(page.getByRole('checkbox', { name: /I have checked the review flags/ }), 500)
  await click(page.getByRole('button', { name: 'Confirm & continue', exact: true }), 500)
  await page.getByRole('button', { name: 'Download PDF report', exact: true }).waitFor()
  await sleep(1000)

  // 4 · Export, then open the exported PDF
  mark('step4_export')
  await scrollBy(-2000, 700)
  await hover(page.getByRole('button', { name: 'Download PDF report', exact: true }), 600)
  const dl = page.waitForEvent('download')
  await click(page.getByRole('button', { name: 'Download PDF report', exact: true }), 200)
  const pdfPath = path.join(out, 'LaporSiap-report.pdf')
  await (await dl).saveAs(pdfPath)
  await sleep(900)
  await page.mouse.move(cursor.x + 1, cursor.y) // emit a frame before switching tabs
  const png = path.join(out, 'LaporSiap-report-page1.png')
  spawnSync('sips', ['-s', 'format', 'png', '-Z', '2200', pdfPath, '--out', png])
  const viewer = await context.newPage()
  const img = (await fs.readFile(png)).toString('base64')
  await viewer.setContent(`<title>LaporSiap-report.pdf</title><body style="margin:0;background:#525659;text-align:center;padding:28px 0 60px">
    <img src="data:image/png;base64,${img}" style="display:inline-block;width:980px;height:auto;box-shadow:0 4px 18px rgba(0,0,0,.45);background:#fff"></body>`)
  await viewer.evaluate(cursorScript)
  await viewer.bringToFront()
  await viewer.mouse.move(cursor.x, cursor.y)
  page = viewer
  await capture(viewer)
  mark('pdf_open')
  await sleep(1200)
  await moveTo(820, 520); await sleep(600)
  await scrollBy(700, 2600); await sleep(1600)
  await scrollBy(-700, 1800); await sleep(1400)
  mark('end')
  report.end_wall = Date.now() / 1000
  report.result = 'Success'; ok = true
} catch (error) {
  report.result = 'Failed'; report.error = error.message
  console.error('STOP', error.message)
} finally {
  if (cdp) await cdp.send('Page.stopScreencast').catch(() => {})
  await sleep(500)
  await browser.close().catch(() => {})
}

// ---------- encode: variable-timestamp frames → constant 30 fps ----------
report.finished_at = new Date().toISOString()
report.scene_marks = marks
if (frames.length > 1) {
  const list = frames.map((f, i) => {
    const next = frames[i + 1]?.ts ?? Math.max(f.ts + 0.5, report.end_wall ?? 0)
    return `file '${f.file}'\nduration ${Math.max(0.001, next - f.ts).toFixed(4)}`
  }).join('\n') + `\nfile '${frames.at(-1).file}'\n`
  const listFile = path.join(framesDir, 'frames.txt')
  await fs.writeFile(listFile, list)
  const name = ok ? 'LaporSiap-App-Demo.mp4' : 'LaporSiap-App-Demo-FAILED.mp4'
  const enc = spawnSync('ffmpeg', ['-hide_banner', '-loglevel', 'error', '-y', '-f', 'concat', '-safe', '0', '-i', listFile,
    '-vf', 'scale=1920:1080:flags=lanczos,fps=30,format=yuv420p', '-c:v', 'libx264', '-preset', 'slow', '-crf', '16',
    '-movflags', '+faststart', '-an', path.join(out, name)], { encoding: 'utf8' })
  if (enc.status !== 0) console.error('Encoding failed:', enc.stderr)
  report.video = name
  report.frames = frames.length
  report.duration_seconds = +(Math.max(frames.at(-1).ts + 0.5, report.end_wall ?? 0) - frames[0].ts).toFixed(1)
}
await fs.writeFile(path.join(out, 'demo-report.json'), JSON.stringify(report, null, 2))
console.log('OUTPUT', out)
process.exitCode = ok ? 0 : 1
