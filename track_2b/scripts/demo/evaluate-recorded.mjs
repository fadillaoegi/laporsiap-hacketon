// Two bounded fictional development rechecks. Fixture creation uses the local
// API; consent and the sole provider attempt are performed through the real UI.
import fs from 'node:fs/promises'
import path from 'node:path'
import crypto from 'node:crypto'
import { createRequire } from 'node:module'
import { spawnSync } from 'node:child_process'
import assert from 'node:assert/strict'
const require = createRequire(new URL('../../src/laporsiap_FE/package.json', import.meta.url))
const root = path.resolve(new URL('../..', import.meta.url).pathname)
process.env.PLAYWRIGHT_BROWSERS_PATH ||= path.join(root, 'tmp/playwright')
const { chromium, request, expect } = require('@playwright/test')
const argument = (name, fallback) => { const i=process.argv.indexOf(name); if(i<0) return fallback; assert.ok(process.argv[i+1] && !process.argv[i+1].startsWith('--'), 'Missing argument '+name); return process.argv[i+1] }
const origin = new URL(argument('--origin', 'http://127.0.0.1:5175'))
assert.ok(origin.protocol==='http:' && ['127.0.0.1','localhost'].includes(origin.hostname) && !origin.username && !origin.password && !origin.search && !origin.hash && origin.pathname==='/', 'Use a local origin, including a forwarded GPU app')
const raw = await fs.readFile(path.resolve(root, argument('--fixture', 'docs/evaluation/actor-direction-v6-time-preserved.json')))
const fixture = JSON.parse(raw)
assert.equal(fixture.synthetic, true)
assert.equal(fixture.items.length, 2)
const promptRaw = await fs.readFile(path.join(root, 'src/laporsiap_BE/internal/infrastructure/apertus/prompt.json'))
const hash = b => crypto.createHash('sha256').update(b).digest('hex')
if (process.argv.includes('--check')) {
 const browser = await chromium.launch({channel:'chrome',headless:true})
 const folder = path.join(root,'tmp','recorded-check',new Date().toISOString().replaceAll(':','-'))
 const context = await browser.newContext({recordVideo:{dir:folder}})
 const page = await context.newPage();await page.setContent('<p>Technical recorder check. No application/provider call.</p>')
 await context.close();await browser.close()
 console.log('Recorder check passed; no app/provider calls.')
 process.exit(0)
}
if (!process.argv.includes('--live')) {
  console.log('Validated two fictional cases; no browser or API calls. Use --live to record.')
  process.exit(0)
}
const out = path.join(root, 'output/evaluation', new Date().toISOString().replaceAll(':','-'))
await fs.mkdir(out, { recursive: true, mode: 0o700 })
const write = (name, data) => fs.writeFile(path.join(out, name), JSON.stringify(data, null, 2)+'\n', { flag:'wx', mode:0o600 })
const api = await request.newContext({ baseURL:origin.origin, timeout:55000 })
const policyResponse = await api.get('/api/policy')
assert.equal(policyResponse.status(), 200)
const policy = await policyResponse.json()
assert.equal(policy.prompt_sha256, hash(promptRaw))
assert.equal(policy.prompt_version, 'bilingual-v6-time-preserved')
await write('fixture-frozen.json', { sha256:hash(raw), version:fixture.version, items:fixture.items, note:'New session recheck of the same development examples. Prior failures are preserved.' })
await write('prompt-frozen.json', { version:policy.prompt_version, sha256:hash(promptRaw), prompt:JSON.parse(promptRaw), running_backend_verified:true })
const browser = await chromium.launch({ channel:'chrome', headless:true, slowMo:120 })
const results=[]
try {
 for (const item of fixture.items) {
  assert.equal(item.input.synthetic, true)
  const r={id:item.id,language:item.language,category:item.category,started_at:new Date().toISOString(),status:'not_started',api_status:0,contract_passed:false,state_verified:false,transfers_unchanged:false,original_sources_unchanged:false,failed_draft_not_saved:false,semantic_review_status:'not_scored',semantic_review_criteria:item.semantic_review_criteria,provider_attempts:0}
  const folder=path.join(out,item.id)
  await fs.mkdir(folder,{mode:0o700})
  let context, page, video, before, failed=false
  console.log('START',item.id)
  try {
   const created=await api.post('/api/cases',{data:item.input,headers:{'Accept-Language':item.language}})
   assert.equal(created.status(),201)
   before=await created.json()
   assert.ok(before.token)
   r.case_id=before.case.id; r.original_sources=before.case.messages
   context=await browser.newContext({viewport:{width:1920,height:1080},recordVideo:{dir:folder,size:{width:1920,height:1080}}})
   await context.addInitScript(({language,access})=>{localStorage.setItem('laporsiap-language',language);sessionStorage.setItem('laporsiap-access',JSON.stringify([access]))},{language:item.language,access:{id:before.case.id,token:before.token,title:before.case.title}})
   page=await context.newPage();video=page.video();page.setDefaultTimeout(12000)
   await page.goto(origin.href,{waitUntil:'networkidle'})
   await expect(page.getByRole('textbox',{name:item.language==='en'?'Conversation':'Percakapan',exact:true})).toHaveValue(item.input.raw_chat)
   await page.waitForTimeout(800)
   const previewPromise=page.waitForResponse(res=>new URL(res.url()).pathname.endsWith('/ai-preview') && res.request().method()==='POST')
   await page.getByRole('navigation').getByRole('button',{name:item.language==='en'?/Create draft/:/Susun draft/}).click()
   const previewResponse=await previewPromise
   assert.equal(previewResponse.status(),200);r.preview=await previewResponse.json()
   const consent=page.getByRole('checkbox',{name:item.language==='en'?/I consent to sending/:/Saya setuju mengirim/})
   await expect(consent).toBeEnabled();await consent.check()
   const responsePromise=page.waitForResponse(res=>new URL(res.url()).pathname.endsWith('/draft') && res.request().method()==='POST',{timeout:55000})
   const start=performance.now()
   r.provider_attempts=1
   await page.getByRole('button',{name:item.language==='en'?'Create draft with Apertus':'Susun draft dengan Apertus',exact:true}).click()
   const response=await responsePromise
   r.draft_request_seconds=(performance.now()-start)/1000;r.api_status=response.status()
   const data=await response.json()
   if (!response.ok()) {
    r.error=typeof data.error==='string'?data.error:'Local draft request rejected'
    const d=data.provider_diagnostic
    if(d && [d.requested_model,d.response_model].every(v=>v===undefined||/^[A-Za-z0-9][A-Za-z0-9./_:-]{0,199}$/.test(v))) r.provider_diagnostic={requested_model:d.requested_model,response_model:d.response_model,validation_rule:['invalid chronology JSON schema','trailing JSON after chronology','chronology event count differs from sources','chronology source order or identifier invalid','chronology quote differs from original source','chronology description empty or over 35 words','chronology relative time differs from source'].includes(d.validation_rule)?d.validation_rule:undefined}
    await expect(page.locator('.alert.error')).toBeVisible({timeout:2000})
    await page.screenshot({path:path.join(folder,'failure-screen.png'),fullPage:true})
    throw new Error(r.error)
   }
   r.draft=data.case.draft
   assert.equal(r.draft.model,'swiss-ai/apertus-v1.5-70b');assert.equal(r.draft.language,item.language)
   assert.equal(r.draft.prompt_sha256,hash(promptRaw));assert.equal(r.draft.prompt_version,policy.prompt_version)
   assert.equal(r.draft.events.length,before.case.messages.length)
   for(const [i,event] of r.draft.events.entries()) {
    assert.equal(event.source_id,before.case.messages[i].id);assert.equal(event.evidence,before.case.messages[i].text)
    assert.ok(event.description.trim());assert.ok(event.description.trim().split(/\s+/).length<=35)
   }
   if(item.required_timing) {
    const event=r.draft.events.find(e=>e.source_id===item.required_timing.source_id)
    assert.ok(event?.description.toLowerCase().includes(item.required_timing.expression.toLowerCase()),'Required source timing omitted')
    r.required_timing_preserved=true
   }
   assert.equal(data.summary.successful_known_total,item.expected_total)
   assert.equal(data.summary.successful_unknown_amount,item.expected_unknown_success)
   await expect(page.getByRole('textbox',{name:item.language==='en'?'Chronology M1':'Kronologi M1',exact:true})).toBeVisible()
   await page.screenshot({path:path.join(folder,'draft.png'),fullPage:true})
   await page.waitForTimeout(1800)
   r.contract_passed=true;r.status='contract_passed_review_pending';r.semantic_review_status='pending_assistant_review'
  } catch(error) {
   failed=true;r.status='failed_or_rejected';r.error=r.error||error.message
  } finally {
   // Stop recording first, including a failed provider call; no retry/fallback.
   if(context) await context.close()
  }
  if(before?.token) {
   try {
    const state=await api.get('/api/cases/'+before.case.id,{headers:{Authorization:'Bearer '+before.token}})
    assert.equal(state.status(),200)
    const after=await state.json()
    assert.deepEqual(after.case.transfers,before.case.transfers);r.transfers_unchanged=true
    assert.deepEqual(after.case.messages,before.case.messages);assert.equal(after.case.raw_chat,before.case.raw_chat);r.original_sources_unchanged=true
    r.failed_draft_not_saved=after.case.revision===before.case.revision && after.case.draft===null && after.case.reviewed_revision===null
    if(r.api_status>=400) assert.equal(r.failed_draft_not_saved,true)
    r.state_verified=true
   } catch(error) { r.state_error=error.message;r.status='failed_or_rejected';r.contract_passed=false;failed=true }
  }
  if(video) {
   const name='LaporSiap_'+(failed?'Failed':'Success')+'_1'
   const saved=path.join(folder,name+'.webm')
   await video.saveAs(saved)
   const converted=spawnSync('ffmpeg',['-hide_banner','-loglevel','error','-i',saved,'-c:v','libx264','-crf','23','-pix_fmt','yuv420p','-movflags','+faststart',path.join(folder,name+'.mp4')],{encoding:'utf8'})
   assert.equal(converted.status,0,'Video conversion failed')
   r.recording=path.relative(out,path.join(folder,name+'.mp4'))
  }
  await write(item.id+'.json',r);results.push(r);console.log('END',item.id,r.status)
 }
} finally { await browser.close();await api.dispose() }
await write('summary.json',{version:fixture.version,fixture_sha256:hash(raw),prompt_version:policy.prompt_version,prompt_sha256:hash(promptRaw),requested_model:'swiss-ai/apertus-v1.5-70b',cases:results.length,contract_passed:results.filter(r=>r.contract_passed).length,results,retry_count:0,fallback_count:0,independent_human_sessions:0,recorded_sessions:results.filter(r=>r.recording).length, provider_attempts:results.reduce((n,r)=>n+r.provider_attempts,0), limitations:['Same development examples rechecked; not independent holdout.','Fictional fixture creation used the local API; consent and live calls are recorded in the UI.','API model identity is metadata, not a weights audit.','Per-source semantic review required; contract acceptance is not accuracy.','Failure videos stop before read-only state verification.']})
console.log('Artifacts:',out)
if(results.some(r=>!r.contract_passed))process.exitCode=1
