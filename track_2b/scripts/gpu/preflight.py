"""Validate preparation, or inspect real GPU hardware without launching inference."""
import argparse
from datetime import datetime, timezone
import json
from pathlib import Path
import re
import shutil
from hardware import read_gpus, host_memory_gib, requirements


def load_lock(path):
    lock = json.loads(Path(path).read_text())
    if lock["model_repository"] != "swiss-ai/Apertus-v1.5-70B" or lock["served_model"] != "swiss-ai/apertus-v1.5-70b":
        raise ValueError("Use the official 70B model and the application alias.")
    if not re.fullmatch(r"[a-f0-9]{40}", lock["model_revision"]):
        raise ValueError("Model revision must be a pinned commit.")
    if not re.fullmatch(r"ghcr.io/swiss-ai/vllm_apertus_1\.5_release@sha256:[a-f0-9]{64}", lock["runtime_image"]):
        raise ValueError("Runtime image must use an official immutable digest.")
    if lock["tensor_parallel_size"] != 4 or lock["dtype"] != "bfloat16":
        raise ValueError("This prepared profile uses four GPUs and BF16.")
    return lock


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--lock", default=str(Path(__file__).resolve().parents[2] / "deploy/gpu-lock.json"))
    parser.add_argument("--probe", action="store_true", help="Read actual host GPU/RAM/disk; no inference or purchase")
    parser.add_argument("--disk", default=".", help="Existing intended model-cache filesystem")
    parser.add_argument("--output", help="New JSON report file; never overwritten")
    args = parser.parse_args()
    lock = load_lock(args.lock)
    report = {"checked_at": datetime.now(timezone.utc).isoformat(), "mode": "hardware_probe" if args.probe else "preparation_only",
              "model_repository": lock["model_repository"], "model_revision": lock["model_revision"],
              "runtime_image": lock["runtime_image"], "configuration_valid": True,
              "gpu_deployment_proven": False, "inference_requests": 0, "resources_rented": 0,
              "limitations": ["Pinned metadata and preparation do not prove model execution.",
                              "Hardware capacity does not prove loaded weights or runtime compatibility."]}
    ok = True
    if args.probe:
        gpus, error = read_gpus()
        ram = host_memory_gib()
        disk = shutil.disk_usage(args.disk).free / (1024 ** 3)
        report.update(gpus=gpus, host_ram_gib=ram, free_disk_gib=disk)
        if error:
            report["gpu_error"] = error
        report["requirements"] = requirements(lock, gpus, ram, disk)
        ok = all(report["requirements"].values())
        report["hardware_preflight_passed"] = ok
    raw = json.dumps(report, indent=2) + "\n"
    if args.output:
        with Path(args.output).open("x") as f:
            f.write(raw)
    print(raw, end="")
    return 0 if ok else 1


if __name__ == "__main__":
    raise SystemExit(main())
