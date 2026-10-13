import type { FileSaver } from '../../application/ports'
export class BrowserFileSaver implements FileSaver {
 save(bytes: Blob, filename: string) { const url = URL.createObjectURL(bytes); const a = document.createElement('a'); a.href = url; a.download = filename; a.click(); setTimeout(() => URL.revokeObjectURL(url), 1000) }
}
