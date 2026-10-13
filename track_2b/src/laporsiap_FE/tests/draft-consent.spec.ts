import { test, expect } from '@playwright/test'
import type { Response } from '../src/domain/models'

test.beforeEach(async ({ page }) => {
  await page.addInitScript(() => localStorage.setItem('laporsiap-language', 'id'))
})

for (const outcome of ['success', 'failure'] as const) {
  test(`each AI attempt requires fresh consent after ${outcome}`, async ({ page }) => {
    await page.goto('/')
    await page.getByRole('button', { name: 'Coba contoh', exact: true }).click()
    const created = page.waitForResponse(res => res.url().endsWith('/api/cases') && res.request().method() === 'POST')
    await page.getByRole('button', { name: 'Simpan & lanjut' }).click()
    const data: Response = await (await created).json()
    let requests = 0
    // Mock only the UI response: API/provider contracts are tested in Go.
    // No chat is sent to the external provider by this browser test.
    await page.route('**/api/cases/*/draft', async route => {
      requests++
      expect(route.request().postDataJSON().external_consent).toBe(true)
      if (outcome === 'failure') {
        await route.fulfill({ status: 502, json: { error: 'Provider sintetis tidak tersedia. Gunakan draft manual.' } })
      } else {
        await route.fulfill({ json: { ...data, case: { ...data.case, draft: {
          origin: 'apertus', model: 'swiss-ai/apertus-v1.5-70b', checks: [],
          events: data.case.messages.map(m => ({ source_id: m.id, evidence: m.text, description: 'Ringkasan sintetis untuk uji UI.' })),
        } } } })
      }
    })
    const consent = page.getByRole('checkbox', { name: /Saya setuju mengirim/ })
    const ai = page.getByRole('button', { name: /Susun (ulang|draft) dengan Apertus/ })
    await expect(ai).toBeDisabled()
    await consent.check()
    await ai.click()
    if (outcome === 'success') {
      await expect(page.getByText('DRAFT APERTUS', { exact: true })).toBeVisible()
      await page.getByRole('button', { name: /Susun draft.*Apertus atau manual/ }).click()
    } else {
      await expect(page.getByRole('alert')).toContainText('Provider sintetis tidak tersedia')
      await expect(page.getByRole('button', { name: 'Gunakan draft manual' })).toBeEnabled()
    }
    await expect(consent).not.toBeChecked()
    await expect(ai).toBeDisabled()
    expect(requests).toBe(1)
  })
}

test('AI byte limit handles multibyte text while manual drafting stays available', async ({ page }) => {
  await page.goto('/')
  await page.getByRole('textbox', { name: 'Judul kasus', exact: true }).fill('Simulasi batas UTF-8')
  // 8,405 UTF-16 units, 16,805 UTF-8 bytes, and six messages under 6 KB each.
  await page.getByRole('textbox', { name: 'Percakapan', exact: true }).fill(Array(6).fill('🌱'.repeat(700)).join('\n'))
  await page.getByRole('checkbox', { name: 'Ini data simulasi' }).check()
  await page.getByRole('button', { name: 'Simpan & lanjut' }).click()
  await expect(page.getByText('Percakapan ini melebihi batas draft AI', { exact: false })).toBeVisible()
  await page.getByRole('checkbox', { name: /Saya setuju mengirim/ }).check()
  await expect(page.getByRole('button', { name: 'Susun draft dengan Apertus', exact: true })).toBeDisabled()
  await page.getByRole('button', { name: 'Gunakan draft manual' }).click()
  await expect(page.getByText('DRAFT MANUAL', { exact: true })).toBeVisible()
  await expect(page.getByRole('textbox', { name: /^Kronologi M/ })).toHaveCount(6)
})
