import { useWorkspaceField } from '../state/hooks'
import type { Attachment, Case } from '../../domain/models'
import { useLanguage } from '../i18n'

export function AttachmentLinkEditor({ c, attachment, disabled, onSave, onCancel }: {
  c: Case; attachment: Attachment; disabled: boolean;
  onSave: (sourceIDs: string[], transferIDs: string[]) => void; onCancel: () => void;
}) {
  const { tr } = useLanguage()
  const [sources, setSources] = useWorkspaceField('linkedSources')
  const [transfers, setTransfers] = useWorkspaceField('linkedTransfers')
  const toggle = (values: string[], id: string, checked: boolean) => checked ? [...values, id] : values.filter(v => v !== id)
  return <section className="card link-editor" aria-label={tr('Hubungkan lampiran')}>
    <div className="eyebrow">{tr('HUBUNGAN BUKTI')}</div><h2>{attachment.filename}</h2>
    <p className="field-help">{tr('Pilih bagian yang didukung lampiran menurut Anda. Hubungan ini tidak memverifikasi isi atau keaslian file.')}</p>
    <fieldset disabled={disabled}><legend>{tr('Pesan sumber')}</legend>
      {c.messages.map(m => <label className="checkbox link-choice" key={m.id}><input type="checkbox" checked={sources.includes(m.id)} onChange={e => setSources(toggle(sources, m.id, e.target.checked))}/><span><strong>{m.id}</strong> {m.text}</span></label>)}
    </fieldset>
    <fieldset disabled={disabled}><legend>{tr('Catatan transfer')}</legend>
      {c.transfers.length === 0 && <p>{tr('Tidak ada transfer yang diisi pengguna.')}</p>}
      {c.transfers.map((t, i) => <label className="checkbox link-choice" key={t.id}><input type="checkbox" checked={transfers.includes(t.id)} onChange={e => setTransfers(toggle(transfers, t.id, e.target.checked))}/><span><strong>Transfer {String(i + 1).padStart(2, '0')}</strong> · {t.to_bank || tr('Bank belum diketahui')} / {t.to_account || tr('Rekening belum diketahui')}</span></label>)}
    </fieldset>
    <div className="link-actions"><button className="button primary" disabled={disabled} onClick={() => onSave(sources, transfers)}>{tr('Simpan hubungan bukti')}</button><button className="button secondary" disabled={disabled} onClick={onCancel}>{tr('Batal')}</button></div>
  </section>
}
