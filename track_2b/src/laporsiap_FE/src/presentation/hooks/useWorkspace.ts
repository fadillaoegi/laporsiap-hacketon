import { useEffect, useRef } from 'react'
import type { Input, CaseResponse, PrivacyOptions, Transfer } from '../../domain/models'
import { toInput } from '../../domain/case'
import { useLanguage } from '../i18n'
import { useWorkflow } from '../services'
import { useWorkspaceField, useAppDispatch, useAppSelector } from '../state/hooks'
import { workspaceActions } from '../state/workspaceSlice'
export function useWorkspace() {
  const { language, setLanguage, tr } = useLanguage()
  const [accessList, setAccessList] = useWorkspaceField('accessList')
  const [access, setAccess] = useWorkspaceField('access')
  const [c, setCase] = useWorkspaceField('c')
  const [summary, setSummary] = useWorkspaceField('summary')
  const [input, setInput] = useWorkspaceField('input')
  const [events, setEvents] = useWorkspaceField('events')
  const [step, setStep] = useWorkspaceField('step')
  const [busy, setBusy] = useWorkspaceField('busy')
  const [error, setError] = useWorkspaceField('error')
  const [notice, setNotice] = useWorkspaceField('notice')
  const [consent, setConsent] = useWorkspaceField('consent')
  const [maskPrivate, setMaskPrivate] = useWorkspaceField('maskPrivate')
  const [privateTerms, setPrivateTerms] = useWorkspaceField('privateTerms')
  const [aiPreview, setAIPreview] = useWorkspaceField('aiPreview')
  const [previewPending, setPreviewPending] = useWorkspaceField('previewPending')
  const [previewError, setPreviewError] = useWorkspaceField('previewError')
  const [policy, setPolicy] = useWorkspaceField('policy')
  const [reviewedSources, setReviewedSources] = useWorkspaceField('reviewedSources')
  const [reviewedTransfers, setReviewedTransfers] = useWorkspaceField('reviewedTransfers')
  const [acceptedLimits, setAcceptedLimits] = useWorkspaceField('acceptedLimits')
  const [editingAttachment, setEditingAttachment] = useWorkspaceField('editingAttachment')
  const [history, setHistory] = useWorkspaceField('history')
  const initialLoad = useRef(false)
  const errorRef = useRef<HTMLDivElement>(null)
  const inputDirty = c !== null && JSON.stringify(input) !== JSON.stringify(toInput(c))
  const draftDirty = c?.draft !== null && c?.draft !== undefined && JSON.stringify(events) !== JSON.stringify(c.draft.events)
  const dirty = inputDirty || draftDirty || editingAttachment !== null
  const aiInputTooLarge = !!c && (c.messages.length > 20 || new TextEncoder().encode(c.raw_chat).byteLength > 12000)
  const confirmed = !!c && c.reviewed_revision === c.revision && !dirty
  const completedSteps = [!!c && !inputDirty, !!c?.draft && !inputDirty && !draftDirty, confirmed]
  const allReviewed = !!c && reviewedSources.length === c.messages.length && reviewedTransfers && acceptedLimits && !dirty

  const dispatch = useAppDispatch()
  const workflow = useWorkflow()
  const accessEpoch = useAppSelector(state => state.workspace.accessEpoch)
  function apply(data: CaseResponse) { dispatch(workspaceActions.applyCase(data)) }
  async function run(label: string, action: () => Promise<void>) { setBusy(label); setError(''); setNotice(''); try { await action() } catch (e) { setError(e instanceof Error ? e.message : 'Terjadi kesalahan.') } finally { setBusy('') } }
  useEffect(() => { if (initialLoad.current) return; initialLoad.current = true; if (access) void run('Memuat kasus…', async () => { const data = await workflow.load(access.id, language); apply(data); setStep(data.case.reviewed_revision === data.case.revision ? 3 : data.case.draft ? 2 : 0) }) }, [])
  useEffect(() => { void workflow.policy(language).then(setPolicy).catch(() => {}) }, [])
  useEffect(() => {
    setConsent(false); setAIPreview(null); setPreviewError('')
    if (!c || !access || step !== 1 || aiInputTooLarge) { setPreviewPending(false); return }
    let cancelled = false; setPreviewPending(true)
    const timer = window.setTimeout(() => {
      void workflow.preview(c, privacy(), language)
        .then(data => { if (!cancelled) setAIPreview(data) })
        .catch(e => { if (!cancelled) setPreviewError(e instanceof Error ? e.message : 'Terjadi kesalahan.') })
        .finally(() => { if (!cancelled) setPreviewPending(false) })
    }, 300)
    return () => { cancelled = true; window.clearTimeout(timer) }
  }, [c?.id, c?.revision, accessEpoch, step, language, maskPrivate, privateTerms, aiInputTooLarge])
  useEffect(() => {
    if (!error) return
    errorRef.current?.scrollIntoView({ block: 'center', behavior: 'auto' })
    errorRef.current?.focus({ preventScroll: true })
  }, [error])
  useEffect(() => { const handler = (e: BeforeUnloadEvent) => { if (dirty) e.preventDefault() }; window.addEventListener('beforeunload', handler); return () => window.removeEventListener('beforeunload', handler) }, [dirty])

  useEffect(() => { setConsent(false) }, [language, c?.revision])

  function updateInput(patch: Partial<Input>) { setInput(old => ({ ...old, ...patch })); setConsent(false); setReviewedSources([]); setReviewedTransfers(false); setAcceptedLimits(false) }
  function transferChange(index: number, patch: Partial<Transfer>) { updateInput({ transfers: input.transfers.map((t, i) => i === index ? { ...t, ...patch } : t) }) }
  const privacy = (): PrivacyOptions => ({ mode: maskPrivate ? 'redacted' : 'original', terms: privateTerms.split('\n').map(s => s.trim()).filter(Boolean) })
  async function saveInput() { const data = await workflow.save(c, input, language); setAccess({ id: data.case.id, title: data.case.title }); setAccessList(workflow.records()); apply(data); setNotice('Bahan kasus tersimpan.'); setStep(1) }
  async function generate(mode: string) { if (!c || !access) return; const externalConsent = mode === 'apertus' && consent; setConsent(false); apply(await workflow.generate(c, mode, externalConsent, privacy(), aiPreview?.hash, language)); setStep(2); setNotice(mode === 'manual' ? 'Draft manual memuat kutipan pesan. Periksa dan sunting tiap entri.' : 'Draft Apertus tersimpan. Seluruh entri tetap perlu review makna.') }
  async function saveDraft() { if (!c || !access) return; apply(await workflow.edit(c, events, language)); setNotice('Koreksi tersimpan. Periksa revisi terbaru sebelum konfirmasi.') }
  async function confirm() { if (!c || !access) return; apply(await workflow.review(c, allReviewed, reviewedTransfers, acceptedLimits, language)); setStep(3); setNotice('Review revisi terbaru telah dikonfirmasi. Bahan siap diekspor.') }
  async function upload(file: File) { if (!c || !access) return; apply(await workflow.upload(c, file, language)); setNotice('Lampiran tersimpan. Konfirmasi review lama dibatalkan.') }
  function newCase() { dispatch(workspaceActions.newCase()) }
  async function rotateAccess() { if (!c || !access) return; apply(await workflow.rotate(c, language)); setAccessList(workflow.records()); dispatch(workspaceActions.accessRotated()); setNotice('Akses lama dicabut. Akses sesi ini telah diperbarui.') }
  async function deleteCase() { if (!c || !access) return; await workflow.delete(c, language); setAccessList(workflow.records()); newCase(); setNotice('Kasus, riwayat, dan lampiran telah dihapus dari database.') }
  async function selectCase(id: string) { const data = await workflow.load(id, language); apply(data); setAccess({ id, title: data.case.title }); setStep(0); setConsent(false) }
  async function reloadCase() { if (!access) return; apply(await workflow.load(access.id, language)); setNotice('Revisi terbaru dimuat.') }
  async function deleteAttachment(id: string) { if (!c || !access) return; apply(await workflow.removeAttachment(c, id, language)); setNotice('Lampiran dihapus. Review perlu dikonfirmasi ulang.') }
  async function linkAttachment(sourceIDs: string[], transferIDs: string[]) { if (!c || !access || !editingAttachment) return; apply(await workflow.link(c, editingAttachment.id, sourceIDs, transferIDs, language)); setNotice('Hubungan bukti tersimpan. Konfirmasi review lama dibatalkan.') }
  async function loadHistory() { if (!c || !access) return; setHistory(await workflow.history(c, language)) }
  async function downloadAttachment(a: { id: string; filename: string }) { if (c && access) await workflow.downloadAttachment(c, a, language) }
  async function exportCase(format: 'zip' | 'pdf' | 'txt') { if (c && access) await workflow.export(c, format, language) }
  async function importTXT(file: File) { updateInput({ raw_chat: await workflow.readConversation(file) }) }
  const disabled = !!busy
  const locked = disabled || dirty
  return { language, setLanguage, tr, accessList, access, c, summary, input, events, step, busy, error, notice, consent, maskPrivate, privateTerms, aiPreview, previewPending, previewError, policy, reviewedSources, reviewedTransfers, acceptedLimits, editingAttachment, history, errorRef, dirty, inputDirty, draftDirty, aiInputTooLarge, confirmed, completedSteps, allReviewed, disabled, locked, setInput, setEvents, setStep, setError, setNotice, setConsent, setMaskPrivate, setPrivateTerms, setAIPreview, setReviewedSources, setReviewedTransfers, setAcceptedLimits, setEditingAttachment, run, updateInput, transferChange, saveInput, generate, saveDraft, confirm, upload, newCase, rotateAccess, deleteCase, selectCase, reloadCase, deleteAttachment, linkAttachment, loadHistory, downloadAttachment, exportCase, importTXT }
}
