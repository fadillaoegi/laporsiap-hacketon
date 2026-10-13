# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: attachment-validation.spec.ts >> truncated attachments are rejected in both languages without invalidating reviewed data
- Location: tests/attachment-validation.spec.ts:6:1

# Error details

```
Test timeout of 45000ms exceeded.
```

```
Error: locator.click: Test timeout of 45000ms exceeded.
Call log:
  - waiting for getByRole('button', { name: 'Confirm & continue', exact: true })
    - locator resolved to <button disabled class="button primary">…</button>
  - attempting click action
    2 × waiting for element to be visible, enabled and stable
      - element is not enabled
    - retrying click action
    - waiting 20ms
    2 × waiting for element to be visible, enabled and stable
      - element is not enabled
    - retrying click action
      - waiting 100ms
    85 × waiting for element to be visible, enabled and stable
       - element is not enabled
     - retrying click action
       - waiting 500ms

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
        - button "Case materials Text, transfers, and evidence" [ref=e15] [cursor=pointer]:
          - generic [ref=e19]:
            - text: Case materials
            - generic [ref=e20]: Text, transfers, and evidence
        - button "Create draft Apertus or manual" [ref=e21] [cursor=pointer]:
          - generic [ref=e25]:
            - text: Create draft
            - generic [ref=e26]: Apertus or manual
        - button "03 Review & edit Compare with sources" [ref=e27] [cursor=pointer]:
          - generic [ref=e28]: "03"
          - generic [ref=e29]:
            - text: Review & edit
            - generic [ref=e30]: Compare with sources
        - button "04 Export materials For official channels" [ref=e31] [cursor=pointer]:
          - generic [ref=e32]: "04"
          - generic [ref=e33]:
            - text: Export materials
            - generic [ref=e34]: For official channels
      - generic [ref=e35]:
        - strong [ref=e39]: You stay in control.
        - paragraph [ref=e40]: The draft helps organize your materials. You review the contents and submit a report through official channels.
      - generic [ref=e41]:
        - generic [ref=e42]: Cases in this session
        - combobox "Cases in this session" [ref=e43]:
          - option "Select a case" [disabled]
          - option "Online task offer and deposit requests" [selected]
        - button "+ Create a new case" [ref=e44] [cursor=pointer]
      - generic [ref=e45]: No identity card upload.No automatic report submission.
    - main [ref=e46]:
      - generic [ref=e47]:
        - paragraph [ref=e49]:
          - strong [ref=e50]: Need to report now?
          - text: You can access official channels immediately, without waiting for an AI draft.
        - link "Open IASC" [ref=e51] [cursor=pointer]:
          - /url: https://iasc.ojk.go.id/
      - generic [ref=e54]:
        - generic [ref=e55]:
          - generic [ref=e56]: STEP 03 / 04
          - heading "Review every part of the story." [level=1] [ref=e57]
          - paragraph [ref=e58]: Compare the draft with the original messages. Correct any differences, then confirm your review.
        - generic [ref=e59]: Revision 2 · Review required
      - generic [ref=e60]: SIMULATED DATA — identities, conversations and transactions are fictional.
      - status [ref=e61]: The manual draft contains message quotes. Review and edit each entry.
      - generic [ref=e62]:
        - generic [ref=e63]:
          - strong [ref=e64]: "Source messages reviewed: 0 / 5"
          - paragraph [ref=e65]: Save your edits, check every source, then confirm the materials.
        - progressbar "Source messages reviewed" [ref=e66]
      - generic [ref=e67]:
        - generic [ref=e68]:
          - generic [ref=e69]: MANUAL DRAFT
          - generic [ref=e70]: 5 entries · check meaning, amounts in messages, and direction of funds
          - generic [ref=e71]: "Language when draft was created: English"
        - button "Save edits" [disabled] [ref=e72]
      - generic [ref=e73]:
        - generic [ref=e74]: EDITABLE DRAFT
        - generic [ref=e75]: ORIGINAL SOURCE MESSAGE
      - generic [ref=e76]:
        - generic [ref=e77]:
          - generic [ref=e78]:
            - generic [ref=e79]:
              - generic [ref=e80]: M1
              - generic [ref=e81]: Entry 01
            - generic [ref=e82]: Chronology M1
            - textbox "Chronology M1" [ref=e83]: "Message quote: Hi Raka, we offer liking tasks with a promised commission of IDR 50,000 per day."
            - generic [ref=e84] [cursor=pointer]:
              - checkbox "I have compared this entry with its source" [ref=e85]
              - text: I have compared this entry with its source
          - generic [ref=e86]:
            - generic [ref=e87]:
              - strong [ref=e88]: Nadia
              - generic [ref=e89]: 08/10/2026, 09:10
            - blockquote [ref=e90]: Hi Raka, we offer liking tasks with a promised commission of IDR 50,000 per day.
            - generic [ref=e91]: Full quote · M1 · message time is not transfer time
        - generic [ref=e92]:
          - generic [ref=e93]:
            - generic [ref=e94]:
              - generic [ref=e95]: M2
              - generic [ref=e96]: Entry 02
            - generic [ref=e97]: Chronology M2
            - textbox "Chronology M2" [ref=e98]: "Message quote: To continue, please send a deposit of IDR 180,000 to ****4821."
            - generic [ref=e99]: The source contains a request or promise. Check whether the draft presents it as an event that already happened.
            - generic [ref=e100] [cursor=pointer]:
              - checkbox "I have compared this entry with its source" [ref=e101]
              - text: I have compared this entry with its source
          - generic [ref=e102]:
            - generic [ref=e103]:
              - strong [ref=e104]: Nadia
              - generic [ref=e105]: 08/10/2026, 09:14
            - blockquote [ref=e106]: To continue, please send a deposit of IDR 180,000 to ****4821.
            - generic [ref=e107]: Full quote · M2 · message time is not transfer time
        - generic [ref=e108]:
          - generic [ref=e109]:
            - generic [ref=e110]:
              - generic [ref=e111]: M3
              - generic [ref=e112]: Entry 03
            - generic [ref=e113]: Chronology M3
            - textbox "Chronology M3" [ref=e114]: "Message quote: I have just sent IDR 180,000. I have not recorded the transfer time."
            - generic [ref=e115] [cursor=pointer]:
              - checkbox "I have compared this entry with its source" [ref=e116]
              - text: I have compared this entry with its source
          - generic [ref=e117]:
            - generic [ref=e118]:
              - strong [ref=e119]: Raka
              - generic [ref=e120]: 08/10/2026, 09:28
            - blockquote [ref=e121]: I have just sent IDR 180,000. I have not recorded the transfer time.
            - generic [ref=e122]: Full quote · M3 · message time is not transfer time
        - generic [ref=e123]:
          - generic [ref=e124]:
            - generic [ref=e125]:
              - generic [ref=e126]: M4
              - generic [ref=e127]: Entry 04
            - generic [ref=e128]: Chronology M4
            - textbox "Chronology M4" [ref=e129]: "Message quote: Now add a deposit of IDR 420,000 to unlock withdrawal."
            - generic [ref=e130] [cursor=pointer]:
              - checkbox "I have compared this entry with its source" [ref=e131]
              - text: I have compared this entry with its source
          - generic [ref=e132]:
            - generic [ref=e133]:
              - strong [ref=e134]: Nadia
              - generic [ref=e135]: 08/10/2026, 09:35
            - blockquote [ref=e136]: Now add a deposit of IDR 420,000 to unlock withdrawal.
            - generic [ref=e137]: Full quote · M4 · message time is not transfer time
        - generic [ref=e138]:
          - generic [ref=e139]:
            - generic [ref=e140]:
              - generic [ref=e141]: M5
              - generic [ref=e142]: Entry 05
            - generic [ref=e143]: Chronology M5
            - textbox "Chronology M5" [ref=e144]: "Message quote: The additional deposit is still pending. I have not checked my bank statement."
            - generic [ref=e145] [cursor=pointer]:
              - checkbox "I have compared this entry with its source" [ref=e146]
              - text: I have compared this entry with its source
          - generic [ref=e147]:
            - generic [ref=e148]:
              - strong [ref=e149]: Raka
              - generic [ref=e150]: 08/10/2026, 09:40
            - blockquote [ref=e151]: The additional deposit is still pending. I have not checked my bank statement.
            - generic [ref=e152]: Full quote · M5 · message time is not transfer time
      - generic [ref=e153]:
        - generic [ref=e154]:
          - generic [ref=e155]:
            - generic [ref=e156]:
              - heading "Review transfer records" [level=2] [ref=e157]
              - generic [ref=e158]: USER INPUT
            - table [ref=e160]:
              - rowgroup [ref=e161]:
                - row [ref=e162]:
                  - columnheader "Account pair" [ref=e163]
                  - columnheader "Amount & date" [ref=e164]
                  - columnheader "Status" [ref=e165]
              - rowgroup [ref=e166]:
                - row [ref=e167]:
                  - cell "Bank Aruna ****7392 ↓ Bank Lestari ****4821" [ref=e168]:
                    - strong [ref=e169]: Bank Aruna
                    - generic [ref=e170]: "****7392"
                    - generic [ref=e171]: ↓
                    - strong [ref=e172]: Bank Lestari
                    - generic [ref=e173]: "****4821"
                  - cell "IDR 180,000 Transfer date unknown Transfer date is unknown; do not use the message time." [ref=e174]:
                    - strong [ref=e175]: IDR 180,000
                    - generic [ref=e176]: Transfer date unknown
                    - generic [ref=e177]: Transfer date is unknown; do not use the message time.
                  - cell "Successful" [ref=e178]
                - row [ref=e180]:
                  - cell "Bank Aruna ****7392 ↓ Bank Lestari ****4821" [ref=e181]:
                    - strong [ref=e182]: Bank Aruna
                    - generic [ref=e183]: "****7392"
                    - generic [ref=e184]: ↓
                    - strong [ref=e185]: Bank Lestari
                    - generic [ref=e186]: "****4821"
                  - cell "IDR 420,000 Transfer date unknown Not confirmed successful." [ref=e187]:
                    - strong [ref=e188]: IDR 420,000
                    - generic [ref=e189]: Transfer date unknown
                    - generic [ref=e190]: Not confirmed successful.
                  - cell "Pending" [ref=e191]
            - generic [ref=e193] [cursor=pointer]:
              - checkbox "I have checked the accounts, amounts, dates, and transfer status." [checked] [ref=e194]
              - text: I have checked the accounts, amounts, dates, and transfer status.
          - generic [ref=e195]:
            - heading "Fill in or mark unknown details" [level=2] [ref=e196]
            - list [ref=e197]:
              - listitem [ref=e198]: "Transfer 1: missing transfer date."
              - listitem [ref=e199]: "Transfer 2: missing transfer date."
              - listitem [ref=e200]: No evidence attachments yet. Attachments can be added without OCR.
            - paragraph [ref=e201]: Heuristic checks may miss errors or flag valid paraphrases. Review flags do not establish truth.
            - generic [ref=e202] [cursor=pointer]:
              - checkbox "I have checked the review flags and incomplete details. I understand this export is supporting material, not an official report or bank verification." [checked] [active] [ref=e203]
              - text: I have checked the review flags and incomplete details. I understand this export is supporting material, not an official report or bank verification.
          - generic [ref=e204]:
            - generic [ref=e205]: 0 / 5 entries reviewed
            - button "Confirm & continue" [disabled] [ref=e206]
        - complementary [ref=e209]:
          - generic [ref=e210]:
            - text: TRANSFER RECORDS
            - paragraph [ref=e211]: Successful transfer amountreported by the user
            - strong [ref=e212]: IDR 180,000
            - generic [ref=e213]: Known amounts only. Not net losses or bank-verified amounts.
            - generic [ref=e214]:
              - generic [ref=e215]:
                - generic [ref=e216]: Successful
                - strong [ref=e218]: "1"
              - generic [ref=e219]:
                - generic [ref=e220]: Pending
                - strong [ref=e222]: "1"
              - generic [ref=e223]:
                - generic [ref=e224]: Failed
                - strong [ref=e226]: "0"
              - generic [ref=e227]:
                - generic [ref=e228]: Unconfirmed
                - strong [ref=e230]: "0"
            - generic [ref=e231]: Pending, failed, and unconfirmed transfers are excluded from the total.
          - generic [ref=e232]:
            - heading "Changes require a new review." [level=3] [ref=e233]
            - paragraph [ref=e234]: Edits to the draft, case input, or attachments clear the previous confirmation. Export is available only for a confirmed revision.
            - button "Edit case materials →" [ref=e235] [cursor=pointer]
      - generic [ref=e236]:
        - text: LaporSiap · Supporting materials, with sources and human review.
        - generic [ref=e237]: Local development · Data stored on this device
```

# Test source

```ts
  1  | import { test, expect } from '@playwright/test'
  2  | import { readFile } from 'node:fs/promises'
  3  | import path from 'node:path'
  4  | import type { Response } from '../src/types'
  5  | 
  6  | test('truncated attachments are rejected in both languages without invalidating reviewed data', async ({ page }) => {
  7  |   await page.goto('/')
  8  |   await page.getByRole('button', { name: 'Try an example', exact: true }).click()
  9  |   const created = page.waitForResponse(r => r.url().endsWith('/api/cases') && r.request().method() === 'POST')
  10 |   await page.getByRole('button', { name: 'Save & continue', exact: true }).click()
  11 |   const access: Response = await (await created).json()
  12 |   await page.getByRole('button', { name: 'Use a manual draft', exact: true }).click()
  13 |   for (const checkbox of await page.getByRole('checkbox', { name: 'I have compared this entry with its source', exact: true }).all()) await checkbox.check()
  14 |   await page.getByRole('checkbox', { name: /I have checked the accounts/ }).check()
  15 |   await page.getByRole('checkbox', { name: /I have checked the review flags/ }).check()
> 16 |   await page.getByRole('button', { name: 'Confirm & continue', exact: true }).click()
     |                                                                               ^ Error: locator.click: Test timeout of 45000ms exceeded.
  17 |   await expect(page.getByRole('button', { name: 'Download ZIP package', exact: true })).toBeEnabled()
  18 |   const getCase = async () => (await (await page.request.get(`/api/cases/${access.case.id}`, { headers: { Authorization: `Bearer ${access.token}` } })).json()).case
  19 |   const before = await getCase()
  20 |   await page.getByRole('navigation').getByRole('button', { name: /Case materials/ }).click()
  21 |   const pdf = await readFile(path.resolve('tests/fixtures/valid-reference.pdf'))
  22 |   const files = [
  23 |     { name: 'truncated.png', mimeType: 'image/png', buffer: Buffer.from('\x89PNG\r\n\x1a\n', 'binary') },
  24 |     { name: 'truncated.jpg', mimeType: 'image/jpeg', buffer: Buffer.from([255, 216, 255]) },
  25 |     { name: 'truncated.pdf', mimeType: 'application/pdf', buffer: pdf.subarray(0, 80) },
  26 |   ]
  27 |   for (const language of ['en', 'id']) {
  28 |     await page.getByRole('combobox', { name: 'Language / Bahasa' }).selectOption(language)
  29 |     for (const file of files) {
  30 |       const pending = page.waitForResponse(r => r.request().method() === 'POST' && r.url().endsWith('/attachments'))
  31 |       await page.locator('input[type=file][accept="image/png,image/jpeg,application/pdf"]').setInputFiles(file)
  32 |       expect((await pending).status()).toBe(400)
  33 |       await expect(page.getByRole('alert')).toContainText(language === 'en' ? (file.name.endsWith('.pdf') ? 'The PDF is incomplete' : 'The image is damaged') : (file.name.endsWith('.pdf') ? 'PDF tidak lengkap' : 'Gambar rusak'))
  34 |       expect(await getCase()).toEqual(before)
  35 |       await expect(page.locator('.attachment-row')).toHaveCount(0)
  36 |     }
  37 |   }
  38 |   await page.getByRole('combobox', { name: 'Language / Bahasa' }).selectOption('en')
  39 |   await page.getByRole('navigation').getByRole('button', { name: /Export materials/ }).click()
  40 |   await expect(page.getByRole('button', { name: 'Download ZIP package', exact: true })).toBeEnabled()
  41 | })
  42 | 
```