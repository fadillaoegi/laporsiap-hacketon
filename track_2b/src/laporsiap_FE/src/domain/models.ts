export type Message = { id: string; text: string; sender: string; sent_at: string }
export type Transfer = { id: string; from_bank: string; from_account: string; to_bank: string; to_account: string; amount: number | null; occurred_at: string; status: 'successful' | 'pending' | 'failed' | 'unknown'; note: string }
export type Event = { source_id: string; description: string; evidence: string }
export type Attachment = { id: string; filename: string; content_type: string; size: number; sha256: string; source_ids?: string[]; transfer_ids?: string[] }
export type Draft = { origin: string; model: string; language?: 'en' | 'id'; generated_at?: string; duration_ms?: number; usage?: { input_tokens: number; output_tokens: number }; prompt_version?: string; prompt_sha256?: string; privacy_mode?: string; redactions?: number; events: Event[]; checks: { event_index: number; source_id: string; flags: string[] }[] }
export type Case = { id: string; title: string; synthetic: boolean; raw_chat: string; parser: string; messages: Message[]; transfers: Transfer[]; draft: Draft | null; attachments: Attachment[]; revision: number; reviewed_revision: number | null; reviewed_at: string; created_at: string; updated_at: string; expires_at: string; access_expires_at: string }
export type Input = Pick<Case, 'title' | 'synthetic' | 'raw_chat' | 'parser' | 'transfers'>
export type Summary = { successful_known_total: number | null; successful_count: number; successful_unknown_amount: number; status_counts: Record<string, number>; warnings: string[] }
export type Response = { case: Case; summary: Summary; token?: string }
export type Access = { id: string; token: string; title: string }
export type History = { revision: number; action: string; created_at: string; document: Case }
export type SafetyPolicy = { retention_days: number; access_hours: number; ai_case_daily: number; ai_global_daily: number; ai_concurrent: number; provider_label: string }
export type AIPreview = { revision: number; language: 'en' | 'id'; mode: 'redacted' | 'original'; sources: { id: string; text: string }[]; redactions: number; hash: string }

export type Language = 'en' | 'id'
export type CaseAccess = Pick<Access, 'id' | 'title'>
export type CaseResponse = Omit<Response, 'token'>
export type PrivacyOptions = { mode: 'redacted' | 'original'; terms: string[] }
