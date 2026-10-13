import type { Case, CaseResponse, Event, Input, Language, PrivacyOptions, Response } from '../domain/models'
import { validateInput } from '../domain/case'
import type { CaseRepository, FileSaver, SessionVault } from './ports'

// Pure orchestration over injected ports. Browser storage/network/Redux/React
// do not enter this layer. Credentials never enter returned UI state.
export class CaseWorkflow {
 private readonly repository: CaseRepository
 private readonly sessions: SessionVault
 private readonly files: FileSaver
 constructor(repository: CaseRepository, sessions: SessionVault, files: FileSaver) { this.repository = repository; this.sessions = sessions; this.files = files }
 records() { return this.sessions.records() }
 private accept(response: Response): CaseResponse { this.sessions.remember(response); return { case: response.case, summary: response.summary } }
 async load(id: string, language: Language) { return this.accept(await this.repository.load(id, language)) }
 async save(c: Case | null, input: Input, language: Language) { validateInput(input); return this.accept(c ? await this.repository.update(c.id, c.revision, input, language) : await this.repository.create(input, language)) }
 preview(c: Case, privacy: PrivacyOptions, language: Language) { return this.repository.preview(c.id, c.revision, privacy, language) }
 async generate(c: Case, mode: string, consent: boolean, privacy: PrivacyOptions, hash: string | undefined, language: Language) { return this.accept(await this.repository.generate(c.id, { revision: c.revision, mode, external_consent: mode === 'apertus' && consent, ...(mode === 'apertus' ? { privacy, preview_hash: hash } : {}) }, language)) }
 async edit(c: Case, events: Event[], language: Language) { return this.accept(await this.repository.editDraft(c.id, c.revision, events, language)) }
 async review(c: Case, sources: boolean, transfers: boolean, limits: boolean, language: Language) { return this.accept(await this.repository.review(c.id, c.revision, sources, transfers, limits, language)) }
 async upload(c: Case, file: File, language: Language) { if (file.size > 5 * 1024 * 1024) throw new Error('Maksimal ukuran lampiran 5 MB.'); return this.accept(await this.repository.upload(c.id, c.revision, file, language)) }
 async link(c: Case, attachment: string, sources: string[], transfers: string[], language: Language) { return this.accept(await this.repository.linkAttachment(c.id, attachment, c.revision, sources, transfers, language)) }
 async removeAttachment(c: Case, attachment: string, language: Language) { return this.accept(await this.repository.deleteAttachment(c.id, attachment, c.revision, language)) }
 async rotate(c: Case, language: Language) { const response = await this.repository.rotate(c.id, c.revision, language); if (!response.token) throw new Error('Terjadi kesalahan.'); return this.accept(response) }
 async delete(c: Case, language: Language) { await this.repository.delete(c.id, c.revision, language); this.sessions.forget(c.id) }
 policy(language: Language) { return this.repository.policy(language) }
 history(c: Case, language: Language) { return this.repository.history(c.id, language) }
 async downloadAttachment(c: Case, attachment: { id: string; filename: string }, language: Language) { this.files.save(await this.repository.attachment(c.id, attachment.id, language), attachment.filename) }
 async export(c: Case, format: 'zip' | 'txt' | 'pdf', language: Language) { this.files.save(await this.repository.export(c.id, format, language), `laporsiap-revisi-${c.revision}.${format}`) }
 async readConversation(file: File) { if (file.size > 32000) throw new Error('TXT maksimal 32 KB.'); const text = await file.text(); if (text.includes('\uFFFD')) throw new Error('TXT harus menggunakan encoding UTF-8.'); return text }
}
