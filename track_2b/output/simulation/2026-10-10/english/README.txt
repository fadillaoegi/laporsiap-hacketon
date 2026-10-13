LAPORSIAP — ENGLISH SIMULATION RECORDING
10 October 2026

PLAY THE WALKTHROUGH
LaporSiap_Success_1.mp4 — 1 minute 53 seconds, 1920 x 1080, 30 fps.
The app UI, source conversation, transfer notes, draft, reference PDF,
recording captions and exported report contents are in English.
The screen recording has captions and a cursor marker, with no voice-over.

THE CASE
A fictional online task offer promises IDR 50,000 per day for liking posts.
The administrator requests an IDR 180,000 deposit, followed by an additional
IDR 420,000 deposit as a condition for withdrawing the balance.
Demo User reports the first transfer as successful and the second as pending.
The known successful-transfer total is IDR 180,000. The pending IDR 420,000
is excluded. Transfer dates remain unknown; message times are not used as
transfer dates. All names, account identifiers and transactions are fictional.

WHAT THE RECORDING SHOWS
1. Import the English TXT conversation and enter the two transfer records.
2. Check transfer statuses and save the case materials.
3. Upload an English fictional reference PDF.
4. Create a manual draft from the original English message quotes.
5. Edit all five entries while comparing them with their sources.
6. Review the sources, transfers and unknown details; confirm the revision.
7. Download the ZIP package and chronology TXT.
8. Display the revision history.
All eight steps passed against the real local API and PostgreSQL database.
The recording uses a manual draft; no live Apertus call was made.

REUSE THE DATA
synthetic-conversation.txt — the five English source messages.
synthetic-case.json — English case title, transfer fields, notes and descriptions.
synthetic-reference.pdf — fictional supporting attachment, not a bank receipt.
LaporSiap-chronology.txt — reviewed English report exported by the app.
LaporSiap-simulation.zip — complete package exported by the app.
LaporSiap_Success_1.json — step status, language, revision and recording metadata.
01...05 PNG — screenshots of the workflow.
preview-english.jpg — six stills from the successful recording.

TO FOLLOW ALONG
Open http://127.0.0.1:5175 and select English. Enter the case title from
synthetic-case.json, choose the TXT header format, and import
synthetic-conversation.txt. Check "This is synthetic demo data".
Add the two transfer records using synthetic-case.json, leaving dates blank.
Follow the video through saving, attachment upload, manual draft, editing,
source review, confirmation and export.

FAILED ATTEMPT PRESERVED
LaporSiap_Failed_1.mp4 is the first English recording. It stopped immediately
when the recorder timed out while selecting the transfer status in the form.
The script had required an exact label match, but the nested select's options
were included in its label text. The selector was corrected before starting
a separate successful recording. The original failed recording is preserved.
LaporSiap_Failed_1.json contains its stop point and error.
Original WebM recordings are also available in this project folder.

VERIFICATION
Export ZIP integrity passed; its source messages and evidence quotes matched.
English draft descriptions and every transfer field matched the supplied data.
The confirmed snapshot is revision 5 and contains the reference PDF.
Backend PostgreSQL/race tests and three bilingual browser tests passed after
localizing the manual draft's "Message quote" prefix.

RECORD AGAIN FROM THE PROJECT ROOT
node scripts/demo/record-simulation.mjs --english
Requires local Chrome, the app services, Playwright's FFmpeg component,
system FFmpeg, Node and Python. Recordings stop on a failure and are saved
with numbered Success/Failed names. Recording captions do not modify the app.
