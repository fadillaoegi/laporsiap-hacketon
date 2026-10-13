import { defineConfig } from '@playwright/test'
import base from './playwright.config'
const stamp = new Date().toISOString().replaceAll(':', '-')
export default defineConfig(base, {
  testMatch: '**/deployment.spec.ts',
  testIgnore: [],
  outputDir: `../../output/hardening/deployment-${stamp}`,
  timeout: 60000,
  use: { baseURL: 'http://127.0.0.1:8877', video: 'on' },
})
