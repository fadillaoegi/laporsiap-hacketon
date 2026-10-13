# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: strengthening.spec.ts >> evidence links survive reload, appear in exports, and source edits clear stale links
- Location: tests/strengthening.spec.ts:57:1

# Error details

```
Test timeout of 45000ms exceeded.
```

```
Error: locator.click: Test timeout of 45000ms exceeded.
Call log:
  - waiting for getByRole('button', { name: 'Use a manual draft' })

```

# Page snapshot

```yaml
- generic [ref=e2]:
  - banner [ref=e3]:
    - generic [ref=e4]:
      - generic [ref=e5]: LaporSiap
      - generic [ref=e6]: LOCAL PROTOTYPE
    - generic [ref=e8]:
      - generic [ref=e9]: Language / Bahasa
      - combobox "Language / Bahasa" [ref=e10]:
        - option "English" [selected]
        - option "Indonesia"
  - generic [ref=e11]:
    - complementary [ref=e12]:
      - generic [ref=e13]: PREPARATION SPACE
      - navigation "Report preparation steps" [ref=e14]:
        - button "01 Case materials Text, transfers, and evidence" [ref=e15] [cursor=pointer]:
          - generic [ref=e16]: "01"
          - generic [ref=e17]:
            - text: Case materials
            - generic [ref=e18]: Text, transfers, and evidence
        - button "02 Create draft Apertus or manual" [disabled] [ref=e19]:
          - generic [ref=e20]: "02"
          - generic [ref=e21]:
            - text: Create draft
            - generic [ref=e22]: Apertus or manual
        - button "03 Review & edit Compare with sources" [disabled] [ref=e23]:
          - generic [ref=e24]: "03"
          - generic [ref=e25]:
            - text: Review & edit
            - generic [ref=e26]: Compare with sources
        - button "04 Export materials For official channels" [disabled] [ref=e27]:
          - generic [ref=e28]: "04"
          - generic [ref=e29]:
            - text: Export materials
            - generic [ref=e30]: For official channels
      - generic [ref=e31]:
        - strong [ref=e35]: You stay in control.
        - paragraph [ref=e36]: The draft helps organize your materials. You review the contents and submit a report through official channels.
      - generic [ref=e37]: No identity card upload.No automatic report submission.
    - main [ref=e38]:
      - generic [ref=e39]:
        - paragraph [ref=e41]:
          - strong [ref=e42]: Need to report now?
          - text: You can access official channels immediately, without waiting for an AI draft.
        - link "Open IASC" [ref=e43] [cursor=pointer]:
          - /url: https://iasc.ojk.go.id/
      - generic [ref=e47]:
        - generic [ref=e48]: STEP 01 / 04
        - heading "Start with the materials you have." [level=1] [ref=e49]
        - paragraph [ref=e50]: An unorganized conversation can be a starting point. Leave unknown details blank.
      - generic [ref=e51]: SIMULATED DATA — identities, conversations and transactions are fictional.
      - alert [active] [ref=e52]:
        - text: Storage or revision capacity reached. Delete unneeded cases or contact the operator.
        - button "Reload saved case" [disabled] [ref=e53]
      - generic [ref=e54]:
        - generic [ref=e55]:
          - generic [ref=e56]:
            - generic [ref=e57]:
              - generic [ref=e62]:
                - heading "Conversation & initial account" [level=2] [ref=e63]
                - paragraph [ref=e64]: Paste text or import a TXT file.
              - button "Try an example" [ref=e65] [cursor=pointer]
            - group [ref=e66]:
              - generic [ref=e67]:
                - text: Case title
                - textbox "Case title" [ref=e68]:
                  - /placeholder: "Example: Online job offer and top-up requests"
                  - text: Online task offer and deposit requests
              - generic [ref=e69]:
                - generic [ref=e70]:
                  - text: Conversation format
                  - combobox "Conversation format" [ref=e71]:
                    - option "Plain text · one message per line"
                    - 'option "TXT · [DD/MM/YYYY, HH:MM] Sender: message" [selected]'
                - generic [ref=e72] [cursor=pointer]:
                  - text: Import TXT
                  - button "Import TXT" [ref=e75]
              - generic [ref=e76]:
                - text: Conversation
                - textbox "Conversation" [ref=e77]:
                  - /placeholder: "Paste messages here.\nSender and time are not inferred from plain text."
                  - text: "[08/10/2026, 09:10] Nadia: Hi Raka, we offer liking tasks with a promised commission of IDR 50,000 per day. [08/10/2026, 09:14] Nadia: To continue, please send a deposit of IDR 180,000 to ****4821. [08/10/2026, 09:28] Raka: I have just sent IDR 180,000. I have not recorded the transfer time. [08/10/2026, 09:35] Nadia: Now add a deposit of IDR 420,000 to unlock withdrawal. [08/10/2026, 09:40] Raka: The additional deposit is still pending. I have not checked my bank statement."
              - generic [ref=e78]: "Plain text: one message per line, without sender/time metadata. TXT with headers: the selected format must match exactly; continuation lines belong to the message. Maximum 100 messages / 32 KB."
              - generic [ref=e79] [cursor=pointer]:
                - checkbox "This is simulated data" [checked] [ref=e80]
                - text: This is simulated data
          - generic [ref=e81]:
            - generic [ref=e82]:
              - generic [ref=e83]: Rp
              - generic [ref=e85]:
                - heading "Transfer records" [level=2] [ref=e86]
                - paragraph [ref=e87]: Enter details from your records, rather than AI output.
              - generic [ref=e88]: 2 records
            - group [ref=e89]:
              - generic [ref=e90]:
                - strong [ref=e91]: Transfer 01
                - button "Remove record" [ref=e92] [cursor=pointer]
              - generic [ref=e93]:
                - generic [ref=e94]:
                  - text: Sending bank
                  - textbox "Sending bank" [ref=e95]:
                    - /placeholder: Unknown
                    - text: Bank Aruna
                - generic [ref=e96]:
                  - text: Sending account
                  - textbox "Sending account" [ref=e97]:
                    - /placeholder: Unknown
                    - text: "****7392"
                - generic [ref=e98]:
                  - text: Receiving bank
                  - textbox "Receiving bank" [ref=e99]:
                    - /placeholder: Unknown
                    - text: Bank Lestari
                - generic [ref=e100]:
                  - text: Receiving account
                  - textbox "Receiving account" [ref=e101]:
                    - /placeholder: Unknown
                    - text: "****4821"
                - generic [ref=e102]:
                  - text: Amount (rupiah)
                  - spinbutton "Amount (rupiah)" [ref=e103]: "180000"
                - generic [ref=e104]:
                  - text: Transfer status
                  - combobox "Transfer status" [ref=e105]:
                    - option "Successful" [selected]
                    - option "Pending"
                    - option "Failed"
                    - option "Unconfirmed"
                - generic [ref=e106]:
                  - text: Transfer date
                  - textbox "Transfer date May be blank. Chat time is not the transfer date." [ref=e107]
                  - generic [ref=e108]: May be blank. Chat time is not the transfer date.
                - generic [ref=e109]:
                  - text: Note
                  - textbox "Note" [ref=e110]:
                    - /placeholder: Corrections, conflicting accounts, etc.
                    - text: Transfer date is unknown; do not use the message time.
            - group [ref=e111]:
              - generic [ref=e112]:
                - strong [ref=e113]: Transfer 02
                - button "Remove record" [ref=e114] [cursor=pointer]
              - generic [ref=e115]:
                - generic [ref=e116]:
                  - text: Sending bank
                  - textbox "Sending bank" [ref=e117]:
                    - /placeholder: Unknown
                    - text: Bank Aruna
                - generic [ref=e118]:
                  - text: Sending account
                  - textbox "Sending account" [ref=e119]:
                    - /placeholder: Unknown
                    - text: "****7392"
                - generic [ref=e120]:
                  - text: Receiving bank
                  - textbox "Receiving bank" [ref=e121]:
                    - /placeholder: Unknown
                    - text: Bank Lestari
                - generic [ref=e122]:
                  - text: Receiving account
                  - textbox "Receiving account" [ref=e123]:
                    - /placeholder: Unknown
                    - text: "****4821"
                - generic [ref=e124]:
                  - text: Amount (rupiah)
                  - spinbutton "Amount (rupiah)" [ref=e125]: "420000"
                - generic [ref=e126]:
                  - text: Transfer status
                  - combobox "Transfer status" [ref=e127]:
                    - option "Successful"
                    - option "Pending" [selected]
                    - option "Failed"
                    - option "Unconfirmed"
                - generic [ref=e128]:
                  - text: Transfer date
                  - textbox "Transfer date May be blank. Chat time is not the transfer date." [ref=e129]
                  - generic [ref=e130]: May be blank. Chat time is not the transfer date.
                - generic [ref=e131]:
                  - text: Note
                  - textbox "Note" [ref=e132]:
                    - /placeholder: Corrections, conflicting accounts, etc.
                    - text: Not confirmed successful.
            - button "+ Add a transfer record" [ref=e133] [cursor=pointer]
          - generic [ref=e134]:
            - generic [ref=e139]:
              - heading "Evidence attachments" [level=2] [ref=e140]
              - paragraph [ref=e141]: Reference images or PDFs. OCR is not included.
            - generic [ref=e142]: Save your case materials before adding attachments.
          - generic [ref=e143]:
            - generic [ref=e144]: Saved in local PostgreSQL.
            - button "Save & continue" [ref=e145] [cursor=pointer]
        - complementary [ref=e148]:
          - generic [ref=e149]:
            - text: START STEP BY STEP
            - heading "You can fill in more later." [level=3] [ref=e150]
            - paragraph [ref=e151]: Keep the original messages. Enter the transfer details you know, then add evidence if available.
            - list [ref=e152]:
              - listitem [ref=e153]: Conversation text or TXT export
              - listitem [ref=e154]: Amounts, accounts, and transfer status
              - listitem [ref=e155]: Transfer or communication evidence
          - generic [ref=e156]:
            - heading "Know what will be sent." [level=3] [ref=e160]
            - paragraph [ref=e161]: Saving a case does not send text to AI. In the next step, choose a manual draft or consent to sending text to Public AI.
            - paragraph [ref=e162]: Attachments are for reference. Images and PDFs are not read automatically.
      - group [ref=e164]:
        - generic "Your data and access" [ref=e165] [cursor=pointer]
      - generic [ref=e166]:
        - text: LaporSiap · Supporting materials, with sources and human review.
        - generic [ref=e167]: Local development · Data stored on this device
```

# Test source

```ts
  1   | import { test, expect } from '@playwright/test'
  2   | import { mkdir, readFile } from 'node:fs/promises'
  3   | import path from 'node:path'
  4   | import type { Response } from '../src/domain/models'
  5   | 
  6   | test('AI preview keeps exact sources and language changes require fresh consent', async ({ page }) => {
  7   |   await page.goto('/')
  8   |   await page.getByRole('button', { name: 'Try an example', exact: true }).click()
  9   |   const conversation = await page.getByRole('textbox', { name: 'Conversation', exact: true }).inputValue()
  10  |   expect(conversation).toContain('I have just sent IDR 180,000.')
  11  |   expect(conversation).toContain('Nadia: Hi Raka,')
  12  |   expect(conversation).toContain('****4821')
  13  |   await expect(page.getByRole('checkbox', { name: 'This is simulated data', exact: true })).toBeChecked()
  14  |   await expect(page.locator('.simulation-banner')).toContainText('SIMULATED DATA')
  15  |   await expect(page.locator('body')).not.toContainText(/\bdemo\b/i)
  16  |   const created = page.waitForResponse(res => res.url().endsWith('/api/cases') && res.request().method() === 'POST')
  17  |   await page.getByRole('button', { name: 'Save & continue' }).click()
  18  |   const data: Response = await (await created).json()
  19  |   expect(data.case.synthetic).toBe(true)
  20  |   // Explicit original mode remains available after checking the disclosure.
  21  |   await page.getByRole('checkbox', { name: 'Mask personal information before sending to AI', exact: true }).uncheck()
  22  |   await page.getByText('View all messages to be sent', { exact: true }).click()
  23  |   const preview = page.locator('.preview-messages')
  24  |   await expect(preview.locator('blockquote')).toHaveCount(data.case.messages.length)
  25  |   for (let i = 0; i < data.case.messages.length; i++) await expect(preview.locator('blockquote').nth(i)).toHaveText(data.case.messages[i].text)
  26  |   await expect(preview).not.toContainText('****7392')
  27  |   const picker = page.getByRole('combobox', { name: 'Language / Bahasa' })
  28  |   await page.getByRole('checkbox', { name: /I consent to sending/ }).check()
  29  |   await picker.selectOption('id')
  30  |   await expect(page.locator('.simulation-banner')).toContainText('DATA SIMULASI')
  31  |   await expect(page.locator('body')).not.toContainText(/\bdemo\b/i)
  32  |   await expect(page.getByRole('checkbox', { name: /Saya setuju mengirim/ })).not.toBeChecked()
  33  |   await expect(page.getByRole('button', { name: 'Susun draft dengan Apertus', exact: true })).toBeDisabled()
  34  |   await picker.selectOption('en')
  35  |   await expect(preview.locator('blockquote').nth(0)).toHaveText(data.case.messages[0].text)
  36  |   await page.route('**/api/cases/*/draft', async route => {
  37  |     expect(route.request().headers()['accept-language']).toBe('en')
  38  |     expect(route.request().postDataJSON().external_consent).toBe(true)
  39  |     await route.fulfill({ json: { ...data, case: { ...data.case, draft: {
  40  |       origin: 'apertus', model: 'swiss-ai/apertus-v1.5-70b', language: 'en',
  41  |       checks: [{ event_index: 4, source_id: 'M5', flags: ['transfer_status_requires_review', 'relative_time_requires_review'] }],
  42  |       events: data.case.messages.map(m => ({ source_id: m.id, evidence: m.text, description: 'Synthetic response used only for browser regression.' })),
  43  |     } } } })
  44  |   })
  45  |   await page.getByRole('checkbox', { name: /I consent to sending/ }).check()
  46  |   await page.getByRole('button', { name: 'Create draft with Apertus', exact: true }).click()
  47  |   await expect(page.getByText('APERTUS DRAFT', { exact: true })).toBeVisible()
  48  |   await expect(page.getByText('The source describes an incomplete or failed transfer.', { exact: false })).toBeVisible()
  49  |   await expect(page.getByText('Check the timing in the message:', { exact: false })).toBeVisible()
  50  |   await expect(page.locator('.review-toolbar .draft-language')).toContainText('English')
  51  |   await picker.selectOption('id')
  52  |   await expect(page.locator('.review-toolbar .draft-language')).toContainText('English')
  53  |   await expect(page.getByText('Sumber menyebut status belum selesai atau gagal.', { exact: false })).toBeVisible()
  54  |   await expect(page.getByText('Periksa waktu dalam pesan:', { exact: false })).toBeVisible()
  55  | })
  56  | 
  57  | test('evidence links survive reload, appear in exports, and source edits clear stale links', async ({ page }) => {
  58  |   await page.goto('/')
  59  |   await page.getByRole('button', { name: 'Try an example', exact: true }).click()
  60  |   await page.getByRole('button', { name: 'Save & continue' }).click()
> 61  |   await page.getByRole('button', { name: 'Use a manual draft' }).click()
      |                                                                  ^ Error: locator.click: Test timeout of 45000ms exceeded.
  62  |   await page.getByRole('navigation').getByRole('button', { name: /Case materials/ }).click()
  63  |   const validPDF = await readFile(path.resolve('tests/fixtures/valid-reference.pdf'))
  64  |   await page.locator('input[type=file][accept="image/png,image/jpeg,application/pdf"]').setInputFiles({ name: 'linked-synthetic-reference.pdf', mimeType: 'application/pdf', buffer: validPDF })
  65  |   await expect(page.getByText('linked-synthetic-reference.pdf', { exact: true })).toBeVisible()
  66  |   await page.getByRole('button', { name: 'Link evidence', exact: true }).click()
  67  |   const editor = page.getByRole('region', { name: 'Link attachment' })
  68  |   await editor.getByRole('checkbox', { name: /^M3 / }).check()
  69  |   await editor.getByRole('checkbox', { name: /^Transfer 01 / }).check()
  70  |   await expect(page.getByRole('button', { name: 'Save & continue', exact: true })).toBeDisabled()
  71  |   await editor.getByRole('button', { name: 'Save evidence links', exact: true }).click()
  72  |   await expect(page.locator('.attachment-row')).toContainText('M3 · Transfer 01')
  73  |   await page.reload()
  74  |   await page.getByRole('navigation').getByRole('button', { name: /Case materials/ }).click()
  75  |   await expect(page.locator('.attachment-row')).toContainText('M3 · Transfer 01')
  76  |   await page.getByRole('button', { name: 'Link evidence', exact: true }).click()
  77  |   await expect(page.getByRole('region', { name: 'Link attachment' }).getByRole('checkbox', { name: /^M3 / })).toBeChecked()
  78  |   for (const width of [320, 390, 768, 1440]) {
  79  |     await page.setViewportSize({ width, height: 1000 })
  80  |     expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth), `Link editor overflow at ${width}`).toBe(true)
  81  |   }
  82  |   await page.screenshot({ path: '../../docs/screenshots/13-evidence-links.png', fullPage: true })
  83  |   await page.getByRole('button', { name: 'Cancel', exact: true }).click()
  84  |   await page.getByRole('navigation').getByRole('button', { name: /Review & edit/ }).click()
  85  |   await expect(page.locator('.event-card').nth(2)).toContainText('Linked attachment: linked-synthetic-reference.pdf')
  86  |   for (const check of await page.getByRole('checkbox', { name: 'I have compared this entry with its source' }).all()) await check.check()
  87  |   await page.getByRole('checkbox', { name: /I have checked the accounts/ }).check()
  88  |   await page.getByRole('checkbox', { name: /I have checked the review flags/ }).check()
  89  |   await page.getByRole('button', { name: 'Confirm & continue', exact: true }).click()
  90  |   await mkdir(path.resolve('../../output/pdf'), { recursive: true })
  91  |   for (const language of ['en', 'id']) {
  92  |     await page.getByRole('combobox', { name: 'Language / Bahasa' }).selectOption(language)
  93  |     const downloaded = page.waitForEvent('download')
  94  |     await page.getByRole('button', { name: language === 'en' ? 'Download PDF report' : 'Unduh laporan PDF', exact: true }).click()
  95  |     const file = await downloaded
  96  |     await file.saveAs(path.resolve(`../../output/pdf/LaporSiap-synthetic-${language}.pdf`))
  97  |     expect((await readFile(await file.path() as string)).subarray(0, 5).toString()).toBe('%PDF-')
  98  |   }
  99  |   await page.getByRole('combobox', { name: 'Language / Bahasa' }).selectOption('en')
  100 |   const textDownload = page.waitForEvent('download')
  101 |   await page.getByRole('button', { name: 'Download chronology TXT', exact: true }).click()
  102 |   const txt = await readFile(await (await textDownload).path() as string, 'utf8')
  103 |   expect(txt).toContain('User-selected links — messages: M3 | transfers:')
  104 |   await page.getByRole('navigation').getByRole('button', { name: /Case materials/ }).click()
  105 |   await page.getByRole('button', { name: 'Link evidence', exact: true }).click()
  106 |   await page.getByRole('button', { name: 'Save evidence links', exact: true }).click()
  107 |   await page.getByRole('navigation').getByRole('button', { name: /Export materials/ }).click()
  108 |   await expect(page.getByRole('button', { name: 'Download PDF report', exact: true })).toBeDisabled()
  109 |   await page.getByRole('navigation').getByRole('button', { name: /Case materials/ }).click()
  110 |   await page.getByRole('textbox', { name: 'Conversation', exact: true }).fill('Replacement source text, synthetic.')
  111 |   await page.getByLabel('Conversation format').selectOption('plain')
  112 |   await page.getByRole('button', { name: 'Save & continue', exact: true }).click()
  113 |   await page.getByRole('navigation').getByRole('button', { name: /Case materials/ }).click()
  114 |   await expect(page.locator('.attachment-row')).toContainText('No linked messages')
  115 | })
  116 | 
```