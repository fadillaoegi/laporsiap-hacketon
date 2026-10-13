import { createRequire } from 'node:module'
import { fileURLToPath } from 'node:url'
import fs from 'node:fs/promises'
import path from 'node:path'
import { spawnSync } from 'node:child_process'
import { createHash } from 'node:crypto'

const repo = fileURLToPath(new URL('../../', import.meta.url))
const require = createRequire(new URL('../../src/laporsiap_FE/package.json', import.meta.url))
process.env.PLAYWRIGHT_BROWSERS_PATH ||= path.join(repo, 'tmp/playwright')
const { chromium, expect } = require('@playwright/test')
const arg = name => process.argv[process.argv.indexOf(name) + 1]
const id = arg('--scenario'), mode = arg('--mode')
if (!['S01', 'S02', 'S03'].includes(id) || !['manual', 'apertus'].includes(mode)) throw new Error('Provide --scenario S01|S02|S03 --mode manual|apertus')
const inputs = path.join(repo, 'output/testing/2026-10-10/inputs')
const fixture = JSON.parse(await fs.readFile(path.join(inputs, 'cases.json'), 'utf8')).find(c => c.id === id)
const out = path.join(repo, 'output/testing/2026-10-10', id, mode, new Date().toISOString().replaceAll(':', '-'))
await fs.mkdir(out, { recursive: true })
let browser, context, page, video, latest, summary, originalDraft, problem, active = true, stepName
let draftRequestStarted, draftRequestMs, result = 'Failed', recording, manualEntriesAuthored = 0
const steps = [], uploads = []
let rejectFailure
const failure = new Promise((_, reject) => { rejectFailure = reject })
failure.catch(() => {})
function fail(message) { if (!active || problem) return; problem = new Error(message); rejectFailure(problem) }
async function step(name, run) {
  stepName = name
  const record = { name, started_at: new Date().toISOString() }; steps.push(record)
  console.log(`START ${id} ${mode}: ${name}`)
  await Promise.race([run(), failure]); if (problem) throw problem
  record.status = 'passed'; record.finished_at = new Date().toISOString()
  console.log(`PASS ${id} ${mode}: ${name}`)
}
async function caption(text) {
  await page.evaluate(text => {
    let el = document.getElementById('test-caption')
    if (!el) { el = document.createElement('div'); el.id = 'test-caption'; el.style.cssText = 'position:fixed;bottom:12px;left:335px;right:22px;z-index:9999;padding:18px;background:#123c34ed;color:white;border-radius:12px;font:20px Arial;pointer-events:none'; document.body.append(el) }
    el.textContent = text
  }, `${id} / ${mode.toUpperCase()} / FICTIONAL DATA - ${text}`)
}
const click = locator => locator.click()
const sourceChecks = () => page.getByRole('checkbox', { name: 'I have compared this entry with its source', exact: true })
const started = new Date().toISOString()
try {
  browser = await chromium.launch({ channel: 'chrome', headless: true, slowMo: 100 })
  context = await browser.newContext({ viewport: { width: 1920, height: 1080 }, recordVideo: { dir: out, size: { width: 1920, height: 1080 } }, acceptDownloads: true })
  await context.addInitScript(() => { try { localStorage.setItem('laporsiap-language', 'en') } catch {} })
  page = await context.newPage(); video = page.video(); page.setDefaultTimeout(15000)
  page.on('pageerror', e => fail(`Application error: ${e.message}`))
  page.on('requestfailed', r => { if (new URL(r.url()).pathname.startsWith('/api/')) fail(`API request failed: ${r.failure()?.errorText}`) })
  page.on('response', async r => {
    if (!new URL(r.url()).pathname.startsWith('/api/')) return
    if (!r.ok()) { const error = await r.json().catch(() => ({})); fail(`API HTTP ${r.status()}: ${error.error || 'request failed'}`); return }
    try { const data = await r.json(); if (data.case) { latest = data.case; summary = data.summary } } catch {}
  })
  await step('Import conversation TXT and save transfer inputs', async () => {
    await page.goto('http://127.0.0.1:5175/', { waitUntil: 'networkidle' })
    await expect(page.locator('html')).toHaveAttribute('lang', 'en')
    await caption('Import original TXT; statuses and amounts come from user input.')
    await page.getByRole('textbox', { name: 'Case title', exact: true }).fill(fixture.title)
    await page.getByLabel('Conversation format').selectOption('whatsapp_bracket_dmy')
    await page.locator('input[type=file][accept=".txt,text/plain"]').setInputFiles(path.join(inputs, `${id}-conversation.txt`))
    await expect(page.getByRole('textbox', { name: 'Conversation', exact: true })).toHaveValue(fixture.raw_chat)
    await page.getByRole('checkbox', { name: 'This is simulated data', exact: true }).check()
    for (const [i, t] of fixture.transfers.entries()) {
      await click(page.getByRole('button', { name: '+ Add a transfer record', exact: true }))
      for (const [label, key] of [['Sending bank','from_bank'],['Sending account','from_account'],['Receiving bank','to_bank'],['Receiving account','to_account'],['Amount (rupiah)','amount'],['Note','note']]) {
        await page.getByLabel(label, { exact: true }).nth(i).fill(t[key] == null ? '' : String(t[key]))
      }
      await page.getByLabel('Transfer status').nth(i).selectOption(t.status)
    }
    await click(page.getByRole('button', { name: 'Save & continue', exact: true }))
    await expect(page.getByRole('heading', { name: 'Draft with Apertus', exact: true })).toBeVisible()
    await expect.poll(() => summary?.successful_known_total).toBe(fixture.expected_total)
    await expect.poll(() => summary?.successful_unknown_amount).toBe(fixture.expected_unknown_success)
    await expect(page.locator('.summary-card .total')).toContainText(fixture.expected_total == null ? 'Unknown' : fixture.expected_total.toLocaleString('en-US'))
    if (fixture.expected_unknown_success) await expect(page.locator('.summary-card')).toContainText('Partial total:')
  })
  await step('Upload PDF, PNG and JPEG; link and download each original', async () => {
    await click(page.getByRole('navigation').getByRole('button', { name: /Case materials/ }))
    await caption('Upload three real files. These fictional references are not bank verification.')
    for (const [i, filename] of fixture.attachments.entries()) {
      const file = path.join(inputs, filename), original = await fs.readFile(file)
      await page.locator('input[type=file][accept="image/png,image/jpeg,application/pdf"]').setInputFiles(file)
      const row = page.locator('.attachment-row').filter({ has: page.getByText(filename, { exact: true }) })
      await expect(row).toBeVisible(); await expect(page.locator('.attachment-row')).toHaveCount(i + 1)
      await click(row.getByRole('button', { name: 'Link evidence', exact: true }))
      const editor = page.getByRole('region', { name: 'Link attachment' })
      await editor.getByRole('checkbox', { name: /^M1 / }).check()
      if (!filename.endsWith('.png')) await editor.getByRole('checkbox', { name: /^Transfer 01 / }).check()
      await click(editor.getByRole('button', { name: 'Save evidence links', exact: true }))
      await expect(row).toContainText('M1')
      const pending = page.waitForEvent('download'); pending.catch(() => {})
      await click(row.getByRole('button', { name: 'Download', exact: true }))
      const downloaded = await pending, saved = path.join(out, filename)
      await downloaded.saveAs(saved)
      if (await downloaded.failure()) throw new Error('Attachment download failed')
      if (!original.equals(await fs.readFile(saved))) throw new Error('Attachment changed during round trip')
      uploads.push({ filename, size: original.length, sha256: createHash('sha256').update(original).digest('hex'), round_trip_identical: true })
    }
    await page.screenshot({ path: path.join(out, '01-uploaded-materials.png'), fullPage: true })
    await click(page.getByRole('navigation').getByRole('button', { name: /Create draft/ }))
  })
  await step('Create draft with consent gates', async () => {
    await expect(page.getByRole('button', { name: 'Create draft with Apertus', exact: true })).toBeDisabled()
    if (mode === 'apertus') {
      await caption('Read the sources and consent to one actual Apertus request. No automatic retry.')
      await click(page.getByText('View all messages to be sent', { exact: true }))
      await page.getByRole('checkbox', { name: /I consent to sending/ }).check()
    }
    draftRequestStarted = Date.now()
    await click(page.getByRole('button', { name: mode === 'apertus' ? 'Create draft with Apertus' : 'Use a manual draft', exact: true }))
    await expect(page.getByText(mode === 'apertus' ? 'APERTUS DRAFT' : 'MANUAL DRAFT', { exact: true })).toBeVisible({ timeout: 55000 })
    draftRequestMs = Date.now() - draftRequestStarted
    await expect.poll(() => latest?.draft?.origin).toBe(mode)
    originalDraft = structuredClone(latest.draft)
    if (mode === 'apertus' && (originalDraft.model !== 'swiss-ai/apertus-v1.5-70b' || originalDraft.language !== 'en')) throw new Error('Unexpected model or draft language')
    if (originalDraft.events.length !== fixture.descriptions.length) throw new Error('Wrong event count')
    for (const [i, e] of originalDraft.events.entries()) if (e.source_id !== latest.messages[i].id || e.evidence !== latest.messages[i].text) throw new Error('Source/order/quote changed')
    await fs.writeFile(path.join(out, 'original-draft.json'), JSON.stringify({ synthetic: true, draft: originalDraft, messages: latest.messages }, null, 2))
  })
  await step('Review meaning and confirm current revision', async () => {
    await caption('Compare every source. Unknown amounts and times must stay unknown.')
    if (mode === 'manual') {
      for (const [i, description] of fixture.descriptions.entries()) {
        await page.getByRole('textbox', { name: `Chronology M${i+1}`, exact: true }).fill(description)
        manualEntriesAuthored++
      }
      await click(page.getByRole('button', { name: 'Save edits', exact: true }))
      await expect(page.getByText('Edits saved.', { exact: false })).toBeVisible()
    } else {
      console.log(JSON.stringify({ review_required: true, out, messages: latest.messages, events: originalDraft.events, checks: originalDraft.checks }))
      const decisionFile = path.join(out, 'review-decision.json')
      const deadline = Date.now() + 300000
      let decision
      while (Date.now() < deadline && !problem) {
        decision = await fs.readFile(decisionFile, 'utf8').then(JSON.parse).catch(() => null)
        if (decision) break
        await page.waitForTimeout(1000)
      }
      if (decision?.approved !== true) throw new Error(decision?.reason || 'Assistant source review not approved before timeout')
    }
    await expect(sourceChecks()).toHaveCount(fixture.descriptions.length)
    await expect(page.getByRole('button', { name: 'Confirm & continue', exact: true })).toBeDisabled()
    for (let i = 0; i < fixture.descriptions.length; i++) await sourceChecks().nth(i).check()
    await page.getByRole('checkbox', { name: /I have checked the accounts/ }).check()
    await expect(page.getByRole('button', { name: 'Confirm & continue', exact: true })).toBeDisabled()
    await page.getByRole('checkbox', { name: /I have checked the review flags/ }).check()
    await click(page.getByRole('button', { name: 'Confirm & continue', exact: true }))
    await expect(page.getByRole('button', { name: 'Download ZIP package', exact: true })).toBeEnabled()
  })
  await step('Export reviewed ZIP and TXT; verify sources, transfers and attachments', async () => {
    await caption('Download the reviewed package and verify original sources and exact attachment bytes.')
    for (const [button, filename] of [['Download ZIP package', 'reviewed-package.zip'], ['Download chronology TXT', 'chronology.txt']]) {
      const pending = page.waitForEvent('download'); pending.catch(() => {})
      await click(page.getByRole('button', { name: button, exact: true }))
      const download = await pending; await download.saveAs(path.join(out, filename))
      if (await download.failure()) throw new Error('Export download failed')
    }
    const check = spawnSync('python3', ['-c', `
import json,sys
from pathlib import Path
from zipfile import ZipFile
folder,inputs,id,mode=sys.argv[1:]
f=next(c for c in json.loads((Path(inputs)/'cases.json').read_text()) if c['id']==id)
with ZipFile(Path(folder)/'reviewed-package.zip') as z:
 assert z.testzip() is None
 c=json.loads(z.read('kasus-direview.json'))
 assert c['synthetic'] and c['revision']==c['reviewed_revision']
 assert c['raw_chat']==f['raw_chat'] and c['draft']['origin']==mode
 assert len(c['messages'])==len(f['descriptions'])
 assert len(c['attachments'])==3
 assert len(c['draft']['events'])==len(c['messages'])
 for e,m in zip(c['draft']['events'],c['messages']):
  assert e['source_id']==m['id'] and e['evidence']==m['text']
 assert [{k:t[k] for k in expected} for t,expected in zip(c['transfers'],f['transfers'])]==f['transfers']
 assert len(c['transfers'])==len(f['transfers'])
 for i,filename in enumerate(f['attachments']):
  names=[n for n in z.namelist() if n.endswith(filename)]
  assert len(names)==1 and z.read(names[0])==(Path(inputs)/filename).read_bytes()
  a=next(a for a in c['attachments'] if a['filename']==filename)
  assert a['source_ids']==['M1']
  assert (a.get('transfer_ids') or [])==([] if filename.endswith('.png') else [c['transfers'][0]['id']])
 assert z.read('laporan.pdf').startswith(b'%PDF')
 assert z.read('percakapan-asli.txt').decode()==f['raw_chat']
 text=z.read('kronologi.txt').decode()
 assert 'SIMULATED DATA' in text and 'SUPPORTING REPORT MATERIALS' in text
 assert 'token' not in c and 'access_token' not in c
 print(json.dumps({'zip_integrity':True,'exact_quotes':True,'transfer_inputs_unchanged':True,'attachments_byte_identical':True,'evidence_links_preserved':True,'reviewed_revision':c['revision']}))
`, out, inputs, id, mode], { encoding: 'utf8' })
    if (check.status !== 0) throw new Error(`ZIP validation failed: ${check.stderr}`)
    await fs.writeFile(path.join(out, 'export-verification.json'), check.stdout)
    await page.screenshot({ path: path.join(out, '02-reviewed-export.png'), fullPage: true })
    await page.waitForTimeout(1200)
  })
  result = 'Success'
} catch (error) {
  problem ||= error
  if (steps.length) { steps.at(-1).status = 'failed'; steps.at(-1).error = problem.message }
  console.error(`STOP ${stepName}: ${problem.message}`)
  active = false
  if (context) await context.close().catch(() => {})
  context = undefined
} finally {
  active = false
  if (context) await context.close().catch(() => {})
  if (video) {
    const basename = `LaporSiap_${result}_1`
    await video.saveAs(path.join(out, `${basename}.webm`))
    const converted = spawnSync('ffmpeg', ['-hide_banner','-loglevel','error','-n','-i',path.join(out, `${basename}.webm`),'-vf','fps=30','-c:v','libx264','-preset','veryfast','-crf','18','-pix_fmt','yuv420p','-movflags','+faststart','-an',path.join(out, `${basename}.mp4`)], { encoding: 'utf8' })
    recording = `${basename}.${converted.status === 0 ? 'mp4' : 'webm'}`
  }
  const report = { scenario: id, mode, language:'en', synthetic:true, result, started_at:started, finished_at:new Date().toISOString(), recording, stopped_at:problem ? stepName : null, error:problem?.message || null, case_id:latest?.id || null, summary, original_draft:originalDraft, draft_request_ui_ms:draftRequestMs, manual_entries_authored:manualEntriesAuthored, ai_entries_corrected:mode==='apertus' && result==='Success' ? 0 : null, uploads, steps, limitation:'Developer browser automation with assistant source review. Includes deliberate browser delays and review waiting; not human task time or an independent accuracy/speed study.' }
  await fs.writeFile(path.join(out,'result.json'), JSON.stringify(report,null,2))
  if (browser) await browser.close().catch(() => {})
  console.log(JSON.stringify({result, out, recording}))
}
process.exitCode = result==='Success' ? 0 : 1
