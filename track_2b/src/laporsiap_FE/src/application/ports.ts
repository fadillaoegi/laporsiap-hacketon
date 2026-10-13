import type { Access, AIPreview, CaseAccess, Event, History, Input, Language, PrivacyOptions, Response, SafetyPolicy } from '../domain/models'
export interface SessionVault { get(id: string): Access; records(): CaseAccess[]; remember(response: Response): void; forget(id: string): void }
export interface LanguagePreferences { read(): Language; write(language: Language): void }
export interface FileSaver { save(bytes: Blob, filename: string): void }
export interface CaseRepository {
 load(id: string, language: Language): Promise<Response>
 create(input: Input, language: Language): Promise<Response>
 update(id: string, revision: number, input: Input, language: Language): Promise<Response>
 preview(id: string, revision: number, privacy: PrivacyOptions, language: Language): Promise<AIPreview>
 generate(id: string, data: { revision: number; mode: string; external_consent: boolean; privacy?: PrivacyOptions; preview_hash?: string }, language: Language): Promise<Response>
 editDraft(id: string, revision: number, events: Event[], language: Language): Promise<Response>
 review(id: string, revision: number, sources: boolean, transfers: boolean, limits: boolean, language: Language): Promise<Response>
 upload(id: string, revision: number, file: File, language: Language): Promise<Response>
 linkAttachment(id: string, attachment: string, revision: number, sources: string[], transfers: string[], language: Language): Promise<Response>
 deleteAttachment(id: string, attachment: string, revision: number, language: Language): Promise<Response>
 rotate(id: string, revision: number, language: Language): Promise<Response>
 delete(id: string, revision: number, language: Language): Promise<void>
 policy(language: Language): Promise<SafetyPolicy>
 history(id: string, language: Language): Promise<History[]>
 attachment(id: string, attachment: string, language: Language): Promise<Blob>
 export(id: string, format: 'zip' | 'txt' | 'pdf', language: Language): Promise<Blob>
}
