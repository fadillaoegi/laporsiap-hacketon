import { test, expect } from '@playwright/test'
import { readFile, writeFile } from 'node:fs/promises'
import { spawnSync } from 'node:child_process'
import path from 'node:path'
import type { Response } from '../src/domain/models'

test.skip(({ baseURL }) => baseURL !== 'http://127.0.0.1:8877', 'Separate compiled-container verification')
test.afterEach(async ({ page }, info) => {
  if (info.status === 'skipped') return
  const video = page.video()
  await page.context().close()
  if (!video) return
  const name = info.status === 'passed' ? 'LaporSiap_Success_1' : 'LaporSiap_Failed_1'
  const file = info.outputPath(`${name}.webm`)
  await video.saveAs(file)
  const ff = spawnSync('ffmpeg', ['-hide_banner', '-loglevel', 'error', '-n', '-i', file, '-vf', 'fps=30', '-c:v', 'libx264', '-preset', 'veryfast', '-crf', '20', '-pix_fmt', 'yuv420p', '-movflags', '+faststart', '-an', info.outputPath(`${name}.mp4`)], { encoding: 'utf8' })
  if (ff.status !== 0) throw new Error('Recording conversion failed; original WebM preserved')
})

for (const language of ['en', 'id'] as const) {
  test(`compiled private deployment: complete manual workflow in ${language}`, async ({ page }, info) => {
    const errors: string[] = []
    page.on('pageerror', e => errors.push(e.message))
    page.on('console', m => { if (/Content Security Policy|Refused to/.test(m.text())) errors.push(m.text()) })
    const response = await page.goto('/')
    expect(response?.status()).toBe(200)
    expect(response?.headers()['content-security-policy']).toContain("frame-ancestors 'none'")
    await page.getByRole('combobox', { name: 'Language / Bahasa' }).selectOption(language)
    const name = (en: string, id: string) => language === 'en' ? en : id
    await page.getByRole('button', { name: name('Try an example', 'Coba contoh'), exact: true }).click()
    const pending = page.waitForResponse(r => r.request().method() === 'POST' && r.url().endsWith('/api/cases'))
    await page.getByRole('button', { name: name('Save & continue', 'Simpan & lanjut'), exact: true }).click()
    const saved: Response = await (await pending).json()
    await page.getByRole('button', { name: name('Use a manual draft', 'Gunakan draft manual'), exact: true }).click()
    await page.getByRole('navigation').getByRole('button', { name: language === 'en' ? /Case materials/ : /Bahan kasus/ }).click()
    const reference = path.resolve('tests/fixtures/valid-reference.pdf')
    await page.locator('input[type=file][accept="image/png,image/jpeg,application/pdf"]').setInputFiles(reference)
    await expect(page.locator('.attachment-row')).toHaveCount(1)
    await page.getByRole('navigation').getByRole('button', { name: language === 'en' ? /Review & edit/ : /Periksa & koreksi/ }).click()
    const checks = page.getByRole('checkbox', { name: name('I have compared this entry with its source', 'Saya telah mencocokkan entri ini dengan sumber'), exact: true })
    await expect(checks).toHaveCount(5)
    for (const check of await checks.all()) await check.check()
    await page.getByRole('checkbox', { name: language === 'en' ? /I have checked the accounts/ : /Saya telah memeriksa rekening/ }).check()
    await page.getByRole('checkbox', { name: language === 'en' ? /I have checked the review flags/ : /Saya telah memeriksa tanda review/ }).check()
    await page.getByRole('button', { name: name('Confirm & continue', 'Konfirmasi & lanjut'), exact: true }).click()
    for (const [en, id, ext] of [['Download ZIP package', 'Unduh paket ZIP', 'zip'], ['Download chronology TXT', 'Unduh kronologi TXT', 'txt']] as const) {
      const pending = page.waitForEvent('download')
      await page.getByRole('button', { name: name(en, id), exact: true }).click()
      const downloaded = await pending
      await downloaded.saveAs(info.outputPath(`reviewed-${language}.${ext}`))
      const bytes = await readFile(await downloaded.path() as string)
      expect(bytes.length).toBeGreaterThan(100)
      if (ext === 'zip') expect(bytes.subarray(0, 2).toString()).toBe('PK')
      else expect(bytes.toString()).toContain(language === 'en' ? 'SIMULATED DATA' : 'DATA SIMULASI')
    }
    await page.getByText(name('Your data and access', 'Data dan akses Anda'), { exact: true }).click()
    const rotate = page.waitForResponse(r => r.url().endsWith('/access/rotate'))
    await page.getByRole('button', { name: name('Renew session access', 'Perbarui akses sesi'), exact: true }).click()
    const renewed: Response = await (await rotate).json()
    expect(renewed.token !== saved.token).toBe(true)
    expect((await page.request.get(`/api/cases/${saved.case.id}`, { headers: { Authorization: `Bearer ${saved.token}` } })).status()).toBe(404)
    await expect(page.getByRole('button', { name: name('Download ZIP package', 'Unduh paket ZIP'), exact: true })).toBeEnabled()
    await page.getByRole('textbox', { name: name('Type DELETE to delete', 'Ketik DELETE untuk menghapus'), exact: true }).fill('DELETE')
    await page.getByRole('button', { name: name('Permanently delete case', 'Hapus kasus secara permanen'), exact: true }).click()
    await expect(page.getByRole('status')).toContainText(name('deleted from the database', 'dihapus dari database'))
    expect((await page.request.get(`/api/cases/${saved.case.id}`, { headers: { Authorization: `Bearer ${renewed.token}` } })).status()).toBe(404)
    expect(errors).toEqual([])
    await writeFile(info.outputPath('verification.json'), JSON.stringify({ language, compiled_frontend: true, csp_violations: errors, workflow: 'manual upload review export rotation deletion', live_ai_calls: 0, status: 'passed' }, null, 2))
  })
}
