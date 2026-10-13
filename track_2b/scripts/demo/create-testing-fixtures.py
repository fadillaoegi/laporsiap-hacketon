"""Original fictional English input files for the paired browser evaluation."""
import json
import textwrap
from pathlib import Path
from PIL import Image, ImageDraw, ImageFont
from reportlab.pdfgen import canvas
from reportlab.lib.colors import HexColor

root = Path(__file__).resolve().parents[2]
out = root / 'output/testing/2026-10-10/inputs'
out.mkdir(parents=True, exist_ok=True)
if any(out.iterdir()):
    raise SystemExit('Existing input artifacts preserved. Use the saved files or choose a new output folder.')

def transfer(amount, status):
    return dict(from_bank='Bank Aruna', from_account='****7392',
                to_bank='Bank Lestari', to_account='****4821', amount=amount,
                occurred_at='', status=status, note='Fictional user-entered record; transfer time unknown.')

cases = [
    dict(id='S01', title='Online task offer and deposit requests',
         texts=[('Nadia', 'Hi Raka, we offer tasks liking posts and promise a commission of IDR 50,000 per day.'),
                ('Nadia', 'To continue the tasks, please pay a deposit of IDR 180,000 to ****4821.'),
                ('Raka', 'I have just sent IDR 180,000. I have not recorded the transfer time yet.'),
                ('Nadia', 'Now add a deposit of IDR 420,000 so your balance can be withdrawn.'),
                ('Raka', 'The additional deposit is still pending. I have not seen it in my bank statement.')],
         descriptions=['Nadia offers tasks liking posts and promises a commission of IDR 50,000 per day.',
                       'Nadia requests a deposit of IDR 180,000 to ****4821 to continue the tasks.',
                       'Raka says he sent IDR 180,000 and has not recorded the transfer time.',
                       'Nadia requests an additional deposit of IDR 420,000 so the balance can be withdrawn.',
                       'Raka says the additional deposit is pending and has not appeared in his bank statement.'],
         transfers=[transfer(180000, 'successful'), transfer(420000, 'pending')],
         expected_total=180000, expected_unknown_success=0),
    dict(id='S02', title='Failed job registration deposit',
         texts=[('Nadia', 'Please pay a registration deposit of IDR 275,000 to ****4821.'),
                ('Raka', 'My IDR 275,000 transfer failed. The bank application shows failed.'),
                ('Nadia', 'Please try again after checking your balance.')],
         descriptions=['Nadia requests a registration deposit of IDR 275,000 to ****4821.',
                       'Raka says the IDR 275,000 transfer failed and the bank application shows failed.',
                       'Nadia asks Raka to try again after checking his balance.'],
         transfers=[transfer(275000, 'failed')], expected_total=None, expected_unknown_success=0),
    dict(id='S03', title='Two sent deposits with one unknown amount',
         texts=[('Raka', 'The first IDR 225,000 transfer was successful, but I have not recorded its transfer time.'),
                ('Raka', 'My second transfer was successful too. I cannot find its amount or transfer time yet.'),
                ('Raka', 'The next IDR 410,000 transfer is still pending. It is not a completed transfer.')],
         descriptions=['Raka says the first IDR 225,000 transfer was successful but its transfer time has not been recorded.',
                       'Raka says the second transfer was successful but its amount and transfer time are unknown.',
                       'Raka says the IDR 410,000 transfer is pending and has not completed.'],
         transfers=[transfer(225000, 'successful'), transfer(None, 'successful'), transfer(410000, 'pending')],
         expected_total=225000, expected_unknown_success=1),
]
font_path = '/System/Library/Fonts/Supplemental/Arial.ttf'
font = ImageFont.truetype(font_path, 26)
heading = ImageFont.truetype(font_path, 32)
for case in cases:
    case['raw_chat'] = '\n'.join(f'[10/10/2026, 09:{i:02d}:00] {sender}: {text}'
                                 for i, (sender, text) in enumerate(case.pop('texts'))) + '\n'
    (out / f"{case['id']}-conversation.txt").write_text(case['raw_chat'])
    case['attachments'] = ['fictional-transfer-reference.pdf', f"{case['id']}-conversation-reference.png", f"{case['id']}-transfer-notes.jpg"]
    for kind, suffix in [('conversation', 'png'), ('transfer', 'jpg')]:
        im = Image.new('RGB', (1280, 1000), '#f5f7f3')
        d = ImageDraw.Draw(im)
        d.rectangle((0, 0, 1280, 115), fill='#123c34')
        d.text((45, 25), 'SIMULATED DATA - NOT REAL EVIDENCE', font=heading, fill='white')
        lines = [case['id'] + ' - ' + case['title'], '', 'Fictional identities and masked accounts.',
                 'No real bank receipt, victim or transaction.', '']
        if kind == 'conversation':
            lines += case['raw_chat'].splitlines()
        else:
            lines += ['Bank Aruna ****7392 -> Bank Lestari ****4821', '']
            lines += [f"Transfer {i+1}: IDR {t['amount'] if t['amount'] is not None else 'UNKNOWN'} / {t['status'].upper()} / time UNKNOWN" for i, t in enumerate(case['transfers'])]
            lines += ['', 'Statuses and amounts are user-entered test inputs.', 'Attachments are reference files only; no OCR.']
        y = 145
        for line in lines:
            for part in textwrap.wrap(line, 76) or ['']:
                d.text((45, y), part, font=font, fill='#18342d')
                y += 37
        assert y < 960, 'Image content exceeds page'
        name = f"{case['id']}-{'conversation-reference' if kind == 'conversation' else 'transfer-notes'}.{suffix}"
        im.save(out / name, quality=95)
(out / 'cases.json').write_text(json.dumps(cases, indent=2))

c = canvas.Canvas(str(out / 'fictional-transfer-reference.pdf'), pagesize=(595, 842))
c.setTitle('LaporSiap fictional testing reference')
c.setFillColor(HexColor('#123c34')); c.rect(0, 730, 595, 112, fill=1, stroke=0)
c.setFillColor(HexColor('#ffffff')); c.setFont('Helvetica-Bold', 18)
c.drawString(40, 795, 'LaporSiap testing reference')
c.setFont('Helvetica', 12); c.drawString(40, 770, 'SIMULATED DATA - NOT A REAL BANK RECEIPT')
y = 695
c.setFillColor(HexColor('#18342d'))
for line in ['All identities, accounts, conversations and transactions are fictional.',
             'Sending: Bank Aruna ****7392. Receiving: Bank Lestari ****4821.',
             'Transfer times are unknown. Chat timestamps are not transfer times.', '']:
    c.setFont('Helvetica', 10); c.drawString(40, y, line); y -= 22
for case in cases:
    c.setFont('Helvetica-Bold', 12); c.drawString(40, y, case['id'] + ' - ' + case['title']); y -= 26
    c.setFont('Helvetica', 11)
    for i, t in enumerate(case['transfers']):
        amount = f"IDR {t['amount']:,}" if t['amount'] is not None else 'Amount UNKNOWN'
        c.drawString(48, y, f"Transfer {i+1}: {amount} / {t['status'].upper()} / time UNKNOWN"); y -= 24
    total = f"IDR {case['expected_total']:,}" if case['expected_total'] is not None else 'UNKNOWN (not zero)'
    c.drawString(48, y, 'Expected known successful total: ' + total); y -= 23
    if case['expected_unknown_success']:
        c.drawString(48, y, 'PARTIAL TOTAL: one successful transfer has an unknown amount.'); y -= 23
    y -= 20
c.setFont('Helvetica', 10); c.drawString(40, 85, 'Reference only. User-selected links do not verify a transfer.')
c.drawString(40, 65, 'Attachments are not sent to Apertus and are not read using OCR.')
c.save()
print(out)
