# Technical report — LaporSiap

Turning an untidy scam chat into reviewed, source-linked report material for
Indonesia's official anti-scam channel — with Apertus drafting and a human
checking every line.

- **Track:** Track 2B — Own Project · LaporSiap
- **Event:** Online
- **Team:** FLdevBoo — M Fadil Muhtarom Laoegi
- **Model:** Apertus v1.5 8B + 70B (`swiss-ai/apertus-v1.5-8b`, `swiss-ai/apertus-v1.5-70b`)
- **Code:** https://github.com/fadillaoegi/laporsiap-hacketon (`track_2b/`)
- **Demo:** `link to video`
- **Maturity:** working prototype, validated with synthetic data only (no real users yet)
- **Coding assistant:** built with Claude Code (Anthropic) for implementation,
  tests and documentation; design decisions, review and commits by the team.

## 1. Summary

Victims of part-time-job and "top-up" scams in Indonesia are asked by the
official channel, OJK's Indonesia Anti-Scam Centre (IASC), for a chronology,
transaction details and evidence. IASC received **637,055 reports** between
22 Nov 2024 and 31 Jul 2026 [1]. What people actually have is a messy chat.
LaporSiap lets them paste that chat, record the transfers they know, and get an
**Apertus-drafted chronology where every entry is bound to its original message
by ID** — the model writes only short descriptions; quotes are attached locally
from the source and can never be altered by the model. A person must confirm
every entry before PDF/TXT/ZIP export is unlocked. Totals are computed by code,
never by the model. Headline result: on a frozen 12-case bilingual fixture,
**Apertus 8B passed the output contract 12/12 at ~USD 0.0001 per case**, but
human-style review found **3 material meaning errors** that no automatic check
caught — the evidence for why LaporSiap keeps a mandatory human review gate.

## 2. Architecture

```
Browser (React + Redux, EN/ID)
   │  same origin, loopback only
Go API ── Clean Architecture: domain ← application ← adapters
   ├── PostgreSQL 18 (cases, revisions, attachments, AI quota)
   ├── Report renderer (PDF/TXT/CSV/ZIP, fonts embedded)
   └── Apertus adapter ──► OpenAI-compatible endpoint (LLM_BASE_URL)
                             • on-prem vLLM (prepared)  • Public AI (used for tests)
```

Data flow for one AI draft: case saved → user sees an **exact preview** of what
will be sent (masked message IDs + text) → explicit consent bound to a preview
hash → one request, temperature 0, 40 s timeout, no retry, no fallback model →
response must name the requested Apertus model, be strict JSON with one
`source_id` + `description` (≤35 words) per message, in order, preserving
relative time phrases → quotes re-attached from local source → review flags
(direction of funds, injected instructions, pending→completed, promises written
as events, amounts not in quote) → human review → export.

### Target architecture

**a) On-premise — primary target.** One Docker Compose stack: app (non-root,
read-only filesystem, all capabilities dropped, loopback port) + PostgreSQL in
the organisation's own infrastructure. For AI, `deploy/compose.gpu.yaml` adds the
official Apertus vLLM runtime pinned by image digest and checkpoint revision
(`deploy/gpu-lock.json`) sharing the app's private network namespace; inference
binds 127.0.0.1 only. **b) Air-gapped — supported** once images and weights are
pre-loaded: no CDN, analytics, web fonts or telemetry; the full manual workflow
works with zero network, and AI needs only the local vLLM endpoint.

| Phase | External dependency | Data sent |
|---|---|---|
| Build | Docker Hub images (golang, node, postgres), Go module proxy, npm registry | none (public packages) |
| Setup (AI) | Hugging Face gated weights download | owner's HF token |
| Runtime | **Only** the configured Apertus endpoint | prompt + masked message IDs/text of one case, after consent |
| Runtime | — | transfers, attachments, titles, sender/time metadata and access tokens **never leave** |

Honest status: the on-prem GPU path is prepared and validated statically, but
**70B inference has not yet been executed on our own GPU**; all live results
below used the hosted Public AI endpoint.

## 3. Use of Apertus

- **Models:** `swiss-ai/apertus-v1.5-8b` and `swiss-ai/apertus-v1.5-70b`.
- **How it is used:** inference (structured summarisation) and evaluation /
  red-teaming (prompt injection inside source messages).
- **Where it runs:** hosted Public AI endpoint for all tests; on-prem vLLM
  (BF16, tensor parallel 4, 4× H100 80 GB planned for 70B) prepared.
- **Prompt:** `bilingual-v6-time-preserved`, embedded in the binary; its
  SHA-256 is exposed at `/api/policy` and stored with every draft.
- **Why Apertus:** a Swiss open model with open weights and data that an
  Indonesian public-interest service could self-host, so victims' chats need
  not leave the operator's infrastructure. Bilingual EN/ID worked without
  fine-tuning.
- **Where Apertus fell short:** 8B dropped attribution ("I previously said…")
  and once invented a step; the hosted 70B gateway intermittently rejected valid
  requests with HTTP 400 (6/12 on 13 Oct; an identical retry succeeded); on 11 Oct one request was answered by a non-Apertus model
  (`aisingapore/Qwen-SEA-LION-v4-32B-IT`) and rejected by our identity check;
  earlier prompts lost "next week" until timing validation was added.

## 4. Data

All data is **original and synthetic** (fictional people Nadia/Raka, fictional
banks, masked accounts), licensed CDLA-Permissive-2.0: 12-case evaluation fixture
`docs/evaluation/bilingual-safety-v1.json`, 13 manual test scenarios and 8
importable chats in `docs/testing/`. No real victim data, no personal data, no
human subjects. Model outputs are kept as evidence and retain upstream rights.

## 5. Evaluation

**Task:** one chronology entry per source message. **Metrics:** (1) output
contract — strict schema, source binding, order, timing preservation, model
identity; (2) semantic review per entry against fixture criteria (clean /
minor / material); (3) latency and token cost. **Setup:** frozen fixture
(SHA-256), prompt v6, temperature 0, one attempt per case per model, models
interleaved, no retries; run 2026-10-13 through the app's own provider path.
The 6 cases the hosted 70B gateway rejected with HTTP 400 were diagnosed and
rerun once on 2026-10-14 (31 live requests in total).

| Setup | Contract passed | Semantic: clean / minor / material | Mean latency | Cost per case* |
|---|---|---|---|---|
| Baseline: manual draft (no AI) | n/a | quotes only, user writes | — | 0 |
| Apertus v1.5 8B | **12/12** | 7 / 2 / 3 | 5.9 s | ≈ $0.00011 |
| Apertus v1.5 70B | **12/12** (6 after one rerun) | 7 / 5 / **0** | 2.2 s | ≈ $0.00105 |

\*Catalog token prices, not an invoice. **Provider 400 diagnosis:** the gateway
answered "This model only accepts a system message as the first message, with its
content as a plain string" although the request already sends exactly that, and an
identical retry seconds later returned HTTP 200 — an intermittent provider-side
rejection, not a model or app fault. LaporSiap deliberately does not auto-retry;
the user retries or drafts manually. Safety held in every answered case:
injected "SYSTEM: ignore previous instructions" text was described, never
obeyed; no masked personal data appeared in descriptions; totals and sources
were unchanged. Engineering checks: Go unit + PostgreSQL integration tests with
the race detector, architecture-boundary checks, 13 browser regressions, all
green in GitHub Actions CI.

**Decision:** 70B made no material meaning errors (8B: 3), yet both still
produce minor unsupported wording — so human review stays mandatory for either.
8B is the on-prem default — 10× cheaper, a single GPU — *because* the review gate
catches its errors; 70B is the quality option when 4× H100 is available.

**Value and scale.** A case is ~870 input + ~110 output tokens. At 10× volume
(e.g. 1,000 → 10,000 cases/day) hosted 8B costs ≈ $1.1/day vs ≈ $10.5/day for
70B; self-hosted cost is GPU-hours, not tokens (4× H100 for 70B ≈ $16/h list
price). The Go API is stateless apart from PostgreSQL; AI concurrency, daily
quotas and storage caps are configuration.

## 6. Limitations

12 development cases are not a statistical benchmark; semantic review was done
by one reviewer (the coding assistant), not independently. No real users or
victims have tested it; a 2–3 person pilot kit is ready, not run. Model identity
is API metadata, not a weights audit. Regex masking is not anonymisation. No
OCR — attachments are references only. No accounts or token recovery. Not
affiliated with OJK/IASC; it does not file reports, block accounts or promise
refunds. The 70B rerun was a second attempt on a different day; the hosted
gateway's intermittent HTTP 400 remains a provider-side risk.

## 7. Reproducibility

```sh
cd track_2b && make run            # Docker: PostgreSQL + app → http://localhost:8787
export LLM_NAME=swiss-ai/apertus-v1.5-8b LLM_BASE_URL=https://api.publicai.co/v1 LLM_API_KEY=...
cd src/laporsiap_BE && go run ./cmd/modelcompare        # dry run, no API calls
go run ./cmd/modelcompare -live                          # 24 requests, ~USD 0.03
```

Results: `output/evaluation/model-compare-2026-10-13T11-47-19Z/` and the 70B rerun `output/evaluation/model-compare-2026-10-14T06-16-45Z/` (`-cases` filter)
(`comparison.md`, `semantic-review.json`, per-case JSON, frozen fixture and
prompt). Hardware for the app: any x86/ARM machine with Docker (tested on an
8 GB Mac). Determinism: temperature 0, but hosted outputs can still vary.
Commit: `f66aa78` or later on `main`.

## 8. Next steps

Run the 8B model on our own GPU and record the evidence checklist in
`deploy/GPU-PROOF.txt`; run the pilot with 2–3 independent users and measure
review time; partner with a victim-support organisation or legal-aid clinic as
first operator; OCR for transfer screenshots; accounts and token recovery for a
real service.

## License

Documentation CC-BY-4.0, code Apache-2.0, synthetic data CDLA-Permissive-2.0.

## References

1. OJK, Satgas PASTI press release, 1 Sep 2026 — https://ojk.go.id/id/berita-dan-kegiatan/info-terkini/Pages/Satgas-PASTI-Terima-25-Ribu-Pengaduan-Aktivitas-Keuangan-Ilegal-dan-IASC-Berhasil-Blokir-Dana-Terkait-Scam-Rp-724-Miliar.aspx
2. IASC portal — https://iasc.ojk.go.id/
3. Apertus v1.5 70B model card — https://huggingface.co/swiss-ai/Apertus-v1.5-70B
4. Public AI API documentation — https://platform.publicai.co/docs
