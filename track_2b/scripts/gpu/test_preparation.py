"""Offline regressions for hardware gating and safe launcher arguments."""
import unittest
from unittest.mock import patch
from types import SimpleNamespace
import json
from pathlib import Path
from hardware import requirements
from preflight import load_lock
from serve import command
from collect import inspect_container


class PreparationTests(unittest.TestCase):
    def setUp(self):
        self.lock = load_lock(Path(__file__).resolve().parents[2] / "deploy/gpu-lock.json")
        self.gpus = [{"memory_total_mib": 81000, "driver": "580.65.06"} for _ in range(4)]

    def test_adequate_linux_host(self):
        self.assertTrue(all(requirements(self.lock, self.gpus, 512, 300, "Linux").values()))

    def test_inadequate_hardware_and_driver_rejected(self):
        self.assertFalse(all(requirements(self.lock, [], None, 300, "Darwin").values()))
        self.assertFalse(all(requirements(self.lock, self.gpus[:2], 512, 300, "Linux").values()))
        old = [{"memory_total_mib": 81000, "driver": "570.0"} for _ in range(4)]
        self.assertFalse(requirements(self.lock, old, 512, 300, "Linux")["driver"])
        self.assertFalse(requirements(self.lock, self.gpus, 128, 300, "Linux")["host_ram"])
        self.assertFalse(requirements(self.lock, self.gpus, 512, 100, "Linux")["disk"])

    def test_launcher_pins_weights_and_keeps_credentials_outside_argv(self):
        cmd = command(self.lock)
        self.assertEqual(cmd[cmd.index("--revision") + 1], self.lock["model_revision"])
        self.assertEqual(cmd[cmd.index("--served-model-name") + 1], self.lock["served_model"])
        self.assertEqual(cmd[cmd.index("--host") + 1], "127.0.0.1")
        self.assertNotIn("--api-key", cmd)
        self.assertIn("--disable-log-requests", cmd)

    def test_evidence_excludes_container_credentials_and_raw_metadata(self):
        raw = [{"Config": {"Image": self.lock["runtime_image"], "Cmd": command(self.lock)[2:],
                           "Env": ["HF_TOKEN=PRIVATE_TEST_TOKEN"], "Labels": {"private": "PRIVATE_LABEL"}},
                "State": {"Running": True, "Error": "PRIVATE_ERROR"}}]
        responses = [SimpleNamespace(stdout=json.dumps(raw)), SimpleNamespace(stdout=json.dumps({"inference_health": True}))]
        with patch("collect.subprocess.run", side_effect=responses):
            result = inspect_container("fictional-container", self.lock)
        self.assertTrue(result["runtime_digest_matches"])
        self.assertTrue(result["revision_argument_matches"])
        self.assertNotIn("PRIVATE", json.dumps(result))


if __name__ == "__main__":
    unittest.main()
