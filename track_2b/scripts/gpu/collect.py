"""Record GPU/runtime/cache/health observations; never send a completion."""
import argparse
from datetime import datetime, timezone
import importlib.metadata
import json
import os
from pathlib import Path
import re
import subprocess
import time
import urllib.request
from hardware import read_gpus
from preflight import load_lock


def utc():
    return datetime.now(timezone.utc).isoformat()


def runtime(lock):
    info = {"python_packages": {}, "checkpoint_files_present": False,
            "inference_health": False, "served_model_listed": False}
    for package in ("torch", "vllm", "transformers"):
        try:
            version = importlib.metadata.version(package)
            if re.fullmatch(r"[A-Za-z0-9.+_-]{1,100}", version):
                info["python_packages"][package] = version
        except importlib.metadata.PackageNotFoundError:
            pass
    cache = Path(os.environ.get("HF_HOME", str(Path.home() / ".cache/huggingface"))) / "hub"
    snapshot = cache / "models--swiss-ai--Apertus-v1.5-70B" / "snapshots" / lock["model_revision"]
    try:
        index = json.loads((snapshot / "model.safetensors.index.json").read_text())
        names = sorted(set(index["weight_map"].values()))
        if not names or any(Path(n).name != n or not n.endswith(".safetensors") for n in names):
            raise ValueError("Invalid checkpoint file list")
        info["checkpoint_files"] = [{"filename": n, "bytes": (snapshot / n).stat().st_size} for n in names]
        info["checkpoint_files_present"] = all(item["bytes"] > 0 for item in info["checkpoint_files"])
    except (OSError, KeyError, ValueError):
        info["checkpoint_error"] = "Expected pinned checkpoint files not verified in this cache."
    key = os.environ.get("VLLM_API_KEY")
    if key:
        # No proxy or redirects: keep the inference credential on loopback.
        class NoRedirect(urllib.request.HTTPRedirectHandler):
            def redirect_request(self, *args, **kwargs):
                return None
        opener = urllib.request.build_opener(urllib.request.ProxyHandler({}), NoRedirect())
        try:
            req = urllib.request.Request("http://127.0.0.1:9000/health", headers={"Authorization": "Bearer " + key})
            with opener.open(req, timeout=5) as response:
                info["inference_health"] = response.status == 200
            req = urllib.request.Request("http://127.0.0.1:9000/v1/models", headers={"Authorization": "Bearer " + key})
            with opener.open(req, timeout=5) as response:
                body = response.read(256001)
                if len(body) > 256000:
                    raise ValueError("Oversized catalog")
                catalog = json.loads(body)
                info["served_model_listed"] = any(m.get("id") == lock["served_model"] for m in catalog.get("data", []))
        except Exception:
            info["inference_error"] = "Loopback health/catalog not verified; raw responses are not retained."
    else:
        info["inference_error"] = "Private inference access key not available in this process."
    return info


def inspect_container(name, lock):
    # Decode locally but export only whitelisted facts, never Env/labels/logs.
    result = subprocess.run(["docker", "inspect", name], capture_output=True, text=True, timeout=10, check=True)
    item = json.loads(result.stdout)[0]
    argv = item.get("Config", {}).get("Cmd", [])
    def option(name):
        try:
            return argv[argv.index(name) + 1]
        except (ValueError, IndexError):
            return None
    observed = {
        "running": item.get("State", {}).get("Running") is True,
        "runtime_digest_matches": item.get("Config", {}).get("Image") == lock["runtime_image"],
        "model_argument_matches": lock["model_repository"] in argv,
        "revision_argument_matches": option("--revision") == lock["model_revision"],
        "alias_argument_matches": option("--served-model-name") == lock["served_model"],
        "dtype_argument_matches": option("--dtype") == lock["dtype"],
    }
    result = subprocess.run(["docker", "exec", name, "python3", "/opt/laporsiap-gpu/collect.py",
                             "--inside", "--lock", "/opt/gpu-lock.json"],
                            capture_output=True, text=True, timeout=20, check=True)
    observed["runtime"] = json.loads(result.stdout)
    return observed


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--lock", default=str(Path(__file__).resolve().parents[2] / "deploy/gpu-lock.json"))
    parser.add_argument("--container", help="Running Compose inference container on the current host")
    parser.add_argument("--inside", action="store_true", help=argparse.SUPPRESS)
    parser.add_argument("--seconds", type=int, default=30, help="Telemetry window (1–300 seconds)")
    parser.add_argument("--output", help="New report file; never overwritten")
    args = parser.parse_args()
    lock = load_lock(args.lock)
    if args.inside:
        print(json.dumps(runtime(lock)))
        return 0
    if not 1 <= args.seconds <= 300:
        parser.error("Use 1–300 seconds")
    report = {"started_at": utc(), "model_revision": lock["model_revision"], "runtime_image": lock["runtime_image"],
              "completion_requests": 0, "gpu_deployment_proven": False, "samples": [],
              "limitations": ["Telemetry and health must be paired with actual LaporSiap AI/review/export evidence.",
                              "Cache file presence and launch arguments are provenance evidence, not an independent weights audit."]}
    first, error = read_gpus()
    ok = bool(first) and error is None
    if error:
        report["gpu_error"] = error
    if ok:
        end = time.monotonic() + args.seconds
        while True:
            gpus, error = read_gpus()
            report["samples"].append({"at": utc(), "gpus": gpus})
            if error:
                ok = False
                report["gpu_error"] = error
                break
            remaining = end - time.monotonic()
            if remaining <= 0:
                break
            time.sleep(min(2, remaining))
        try:
            if args.container:
                report["container"] = inspect_container(args.container, lock)
                info = report["container"]["runtime"]
                ok = ok and all(value is True for key, value in report["container"].items() if key != "runtime")
            else:
                report["runtime"] = runtime(lock)
                info = report["runtime"]
                report["runtime_digest_verified"] = False
                report["limitations"].append("Pod mode needs provider-console image digest and launch configuration evidence separately.")
            ok = ok and info["inference_health"] and info["served_model_listed"] and info["checkpoint_files_present"]
        except (OSError, subprocess.SubprocessError, KeyError, ValueError):
            ok = False
            report["runtime_error"] = "Runtime observation failed; no raw inspect, stderr or private configuration retained."
    report.update(finished_at=utc(), observation_passed=ok)
    raw = json.dumps(report, indent=2) + "\n"
    if args.output:
        with Path(args.output).open("x") as f:
            f.write(raw)
    print(raw, end="")
    return 0 if ok else 1


if __name__ == "__main__":
    raise SystemExit(main())
