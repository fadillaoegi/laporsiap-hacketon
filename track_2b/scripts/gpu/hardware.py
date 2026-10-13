"""Read-only host measurements and conservative checks for the pinned GPU plan."""
import csv
import io
import platform
import re
import shutil
import subprocess
from pathlib import Path


def read_gpus():
    if not shutil.which("nvidia-smi"):
        return [], "NVIDIA GPU tooling is unavailable on this host."
    try:
        result = subprocess.run([
            "nvidia-smi", "--query-gpu=index,name,memory.total,memory.used,utilization.gpu,driver_version",
            "--format=csv,noheader,nounits",
        ], capture_output=True, text=True, timeout=10, check=True)
        gpus = []
        for row in csv.reader(io.StringIO(result.stdout)):
            index, name, total, used, utilization, driver = [v.strip() for v in row]
            if not re.fullmatch(r"[A-Za-z0-9 .()_-]{1,120}", name) or not re.fullmatch(r"\d+(?:\.\d+)+", driver):
                raise ValueError("Invalid GPU metadata")
            gpus.append({"index": int(index), "name": name, "memory_total_mib": int(total),
                         "memory_used_mib": int(used), "utilization_percent": int(utilization), "driver": driver})
        return gpus, None
    except (OSError, subprocess.SubprocessError, ValueError):
        return [], "NVIDIA measurements failed; no raw command output retained."


def host_memory_gib():
    try:
        data = Path("/proc/meminfo").read_text()
        return int(re.search(r"^MemTotal:\s+(\d+) kB", data, re.M).group(1)) / (1024 ** 2)
    except (OSError, AttributeError, ValueError):
        return None


def requirements(lock, gpus, ram_gib, free_disk_gib, system=None):
    system = system or platform.system()
    checks = {
        "linux": system == "Linux",
        "gpu_count": len(gpus) >= lock["tensor_parallel_size"],
        "gpu_memory": len(gpus) >= lock["tensor_parallel_size"] and all(
            g["memory_total_mib"] >= lock["minimum_gpu_memory_mib_each"]
            for g in gpus[:lock["tensor_parallel_size"]]),
        "driver": bool(gpus) and all(int(g["driver"].split(".")[0]) >= lock["driver_minimum_major"] for g in gpus),
        "host_ram": ram_gib is not None and ram_gib >= lock["minimum_host_ram_gib"],
        "disk": free_disk_gib >= lock["minimum_free_disk_gib"],
    }
    return checks
