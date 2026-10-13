import { createRequire } from 'node:module'
import { fileURLToPath } from 'node:url'
import path from 'node:path'
import fs from 'node:fs/promises'
import { existsSync } from 'node:fs'
import { spawnSync } from 'node:child_process'

// Records a person operating the app. Does not enter data, consent, or review.
const repo = fileURLToPath(new URL('../../', import.meta.url))
const value = flag => {
  const index = process.argv.indexOf(flag)
  return index < 0 ? undefined : process.argv[index + 1]
}
const participant = value('--participant') || 'INTERNAL01'
const taskId = value('--task') || (participant === 'INTERNAL01' ? 'INTERNAL01_T01' : undefined)
const plan = JSON.parse(await fs.readFile(path.join(repo, 'docs/testing/pilot/plan.json'), 'utf8'))
const assignment = participant === 'INTERNAL01'
  ? plan.internal_rehearsal
  : plan.participants.find(p => p.participant_id === participant)?.tasks.find(t => t.task_id === taskId)
if (!assignment || assignment.task_id !== taskId) throw new Error('Choose INTERNAL01 or a participant/task pair from pilot/plan.json.')
const language = value('--language') || 'en'
if (!['en', 'id'].includes(language)) throw new Error('--language must be en or id')
const require = createRequire(new URL('../../src/laporsiap_FE/package.json', import.meta.url))
const localCache = path.join(repo, 'tmp/playwright')
if (existsSync(localCache)) process.env.PLAYWRIGHT_BROWSERS_PATH ||= localCache
const { chromium } = require('@playwright/test')
if (process.argv.includes('--check')) {
  // Technical recorder check only: no app, people, API, or provider calls.
  const checkOut = path.join(repo, 'tmp/user-recorder-check', new Date().toISOString().replaceAll(':', '-'))
  await fs.mkdir(checkOut, { recursive: true })
  const checkBrowser = await chromium.launch({ channel: 'chrome', headless: true })
  try {
    const checkContext = await checkBrowser.newContext({
      viewport: { width: 1440, height: 1050 },
      recordVideo: { dir: checkOut, size: { width: 1440, height: 1050 } }
    })
    await checkContext.exposeBinding('checkControl', () => 'control works')
    const checkPage = await checkContext.newPage()
    const checkVideo = checkPage.video()
    await checkPage.setContent('<p>Technical recorder check — no participant session</p>')
    if (await checkPage.evaluate(() => window.checkControl()) !== 'control works') throw new Error('Control binding failed')
    await checkPage.waitForTimeout(500)
    await checkContext.close()
    const target = path.join(checkOut, 'recorder-check.webm')
    await checkVideo.saveAs(target)
    if ((await fs.stat(target)).size === 0) throw new Error('Empty recording')
    const converted = spawnSync('ffmpeg', ['-hide_banner', '-loglevel', 'error', '-n', '-i', target, '-c:v', 'libx264', '-pix_fmt', 'yuv420p', path.join(checkOut, 'recorder-check.mp4')], { encoding: 'utf8' })
    if (converted.status !== 0) throw new Error('MP4 conversion failed')
    console.log(JSON.stringify({ technical_recorder_check: 'passed', output: checkOut, human_sessions: 0, app_or_provider_calls: 0 }))
  } finally { await checkBrowser.close() }
  process.exit(0)
}
const out = path.join(repo, 'output/user-testing/sessions', taskId, new Date().toISOString().replaceAll(':', '-'))
await fs.mkdir(path.join(out, 'inputs'), { recursive: true })
await fs.mkdir(path.join(out, 'downloads'))
for (const suffix of ['conversation.txt', 'conversation-reference.png', 'transfer-notes.jpg']) {
  const name = `${assignment.case_id}-${suffix}`
  await fs.copyFile(path.join(repo, 'output/testing/2026-10-10/inputs', name), path.join(out, 'inputs', name), fs.constants.COPYFILE_EXCL)
}
await fs.copyFile(path.join(repo, 'docs/testing/pilot/participant-brief.txt'), path.join(out, 'participant-brief.txt'), fs.constants.COPYFILE_EXCL)

let browser, context, page, video, closing = false, pauseStart = null, downloadIndex = 0
let resolveEnd
const end = new Promise(resolve => { resolveEnd = resolve })
const requests = new Map(), pendingDownloads = new Set()
const report = {
  task_id: taskId, participant_id: participant, case_id: assignment.case_id,
  planned_mode: assignment.mode,
  participant_role: participant === 'INTERNAL01' ? 'internal_rehearsal' : 'unverified_planned_independent',
  independent_participation_verified: false,
  synthetic: true, initial_ui_language: language, status: 'not_started',
  recorder_opened_at: new Date().toISOString(), start_at: null, finish_at: null,
  stopped_at: null, elapsed_seconds: null, off_task_pause_seconds: 0,
  total_task_seconds: null, ai_wait_seconds: null,
  authored_entries: null, corrected_entries: null, errors_remaining: null,
  ease_1_to_7: null, comprehension: null, assistance_count: null,
  observed_draft_requests: [], downloads: [], events: [], recording: null,
  quality_review_status: 'not_reviewed',
  limitations: [
    'Participant identity/independence and task quality require facilitator verification.',
    'Completion is the participant declaration after ZIP/TXT downloads, not semantic validation.',
    'Draft request duration measures HTTP completion, not draft visibility or provider-only latency.',
    'AI UI wait, authored/corrected entries and comprehension remain unobserved until facilitator review.',
    'Video begins when browser opens; task timer begins only on Start. No audio is captured.'
  ]
}
function finish(status, reason) {
  if (closing) return
  closing = true
  const now = Date.now()
  if (pauseStart !== null) report.off_task_pause_seconds += (now - pauseStart) / 1000
  report.status = status
  if (status === 'completed_declared') report.finish_at = new Date(now).toISOString()
  else report.stopped_at = new Date(now).toISOString()
  if (report.start_at) {
    report.elapsed_seconds = (now - Date.parse(report.start_at)) / 1000
    report.total_task_seconds = report.elapsed_seconds - report.off_task_pause_seconds
  }
  report.stop_reason = reason
  // Initiate context closure immediately on failure; no retry or fallback.
  const close = context?.close().catch(() => {})
  resolveEnd(close)
}
const operation = request => {
  const url = new URL(request.url())
  if (url.origin !== 'http://127.0.0.1:5175' || !url.pathname.startsWith('/api/')) return null
  if (/\/draft$/.test(url.pathname) && request.method() === 'POST') return 'draft'
  if (/\/export$/.test(url.pathname)) return 'export'
  return 'api'
}
try {
  browser = await chromium.launch({ channel: 'chrome', headless: false })
  context = await browser.newContext({
    viewport: { width: 1440, height: 1050 }, acceptDownloads: true,
    recordVideo: { dir: path.join(out, 'raw-video'), size: { width: 1440, height: 1050 } }
  })
  await context.addInitScript(lang => {
    try { if (!localStorage.getItem('laporsiap-language')) localStorage.setItem('laporsiap-language', lang) } catch {}
  }, language)
  await context.exposeBinding('pilotControl', async ({ page: sender }, action) => {
    if (sender !== page || closing) return 'Session ended.'
    if (action === 'start' && report.status === 'not_started') {
      report.start_at = new Date().toISOString()
      report.status = 'running'
      console.log(`START ${taskId}`)
      return 'Recording task — you operate the app. / Tugas sedang direkam.'
    }
    if (action === 'pause' && ['running', 'paused'].includes(report.status)) {
      if (pauseStart === null) { pauseStart = Date.now(); report.status = 'paused'; return 'Paused / Jeda. Click Pause again to resume.' }
      report.off_task_pause_seconds += (Date.now() - pauseStart) / 1000
      pauseStart = null; report.status = 'running'
      return 'Resumed / Dilanjutkan.'
    }
    if (action === 'complete' && report.start_at) {
      const extensions = new Set(report.downloads.filter(d => d.saved).map(d => path.extname(d.filename).toLowerCase()))
      if (!extensions.has('.zip') || !extensions.has('.txt')) return 'Download ZIP and TXT first. / Unduh ZIP dan TXT dahulu.'
      finish('completed_declared', 'Participant declares TXT opened and task finished; quality review pending.')
    } else if (action === 'failed') finish('failed', 'Participant/facilitator reported unexpected blocking failure.')
    else if (action === 'stop') finish(report.start_at ? 'stopped' : 'not_started', 'Session stopped by participant/facilitator.')
    return 'Use Start before beginning the task. / Klik Start sebelum memulai.'
  })
  page = await context.newPage()
  video = page.video()
  page.on('pageerror', () => {
    if (report.start_at) finish('failed', 'Application JavaScript error. No sensitive error text logged.')
  })
  page.on('close', () => { if (!closing) finish(report.start_at ? 'stopped' : 'not_started', 'Browser page closed.') })
  page.on('request', request => {
    if (operation(request) !== 'draft') return
    let mode = 'unknown'
    try { mode = request.postDataJSON()?.mode || mode } catch {}
    if (!['manual', 'apertus'].includes(mode)) mode = 'unknown'
    const entry = { mode, started_at: new Date().toISOString(), duration_seconds: null, http_status: null }
    report.observed_draft_requests.push(entry)
    requests.set(request, entry)
  })
  page.on('requestfinished', request => {
    const entry = requests.get(request)
    if (entry) entry.duration_seconds = (Date.now() - Date.parse(entry.started_at)) / 1000
  })
  page.on('response', response => {
    const entry = requests.get(response.request())
    if (entry) entry.http_status = response.status()
    if (operation(response.request()) && report.start_at) {
      if (response.status() >= 400) report.events.push({ type: 'api_rejection', status: response.status(), at: new Date().toISOString() })
      if (response.status() >= 500) finish('failed', `API HTTP ${response.status()}. Recording stopped without automatic retry.`)
    }
  })
  page.on('requestfailed', request => {
    if (operation(request) && report.start_at && !closing) finish('failed', 'Local API request failed. No automatic retry.')
  })
  page.on('download', download => {
    const index = ++downloadIndex
    const filename = path.basename(download.suggestedFilename()).replace(/[^a-zA-Z0-9_.-]/g, '_')
    const entry = { filename, saved: false, path: `downloads/${index}-${filename}`, at: new Date().toISOString() }
    report.downloads.push(entry)
    const pending = download.saveAs(path.join(out, entry.path)).then(() => { entry.saved = true }).catch(() => {
      entry.saved = false
      if (!closing) finish('failed', 'Download could not be saved.')
    })
    pendingDownloads.add(pending)
    pending.finally(() => pendingDownloads.delete(pending))
  })
  await page.goto('http://127.0.0.1:5175/', { waitUntil: 'networkidle' })
  await page.evaluate(({ task, mode, caseId }) => {
    const host = document.createElement('div')
    host.id = 'human-session-recorder'
    host.style.cssText = 'position:fixed;bottom:8px;right:8px;z-index:2147483647;max-width:360px;'
    const shadow = host.attachShadow({ mode: 'open' })
    shadow.innerHTML = '<style>:host{font:13px Arial;color:#fff}details{padding:10px;background:#163c35;border:2px solid #e7b552;border-radius:10px}summary{cursor:pointer}button{margin:6px 4px 0 0;padding:7px;cursor:pointer}p{margin:8px 0;line-height:1.4}</style>'
    const panel = document.createElement('details'); panel.open = true
    const title = document.createElement('summary'); title.textContent = `Recorder · ${task} · ${caseId} ${mode}`
    const status = document.createElement('p'); status.textContent = 'Read brief, then Start. / Baca panduan, lalu Start. Collapse this panel if it covers the app.'
    panel.append(title, status)
    for (const [action, label] of [['start', 'Start'], ['pause', 'Pause / Resume'], ['complete', 'TXT opened; finish'], ['failed', 'Failure — stop'], ['stop', 'Stop unfinished']]) {
      const button = document.createElement('button'); button.textContent = label
      button.addEventListener('click', async () => { status.textContent = await window.pilotControl(action) })
      panel.append(button)
    }
    shadow.append(panel); document.body.append(host)
  }, { task: taskId, mode: assignment.mode, caseId: assignment.case_id })
  console.log(JSON.stringify({ status: 'waiting_for_person_to_click_Start', task: taskId, mode: assignment.mode, files: path.join(out, 'inputs'), output: out }))
  process.once('SIGINT', () => finish(report.start_at ? 'stopped' : 'not_started', 'Interrupted by facilitator.'))
  process.once('SIGTERM', () => finish(report.start_at ? 'stopped' : 'not_started', 'Terminated by facilitator.'))
  await end
} catch {
  finish('failed', 'Recorder setup or browser operation failed. Inspect local setup; no sensitive error text logged.')
  await end
} finally {
  await context?.close().catch(() => {})
  await Promise.allSettled([...pendingDownloads])
  if (video) {
    let stem = `LaporSiap_${taskId}_${report.status === 'completed_declared' ? 'Completed' : 'Stopped'}`
    if (report.status === 'failed') stem = 'LaporSiap_Failed_1'
    const webm = path.join(out, `${stem}.webm`)
    try {
      await video.saveAs(webm)
      report.recording = path.basename(webm)
      const conversion = spawnSync('ffmpeg', ['-hide_banner', '-loglevel', 'error', '-n', '-i', webm, '-c:v', 'libx264', '-preset', 'veryfast', '-crf', '18', '-pix_fmt', 'yuv420p', '-movflags', '+faststart', '-an', path.join(out, `${stem}.mp4`)], { encoding: 'utf8' })
      if (conversion.status === 0) report.recording = `${stem}.mp4`
    } catch { report.events.push({ type: 'recording_finalize_failed' }) }
  }
  await fs.writeFile(path.join(out, 'result.json'), JSON.stringify(report, null, 2) + '\n', { flag: 'wx' })
  await browser?.close().catch(() => {})
  console.log(JSON.stringify({ status: report.status, out, recording: report.recording }))
}
process.exitCode = report.status === 'failed' ? 1 : 0
