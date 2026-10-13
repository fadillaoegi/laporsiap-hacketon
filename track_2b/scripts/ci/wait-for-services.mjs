// Bounded read-only readiness check, shared by CI and local verification.
import { setTimeout as delay } from 'node:timers/promises'
const end = Date.now() + 45000
while (Date.now() < end) {
  try {
    const health = await fetch('http://127.0.0.1:8787/api/health', { signal: AbortSignal.timeout(3000), redirect: 'error' })
    if (!health.ok || (await health.json()).status !== 'ok') throw new Error('Backend not ready')
    const frontend = await fetch('http://127.0.0.1:5175/', { signal: AbortSignal.timeout(3000), redirect: 'error' })
    if (!frontend.ok || !(await frontend.text()).includes('id="root"')) throw new Error('Frontend not ready')
    console.log('Local backend/database and frontend ready; no case writes/provider calls.')
    process.exit(0)
  } catch {
    await delay(500)
  }
}
console.error('Local services did not become ready within 45 seconds. Inspect CI startup logs.')
process.exit(1)
