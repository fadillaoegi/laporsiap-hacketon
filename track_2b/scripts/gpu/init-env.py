"""Create a separate private GPU deployment env without touching the local API."""
import argparse
import os
from pathlib import Path
import secrets
from preflight import load_lock

parser = argparse.ArgumentParser(description=__doc__)
parser.add_argument("--write", action="store_true", help="Create deploy/.env.gpu; refuse overwrite")
args = parser.parse_args()
root = Path(__file__).resolve().parents[2]
lock = load_lock(root / "deploy/gpu-lock.json")
if not args.write:
    print("Dry run: would create deploy/.env.gpu with pinned runtime/revision and fresh private keys. No file written.")
else:
    path = root / "deploy/.env.gpu"
    content = "\n".join([
        "DEPLOY_DB_PASSWORD=" + secrets.token_hex(24),
        "DEPLOY_APERTUS_API_KEY=" + secrets.token_hex(24),
        "APERTUS_RUNTIME_IMAGE=" + lock["runtime_image"],
        "APERTUS_WEIGHTS_REVISION=" + lock["model_revision"],
        "APERTUS_TENSOR_PARALLEL=4", "GPU_APP_ORIGIN=http://127.0.0.1:8878", "HF_TOKEN=",
    ]) + "\n"
    fd = os.open(path, os.O_WRONLY | os.O_CREAT | os.O_EXCL, 0o600)
    with os.fdopen(fd, "w") as f:
        f.write(content)
    print("Created private deploy/.env.gpu. Existing local/provider env unchanged. Owner must add authorized HF_TOKEN.")
