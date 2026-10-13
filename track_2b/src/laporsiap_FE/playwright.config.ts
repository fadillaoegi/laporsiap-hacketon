import { defineConfig } from '@playwright/test'
import { existsSync } from 'node:fs'
import { fileURLToPath } from 'node:url'

// Use the recorder's installed tools when present; retain Playwright's
// normal cache location on other developer machines.
const localBrowsers = fileURLToPath(new URL('../../tmp/playwright/', import.meta.url))
if (existsSync(localBrowsers)) process.env.PLAYWRIGHT_BROWSERS_PATH ||= localBrowsers

export default defineConfig({
  testDir: './tests',
  testIgnore: ['**/deployment.spec.ts', '**/unit/**'],
  timeout: 45000,
  workers: 1,
  use: {
    baseURL: 'http://127.0.0.1:5175',
    channel: 'chrome',
    headless: true,
    viewport: { width: 1440, height: 1050 },
    screenshot: 'only-on-failure',
    video: 'on',
  },
})
