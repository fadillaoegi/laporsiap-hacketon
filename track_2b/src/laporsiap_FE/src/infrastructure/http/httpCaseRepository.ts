import type { AIPreview, Event, History, Input, Language, PrivacyOptions, Response, SafetyPolicy } from '../../domain/models'
import type { CaseRepository, SessionVault } from '../../application/ports'
export class HttpCaseRepository implements CaseRepository {
 private readonly sessions: SessionVault
 constructor(sessions: SessionVault) { this.sessions = sessions }
 private async response(path: string, language: Language, method = 'GET', body?: unknown, id?: string, revision?: number) {
  const headers: Record<string, string> = { 'Accept-Language': language }
  if (id) headers.Authorization = `Bearer ${this.sessions.get(id).token}`
  if (revision !== undefined) headers['X-Case-Revision'] = String(revision)
  const form = body instanceof FormData
  if (body !== undefined && !form) headers['Content-Type'] = 'application/json'
  let response: globalThis.Response
  try { response = await fetch(`/api${path}`, { method, headers, body: body === undefined ? undefined : form ? body : JSON.stringify(body) }) }
  catch { throw new Error('Backend belum terhubung. Periksa server lokal; kanal resmi tetap dapat diakses.') }
  if (!response.ok) { const data = await response.json().catch(() => ({})); throw new Error(data.error || `Permintaan gagal (${response.status}).`) }
  return response
 }
 private async json<T>(path: string, language: Language, method = 'GET', body?: unknown, id?: string, revision?: number): Promise<T> { return (await this.response(path, language, method, body, id, revision)).json() }
 load(id: string, language: Language) { return this.json<Response>(`/cases/${id}`, language, 'GET', undefined, id) }
 create(input: Input, language: Language) { return this.json<Response>('/cases', language, 'POST', input) }
 update(id: string, revision: number, input: Input, language: Language) { return this.json<Response>(`/cases/${id}`, language, 'PUT', { revision, input }, id) }
 preview(id: string, revision: number, privacy: PrivacyOptions, language: Language) { return this.json<AIPreview>(`/cases/${id}/ai-preview`, language, 'POST', { revision, privacy }, id) }
 generate(id: string, data: { revision: number; mode: string; external_consent: boolean; privacy?: PrivacyOptions; preview_hash?: string }, language: Language) { return this.json<Response>(`/cases/${id}/draft`, language, 'POST', data, id) }
 editDraft(id: string, revision: number, events: Event[], language: Language) { return this.json<Response>(`/cases/${id}/draft`, language, 'PUT', { revision, events }, id) }
 review(id: string, revision: number, sources: boolean, transfers: boolean, limits: boolean, language: Language) { return this.json<Response>(`/cases/${id}/review`, language, 'POST', { revision, sources_reviewed: sources, transfers_reviewed: transfers, limitations_accepted: limits }, id) }
 upload(id: string, revision: number, file: File, language: Language) { const form = new FormData(); form.append('file', file); return this.json<Response>(`/cases/${id}/attachments`, language, 'POST', form, id, revision) }
 linkAttachment(id: string, attachment: string, revision: number, sources: string[], transfers: string[], language: Language) { return this.json<Response>(`/cases/${id}/attachments/${attachment}`, language, 'PUT', { revision, source_ids: sources, transfer_ids: transfers }, id) }
 deleteAttachment(id: string, attachment: string, revision: number, language: Language) { return this.json<Response>(`/cases/${id}/attachments/${attachment}`, language, 'DELETE', undefined, id, revision) }
 rotate(id: string, revision: number, language: Language) { return this.json<Response>(`/cases/${id}/access/rotate`, language, 'POST', { revision }, id) }
 async delete(id: string, revision: number, language: Language) { await this.json(`/cases/${id}`, language, 'DELETE', { revision, confirm_delete: true }, id) }
 policy(language: Language) { return this.json<SafetyPolicy>('/policy', language) }
 async history(id: string, language: Language) { return (await this.json<{ history: History[] }>(`/cases/${id}/history`, language, 'GET', undefined, id)).history }
 async attachment(id: string, attachment: string, language: Language) { return (await this.response(`/cases/${id}/attachments/${attachment}`, language, 'GET', undefined, id)).blob() }
 async export(id: string, format: 'zip' | 'txt' | 'pdf', language: Language) { return (await this.response(`/cases/${id}/export?format=${format}`, language, 'GET', undefined, id)).blob() }
}
