import { createSlice } from '@reduxjs/toolkit'
import type { PayloadAction } from '@reduxjs/toolkit'
import { emptyInput, toInput } from '../../domain/case'
import type { AIPreview, Case, CaseAccess, CaseResponse, Event, History, Input, SafetyPolicy, Summary } from '../../domain/models'
export type WorkspaceState = {
 accessList: CaseAccess[]; access: CaseAccess | null; accessEpoch: number
 c: Case | null; summary: Summary | null; input: Input; events: Event[]
 step: number; busy: string; error: string; notice: string
 consent: boolean; maskPrivate: boolean; privateTerms: string
 aiPreview: AIPreview | null; previewPending: boolean; previewError: string; policy: SafetyPolicy | null
 reviewedSources: string[]; reviewedTransfers: boolean; acceptedLimits: boolean
 editingAttachment: Case['attachments'][number] | null; history: History[] | null
 linkedSources: string[]; linkedTransfers: string[]; deleteConfirmation: string
}
export const initialWorkspace = (records: CaseAccess[] = []): WorkspaceState => ({
 accessList: records, access: records.at(-1) || null, accessEpoch: 0, c: null, summary: null, input: emptyInput(), events: [], step: 0, busy: '', error: '', notice: '', consent: false, maskPrivate: true, privateTerms: '', aiPreview: null, previewPending: false, previewError: '', policy: null, reviewedSources: [], reviewedTransfers: false, acceptedLimits: false, editingAttachment: null, history: null, linkedSources: [], linkedTransfers: [], deleteConfirmation: '',
})
export const workspaceSlice = createSlice({ name: 'workspace', initialState: initialWorkspace(), reducers: {
 patch(state, action: PayloadAction<Partial<WorkspaceState>>) {
  const attachment = action.payload.editingAttachment
  if (attachment) { state.linkedSources = [...(attachment.source_ids || [])]; state.linkedTransfers = [...(attachment.transfer_ids || [])] }
  Object.assign(state, action.payload)
 },
 applyCase(state, action: PayloadAction<CaseResponse>) {
  const data = action.payload
  if (state.c?.id !== data.case.id) { state.maskPrivate = true; state.privateTerms = ''; state.consent = false; state.aiPreview = null; state.deleteConfirmation = '' }
  state.c = data.case; state.summary = data.summary; state.input = toInput(data.case); state.events = data.case.draft?.events || []
  state.reviewedSources = []; state.reviewedTransfers = false; state.acceptedLimits = false; state.history = null; state.editingAttachment = null
 },
 newCase(state) { const next = initialWorkspace(state.accessList); next.access = null; next.busy = state.busy; next.policy = state.policy; next.accessEpoch = state.accessEpoch; return next },
 accessRotated(state) { state.accessEpoch += 1; state.consent = false; state.aiPreview = null },
} })
export const workspaceActions = workspaceSlice.actions
