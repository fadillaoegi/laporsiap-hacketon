import { test, expect } from '@playwright/test'
import { readFile } from 'node:fs/promises'
import path from 'node:path'
import type { Response } from '../src/domain/models'

test('truncated attachments are rejected in both languages without invalidating reviewed data', async ({ page }) => {
  await page.goto('/')
  await page.getByRole('button', { name: 'Try an example', exact: true }).click()
  const created = page.waitForResponse(r => r.url().endsWith('/api/cases') && r.request().method() === 'POST')
  await page.getByRole('button', { name: 'Save & continue', exact: true }).click()
  const access: Response = await (await created).json()
  await page.getByRole('button', { name: 'Use a manual draft', exact: true }).click()
  await expect(page.getByText('MANUAL DRAFT', { exact: true })).toBeVisible()
  const sources = page.getByRole('checkbox', { name: 'I have compared this entry with its source', exact: true })
  await expect(sources).toHaveCount(access.case.messages.length)
  for (const checkbox of await sources.all()) await checkbox.check()
  await page.getByRole('checkbox', { name: /I have checked the accounts/ }).check()
  await page.getByRole('checkbox', { name: /I have checked the review flags/ }).check()
  await page.getByRole('button', { name: 'Confirm & continue', exact: true }).click()
  await expect(page.getByRole('button', { name: 'Download ZIP package', exact: true })).toBeEnabled()
  const getCase = async () => (await (await page.request.get(`/api/cases/${access.case.id}`, { headers: { Authorization: `Bearer ${access.token}` } })).json()).case
  const before = await getCase()
  await page.getByRole('navigation').getByRole('button', { name: /Case materials/ }).click()
  const pdf = await readFile(path.resolve('tests/fixtures/valid-reference.pdf'))
  const files = [
    { name: 'truncated.png', mimeType: 'image/png', buffer: Buffer.from('\x89PNG\r\n\x1a\n', 'binary') },
    { name: 'truncated.jpg', mimeType: 'image/jpeg', buffer: Buffer.from([255, 216, 255]) },
    { name: 'truncated.pdf', mimeType: 'application/pdf', buffer: pdf.subarray(0, 80) },
  ]
  for (const language of ['en', 'id']) {
    await page.getByRole('combobox', { name: 'Language / Bahasa' }).selectOption(language)
    for (const file of files) {
      await page.locator('input[type=file][accept="image/png,image/jpeg,application/pdf"]').evaluate(input => input.parentElement?.scrollIntoView({ block: 'center' }))
      const pending = page.waitForResponse(r => r.request().method() === 'POST' && r.url().endsWith('/attachments'))
      await page.locator('input[type=file][accept="image/png,image/jpeg,application/pdf"]').setInputFiles(file)
      expect((await pending).status()).toBe(400)
      await expect(page.getByRole('alert')).toContainText(language === 'en' ? (file.name.endsWith('.pdf') ? 'The PDF is incomplete' : 'The image is damaged') : (file.name.endsWith('.pdf') ? 'PDF tidak lengkap' : 'Gambar rusak'))
      await expect(page.getByRole('alert')).toBeInViewport()
      await expect(page.getByRole('alert')).toBeFocused()
      expect(await getCase()).toEqual(before)
      await expect(page.locator('.attachment-row')).toHaveCount(0)
    }
  }
  await page.getByRole('combobox', { name: 'Language / Bahasa' }).selectOption('en')
  await page.getByRole('navigation').getByRole('button', { name: /Export materials/ }).click()
  await expect(page.getByRole('button', { name: 'Download ZIP package', exact: true })).toBeEnabled()
})
