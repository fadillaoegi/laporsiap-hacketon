import { useWorkspaceField } from '../state/hooks'
import type { Case, SafetyPolicy } from '../../domain/models'
import { useLanguage } from '../i18n'

export function CaseControls({ c, policy, busy, dirty, onRotate, onDelete }: { c: Case | null; policy: SafetyPolicy | null; busy: boolean; dirty: boolean; onRotate: () => void; onDelete: () => void }) {
  const { language, tr } = useLanguage()
  const [confirmation, setConfirmation] = useWorkspaceField('deleteConfirmation')
  const date = (value: string) => value ? new Date(value).toLocaleString(language === 'en' ? 'en-US' : 'id-ID') : tr('Belum tercatat')
  return <section className="card data-controls"><details><summary>{tr('Data dan akses Anda')}</summary>
    {policy && <p>{tr('Retensi kasus:')} <strong>{policy.retention_days} {tr('hari sejak dibuat.')}</strong> {tr('Batas AI per kasus:')} <strong>{policy.ai_case_daily} {tr('percobaan per 24 jam.')}</strong></p>}
    <p>{tr('Kasus yang kedaluwarsa tidak dapat diakses dan akan dihapus dari database. Unduhan dan salinan cadangan terpisah tidak ikut dihapus.')}</p>
    {c ? <>
      <dl><dt>{tr('Kasus kedaluwarsa')}</dt><dd>{date(c.expires_at)}</dd><dt>{tr('Akses sesi kedaluwarsa')}</dt><dd>{date(c.access_expires_at)}</dd></dl>
      <p>{tr('Perbarui akses untuk mencabut akses lama. Tab lain dengan akses lama perlu ditutup; sesi ini menerima akses baru. Retensi kasus tidak diperpanjang.')}</p>
      <button className="button secondary" disabled={busy || dirty} onClick={onRotate}>{tr('Perbarui akses sesi')}</button>
      <div className="delete-case-box"><h3>{tr('Hapus kasus ini')}</h3><p>{tr('Penghapusan permanen mencakup percakapan, draft, riwayat, dan lampiran dalam database. Unduhan Anda tetap tersimpan di luar aplikasi.')}</p>
      <label>{tr('Ketik DELETE untuk menghapus')}<input value={confirmation} disabled={busy} onChange={e => setConfirmation(e.target.value)} autoComplete="off"/></label>
      <button className="button danger" disabled={busy || confirmation !== 'DELETE'} onClick={onDelete}>{tr('Hapus kasus secara permanen')}</button></div>
    </> : <p>{tr('Tutup tab akan mengakhiri akses sesi; data tersimpan sampai retensi berakhir. Anda dapat menghapus kasus saat akses masih aktif.')}</p>}
  </details></section>
}
