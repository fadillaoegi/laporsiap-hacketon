"""Create local, inspectable review archives without secrets or model weights."""
from datetime import datetime, timezone
from pathlib import Path
import hashlib
import json
import re
import shlex
import zipfile

repo = Path(__file__).resolve().parents[1]  # track_2b project root
top = repo.parent  # git root: LICENSE, NOTICE, .github
out = repo / "output" / "submission-review" / datetime.now(timezone.utc).strftime("%Y-%m-%dT%H-%M-%S.%fZ")
out.mkdir(parents=True, exist_ok=False)
root_files = [top / n for n in ("LICENSE", "NOTICE", ".gitignore")] + [repo / n for n in ("README.md", ".dockerignore", ".env.example", "Dockerfile", "compose.yaml", "compose.run.yaml", "Makefile")]
excluded = {"node_modules", "dist", "test-results", "playwright-report", "__pycache__", ".git", ".cache", "tmp", "output"}

def eligible(p):
    parts = p.relative_to(top).parts
    if p.is_symlink() or any(part in excluded for part in parts):
        return False
    if p.name.startswith(".env") and p.name != ".env.example":
        return False
    if p.suffix in (".log", ".pyc", ".tsbuildinfo") or p.name in (".DS_Store", "api"):
        return False
    return True

files = [p for p in root_files if p.exists()]
for folder in (top / ".github", top / "LICENSES", repo / "src", repo / "data", repo / "scripts", repo / "docs", repo / "deploy"):
    files.extend(p for p in folder.rglob("*") if p.is_file() and eligible(p))
files = sorted(set(files))
if any(not eligible(p) for p in files):
    raise RuntimeError("Excluded path selected")
source_manifest = []
with zipfile.ZipFile(out / "LaporSiap-source-review.zip", "x", zipfile.ZIP_DEFLATED) as archive:
    for p in files:
        name = p.relative_to(top).as_posix()
        b = p.read_bytes()
        # Known credential fields in examples must be empty/placeholders. The
        # actual .env files are excluded structurally, including nested ones.
        if p.name == ".env.example" and b"APERTUS_API_KEY=" not in b:
            raise RuntimeError("Unexpected environment example")
        archive.writestr("LaporSiap/" + name, b)
        source_manifest.append({"path": name, "bytes": len(b), "sha256": hashlib.sha256(b).hexdigest()})
    archive.writestr("LaporSiap/REVIEW-ARCHIVE-NOTE.txt", "Local source snapshot for owner/reviewer inspection. Not published or submitted. Model weights and credentials are excluded. Historical model-output fixtures retain their upstream rights; NOTICE does not relicense them. Follow README and deploy/GUIDE.txt. Backend/frontend now use Clean Architecture and Redux. See docs/TECHNICAL.md and output/architecture/verification.json in the evidence archive. Local CI workflow is prepared and linted; GitHub execution has not occurred. Independent human testing, actual GPU inference and final event-guide verification remain outstanding.\n")

evidence = []
for folder in ("output/hardening", "output/evaluation", "output/architecture", "output/ci", "output/gpu-preparation"):
    evidence.extend(p for p in (repo / folder).rglob("*") if p.is_file() and not p.is_symlink() and p.name != ".DS_Store")
with zipfile.ZipFile(out / "LaporSiap-hardening-evidence.zip", "x", zipfile.ZIP_DEFLATED) as archive:
    for p in sorted(evidence):
        name = p.relative_to(repo).as_posix()
        if any(part.startswith(".env") for part in p.relative_to(repo).parts):
            raise RuntimeError("Private config selected")
        archive.write(p, name)
    archive.writestr("READ-FIRST.txt", "Fictional developer testing, not independent human sessions or a GPU deployment. Initial failures and subsequent checks are retained. browser-intermediate/context.json explains the backend-restart failure and mislabeled skipped videos; they are not successful deployment attempts. assistant-semantic-review.json records four issues in three accepted drafts. Initial two v4 requests timed out; two later recorded rechecks rejected changed quotes and an alternate response model. The initial recorder cache failure made zero provider calls. Later v5 source-bound EN/ID rechecks passed the contract, with one English timing omission recorded separately. Success video names mean contract acceptance, not semantic perfection. Later v6 rechecks preserved next week/pekan depan with strict generation-time validation. See output/evaluation/timing-v6-2026-10-11/verification.json and docs/evaluation/REVIEW-REPORT.txt for current AI results. GPU pins/plan/collector were prepared, but no GPU deployment or rental occurred; see output/gpu-preparation/2026-10-11/verification.json. output/architecture/verification.json describes the later Clean Architecture/Redux refactor; it made no real provider calls. Catalog/token cost is not an invoice. Third-party provider outputs/advisories retain upstream rights. No credentials, bearer tokens or weights are intentionally included. This is a local review archive; it has not been published/submitted.\n")

manifest = {"purpose": "local review, not published/submitted", "source_files": source_manifest, "archives": []}
for p in sorted(out.glob("*.zip")):
    with zipfile.ZipFile(p) as archive:
        if archive.testzip() is not None:
            raise RuntimeError("Archive CRC verification failed")
    manifest["archives"].append({"file": p.name, "bytes": p.stat().st_size, "sha256": hashlib.sha256(p.read_bytes()).hexdigest()})
(out / "manifest.json").write_text(json.dumps(manifest, indent=2) + "\n")

# Read configured secret values only for in-memory archive screening. Never log
# them or serialize their hashes; the actual private configuration is excluded.
secret_values = set()
secret_keys = {"APERTUS_API_KEY", "DEPLOY_APERTUS_API_KEY", "POSTGRES_PASSWORD", "DATABASE_URL", "HF_TOKEN"}
for config in (repo / ".env", top.parent / ".apertus-api.env", repo / "deploy/.env", repo / "deploy/.env.gpu"):
    if not config.exists():
        continue
    for line in config.read_text().splitlines():
        line = line.strip().removeprefix("export ")
        if "=" not in line or line.startswith("#"):
            continue
        key, value = line.split("=", 1)
        if key.strip() not in secret_keys:
            continue
        try:
            values = shlex.split(value, comments=True)
        except ValueError:
            raise RuntimeError("Cannot screen private configuration safely")
        if len(values) == 1 and len(values[0]) >= 8:
            secret_values.add(values[0].encode())
checked = 0
for archive_path in sorted(out.glob("*.zip")):
    with zipfile.ZipFile(archive_path) as archive:
        for name in archive.namelist():
            if any(part.startswith(".env") and part != ".env.example" for part in Path(name).parts):
                raise RuntimeError("Private environment path in archive")
            data = archive.read(name)
            if any(value in data for value in secret_values):
                raise RuntimeError("Configured secret found in archive; do not share")
            if re.search(rb'"token"\s*:\s*"[a-f0-9]{48}"', data):
                raise RuntimeError("Case bearer token found in archive; do not share")
            checked += 1
(out / "archive-verification.json").write_text(json.dumps({"checked_entries": checked, "crc_verified": True, "configured_secret_matches": 0, "case_bearer_token_matches": 0, "private_environment_paths": 0, "limitations": ["Screens known configured secret values and case-token JSON patterns only; not an independent security audit or image OCR."], "archives": manifest["archives"]}, indent=2)+"\n")

print(json.dumps({"output": str(out), "source_files": len(files), "archives": manifest["archives"]}, indent=2))
