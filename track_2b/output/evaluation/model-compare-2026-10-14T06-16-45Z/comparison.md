# Apertus 8B vs 70B — combined comparison (13–14 Oct 2026)

Same frozen fixture (`bilingual-safety-v1`, 12 fictional cases), same prompt
(`bilingual-v6-time-preserved`), temperature 0, one attempt per request, no retries,
provider path `PrepareAIInput → Generate → RestoreEvidence` via Public AI.

13 Oct (`model-compare-2026-10-13T11-47-19Z`): 24 requests; 70B was rejected with
HTTP 400 on 6 cases. 14 Oct (this folder): 1 diagnostic request + one rerun of those
6 70B cases only (`-cases` filter). Total live requests: 31.

| | Apertus v1.5 8B | Apertus v1.5 70B |
| --- | --- | --- |
| Contract passed | **12/12** first attempt | **12/12** (6 on 13 Oct + 6 rerun on 14 Oct after provider 400) |
| Semantic review — clean / minor / material | 7 / 2 / **3** | 7 / 5 / **0** |
| Mean request time | 5.9 s | 2.2 s |
| Tokens in / out (12 cases) | 10383 / 1350 | 10383 / 1347 |
| Est. cost per case (catalog price) | ≈ $0.0001 | ≈ $0.001 (10×) |

Provider 400 diagnosis: the gateway answered "This model only accepts a system message
as the first message, with its content as a plain string" although the request already
sends exactly that; an identical retry seconds later succeeded. It is an intermittent
provider-side rejection, not a model or app fault. LaporSiap keeps no automatic retry
(the user can retry or write a manual draft).

Reading: 70B made no material meaning errors (8B made 3), but both still produce minor
unsupported wording — human review stays mandatory. 8B remains the low-cost single-GPU
default; 70B is the quality option when 4× H100 is available.

Limitations: 12-case development fixture, single assistant reviewer (owner re-check
pending), the 70B rerun is a second attempt on a different day, timing reflects a shared
hosted endpoint, catalog prices are not an invoice.
