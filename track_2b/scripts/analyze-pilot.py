"""Summarize observed human tasks; templates do not become test results."""
import argparse
import json
from datetime import datetime
from pathlib import Path
from statistics import median

parser = argparse.ArgumentParser()
parser.add_argument("results", type=Path)
args = parser.parse_args()
data = json.loads(args.results.read_text())
tasks = data["tasks"]
observed = [t for t in tasks if t["status"] != "not_started"]
issues = []
complete = []
for t in observed:
    ident = t["task_id"]
    if not t.get("participant_role_actual"):
        issues.append(f"{ident}: participant role unverified")
    if t["status"] not in ("completed_unassisted", "completed_assisted"):
        continue
    try:
        start, end = [datetime.fromisoformat(t[k]) for k in ("start_at", "finish_at")]
        if start.tzinfo is None or end.tzinfo is None:
            raise ValueError("timezone required")
        elapsed = (end - start).total_seconds()
        pause = t["off_task_pause_seconds"]
        total = t["total_task_seconds"]
        if not isinstance(pause, (int, float)) or not isinstance(total, (int, float)) or not 0 <= pause <= elapsed or abs(elapsed - pause - total) > 1:
            raise ValueError("timer fields inconsistent")
        if not t.get("recording_path") or not t.get("quality_review_notes") or t.get("errors_remaining") is None:
            raise ValueError("recording and content review required")
        comp = t.get("comprehension") or {}
        if any(not comp.get(k) for k in ("total_explanation", "unknown_fields_explanation", "source_trace_explanation", "official_submission_explanation", "facilitator_assessment")):
            raise ValueError("comprehension assessment missing")
        complete.append(t)
    except (ValueError, TypeError, KeyError) as e:
        issues.append(f"{ident}: {e}")
independent = [t for t in complete if t["participant_role_actual"] == "independent"]
summary = {
    "observed_tasks": len(observed),
    "completed_with_review": len(complete),
    "independent_participants_with_review": len({t["participant_id"] for t in independent}),
    "role_rule": "Only the explicitly verified role independent counts. Planned participant IDs do not establish independence.",
    "issues": issues,
    "modes": {m: {"reviewed_tasks": len(ts), "median_task_seconds": median([t["total_task_seconds"] for t in ts]) if ts else None} for m in ("manual", "apertus") for ts in [[t for t in independent if t["mode"] == m]]},
    "limitations": ["Small exploratory pilot, no population accuracy or time-saving claim.", "Null means unobserved, not zero.", "Role verification remains the facilitator's responsibility."]
}
print(json.dumps(summary, indent=2))
