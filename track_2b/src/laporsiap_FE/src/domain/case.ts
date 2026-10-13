import type { Case, Input } from './models'
export const emptyInput = (): Input => ({ title: '', synthetic: false, raw_chat: '', parser: 'plain', transfers: [] })
export const toInput = (c: Case): Input => ({ title: c.title, synthetic: c.synthetic, raw_chat: c.raw_chat, parser: c.parser, transfers: c.transfers })
export function validateInput(input: Input) {
 if (!input.title.trim() || !input.raw_chat.trim()) throw new Error('Judul dan percakapan wajib diisi.')
 if (input.transfers.some(t => t.amount !== null && (!Number.isSafeInteger(t.amount) || t.amount < 0))) throw new Error('Nominal harus rupiah bulat nonnegatif atau kosong.')
}
