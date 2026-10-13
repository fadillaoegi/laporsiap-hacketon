import type { AIPreview } from '../../domain/models'
import { useLanguage } from '../i18n'

export function SourcePreview({ preview, masked, terms, pending, error, locked, onMask, onTerms, onEdit }: { preview: AIPreview | null; masked: boolean; terms: string; pending: boolean; error: string; locked: boolean; onMask: (value: boolean) => void; onTerms: (value: string) => void; onEdit: () => void }) {
  const { tr, language } = useLanguage()
  return <section className="source-preview">
    <div className="preview-heading"><strong>{tr('Periksa teks sebelum dikirim')}</strong><span className="small-tag">{language === 'en' ? 'English' : 'Indonesia'}</span></div>
    <p>{tr('Draft baru memakai bahasa pilihan Anda. Kutipan sumber tetap dalam bahasa aslinya.')}</p>
    <label className="checkbox"><input type="checkbox" disabled={locked} checked={masked} onChange={e => onMask(e.target.checked)}/><span>{tr('Samarkan informasi pribadi sebelum dikirim ke AI')}</span></label>
    <p>{tr('Mendeteksi nama pengirim, email, telepon Indonesia, serta nomor rekening/identitas berlabel. Deteksi dapat melewatkan data atau menyamarkan bagian yang tidak Anda maksud; periksa preview.')}</p>
    {masked && <label>{tr('Istilah tambahan untuk disamarkan')}<textarea rows={2} disabled={locked} value={terms} onChange={e => onTerms(e.target.value)} placeholder={tr('Satu nama atau istilah per baris; maksimal 50.')}/></label>}
    {pending && <p role="status">{tr('Menyiapkan preview pribadi…')}</p>}
    {error && <p className="inline-warning">{tr(error)}</p>}
    {preview && <p>{tr(preview.mode === 'redacted' ? 'Informasi pribadi disamarkan:' : 'Teks asli akan dikirim:')}{' '}{preview.redactions}{' '}{tr('bagian disamarkan.')}</p>}
    <details><summary>{tr('Lihat seluruh pesan yang akan dikirim')}</summary>
      <div className="preview-messages">{preview?.sources.map(m => <div key={m.id}><strong>{m.id}</strong><blockquote>{m.text}</blockquote></div>)}</div>
    </details>
    <p>{tr('Hanya ID dan teks preview di atas dikirim, bersama instruksi penyusunan. Sumber asli disimpan untuk review dan ekspor. File lampiran dan isian transfer tidak dikirim.')}</p>
    <button className="text-button" disabled={locked} onClick={onEdit}>{tr('Koreksi teks sebelum mengirim →')}</button>
  </section>
}
