import type { Access, Response } from '../../domain/models'
import type { SessionVault } from '../../application/ports'
export class BrowserSessionVault implements SessionVault {
 private values: Access[]
 constructor() { try { const parsed: unknown = JSON.parse(sessionStorage.getItem('laporsiap-access') || '[]'); this.values = Array.isArray(parsed) ? parsed.filter((a): a is Access => !!a && typeof a.id === 'string' && typeof a.token === 'string' && typeof a.title === 'string') : [] } catch { this.values = [] } }
 get(id: string) { const access = this.values.find(a => a.id === id); if (!access) throw new Error('Kasus tidak ditemukan atau akses tidak valid.'); return { ...access } }
 records() { return this.values.map(({ id, title }) => ({ id, title })) }
 remember(response: Response) { const existing = this.values.find(a => a.id === response.case.id); const token = response.token || existing?.token; if (!token) return; const access = { id: response.case.id, title: response.case.title, token }; this.values = existing ? this.values.map(a => a.id === access.id ? access : a) : [...this.values, access]; this.persist() }
 forget(id: string) { this.values = this.values.filter(a => a.id !== id); this.persist() }
 private persist() { try { sessionStorage.setItem('laporsiap-access', JSON.stringify(this.values)) } catch { /* Keep working in memory if tab storage is disabled. */ } }
}
