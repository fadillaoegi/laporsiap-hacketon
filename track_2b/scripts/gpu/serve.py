"""Dry-run by default; --serve explicitly starts the official pinned GPU model."""
import argparse
import json
import os
from pathlib import Path
import shutil
import subprocess
from preflight import load_lock
from hardware import read_gpus, host_memory_gib, requirements


def command(lock):
    return ["vllm", "serve", lock["model_repository"], "--revision", lock["model_revision"],
            "--served-model-name", lock["served_model"], "--host", "127.0.0.1", "--port", "9000",
            "--chat-template-content-format", "string", "--dtype", lock["dtype"],
            "--tensor-parallel-size", str(lock["tensor_parallel_size"]),
            "--gpu-memory-utilization", str(lock["gpu_memory_utilization"]),
            "--max-model-len", str(lock["max_model_len"]), "--max-num-seqs", str(lock["max_num_seqs"]),
            "--disable-log-requests", "--disable-uvicorn-access-log"]


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--lock", default=str(Path(__file__).resolve().parents[2] / "deploy/gpu-lock.json"))
    parser.add_argument("--serve", action="store_true")
    args = parser.parse_args()
    lock = load_lock(args.lock)
    cmd = command(lock)
    if not args.serve:
        print(json.dumps({"mode": "dry_run", "command": cmd, "inference_started": False,
                          "credentials_required": ["HF_TOKEN", "VLLM_API_KEY"]}, indent=2))
        return 0
    if not os.environ.get("HF_TOKEN") or len(os.environ.get("VLLM_API_KEY", "")) < 24:
        raise SystemExit("Set an authorized HF_TOKEN and a private VLLM_API_KEY; values are never printed.")
    cache = Path(os.environ.get("HF_HOME", str(Path.home() / ".cache/huggingface")))
    disk = cache if cache.exists() else cache.parent
    if not disk.exists():
        disk = Path(".")
    gpus, error = read_gpus()
    if error or not all(requirements(lock, gpus, host_memory_gib(), shutil.disk_usage(disk).free / (1024 ** 3)).values()):
        raise SystemExit("GPU/RAM/disk/driver preflight failed. Run preflight.py --probe before starting.")
    try:
        help_result = subprocess.run(["vllm", "serve", "--help"], capture_output=True, text=True, timeout=30, check=True)
        for flag in [value for value in cmd if value.startswith("--")]:
            if flag not in help_result.stdout:
                raise SystemExit("Pinned runtime does not advertise a required serve option; do not start this profile.")
    except (OSError, subprocess.SubprocessError):
        raise SystemExit("Cannot verify installed vLLM options; no model started.")
    env = dict(os.environ, VLLM_NO_USAGE_STATS="1", DO_NOT_TRACK="1", HF_HUB_DISABLE_TELEMETRY="1")
    os.execvpe("vllm", cmd, env)


if __name__ == "__main__":
    raise SystemExit(main())
