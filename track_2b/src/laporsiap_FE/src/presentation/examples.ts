import type { Input, Transfer } from '../domain/models'
import type { Language } from '../domain/models'

export const newTransfer = (): Transfer => ({ id: crypto.randomUUID(), from_bank: '', from_account: '', to_bank: '', to_account: '', amount: null, occurred_at: '', status: 'unknown', note: '' })
export const emptyInput = (): Input => ({ title: '', synthetic: false, raw_chat: '', parser: 'plain', transfers: [] })
export const demoInput = (language: Language = 'id'): Input => language === 'en' ? ({
  title: 'Online task offer and deposit requests', synthetic: true, parser: 'whatsapp_bracket_dmy',
  raw_chat: '[08/10/2026, 09:10] Nadia: Hi Raka, we offer liking tasks with a promised commission of IDR 50,000 per day.\n[08/10/2026, 09:14] Nadia: To continue, please send a deposit of IDR 180,000 to ****4821.\n[08/10/2026, 09:28] Raka: I have just sent IDR 180,000. I have not recorded the transfer time.\n[08/10/2026, 09:35] Nadia: Now add a deposit of IDR 420,000 to unlock withdrawal.\n[08/10/2026, 09:40] Raka: The additional deposit is still pending. I have not checked my bank statement.',
  transfers: [
    { ...newTransfer(), from_bank: 'Bank Aruna', from_account: '****7392', to_bank: 'Bank Lestari', to_account: '****4821', amount: 180000, status: 'successful', note: 'Transfer date is unknown; do not use the message time.' },
    { ...newTransfer(), from_bank: 'Bank Aruna', from_account: '****7392', to_bank: 'Bank Lestari', to_account: '****4821', amount: 420000, status: 'pending', note: 'Not confirmed successful.' },
  ],
}) : ({
  title: 'Tawaran tugas daring dan permintaan top-up', synthetic: true, parser: 'whatsapp_bracket_dmy',
  raw_chat: '[08/10/2026, 09:10] Nadia: Halo Raka, kami menawarkan tugas memberi like dengan janji komisi Rp50.000 per hari.\n[08/10/2026, 09:14] Nadia: Untuk lanjut tugas, silakan setor deposit Rp180.000 ke ****4821.\n[08/10/2026, 09:28] Raka: Saya baru mengirim Rp180.000. Waktu transfernya belum saya catat.\n[08/10/2026, 09:35] Nadia: Sekarang tambah deposit Rp420.000 agar saldo bisa ditarik.\n[08/10/2026, 09:40] Raka: Status setoran tambahan masih dalam proses. Saya belum melihat mutasi rekening.',
  transfers: [
    { ...newTransfer(), from_bank: 'Bank Aruna', from_account: '****7392', to_bank: 'Bank Lestari', to_account: '****4821', amount: 180000, status: 'successful', note: 'Tanggal transfer belum dicatat; jangan memakai waktu chat.' },
    { ...newTransfer(), from_bank: 'Bank Aruna', from_account: '****7392', to_bank: 'Bank Lestari', to_account: '****4821', amount: 420000, status: 'pending', note: 'Belum dipastikan berhasil.' },
  ],
})
