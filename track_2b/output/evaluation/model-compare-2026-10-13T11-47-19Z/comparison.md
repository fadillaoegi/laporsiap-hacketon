# Apertus 8B vs 70B — model comparison (2026-10-13)

Same frozen fixture (`bilingual-safety-v1`, 12 fictional cases, 6 scenarios × EN/ID),
same prompt (`bilingual-v6-time-preserved`), temperature 0, 40 s timeout, one attempt
per case per model, no retries or fallbacks. Called directly through the app's
provider path (`PrepareAIInput → Generate → RestoreEvidence`) via Public AI;
models interleaved per case. 24 requests total.

| | Apertus v1.5 8B | Apertus v1.5 70B |
| --- | --- | --- |
| Provider answered (HTTP 200) | 12/12 | 6/12 (6 × HTTP 400) |
| Contract passed (schema, source binding, timing) | **12/12** | 6/12 |
| Semantic review — clean | 7 | 5 |
| Semantic review — minor | 2 | 1 |
| Semantic review — material | 3 | 0 |
| Mean request time | 5.9 s | 2.1 s (answered only) |
| Est. cost per answered case (catalog price) | $0.00011 | $0.00105 |

Material 8B issues: lost "previously said" attribution in a corrected amount
(conflicting-status EN+ID); invented a registration step and mixed two messages
(privacy-multiline-ID). These are caught only by human review, which LaporSiap
requires before export.

The 70B HTTP 400 responses arrived in ~1.3–2.0 s with no body recorded (the app
deliberately discards provider error bodies). Cause is undiagnosed — likely a
provider-side rejection, not a model-quality signal. On 10 Oct (earlier prompt)
the same 70B endpoint answered 4 of these 6 scenarios (source-instruction EN/ID
failed then too). 70B semantic quality is therefore compared on 6 cases only.

Limitations: 12-case development fixture, single assistant reviewer (owner re-check
pending), model identity is API metadata, timing reflects a shared hosted endpoint
at one moment, catalog prices are not an invoice, failed-request usage unknown.
