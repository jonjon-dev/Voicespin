import argparse
import json
from collections import Counter
from pathlib import Path
import sys


def summarize_events(file_path: str) -> dict:
    path = Path(file_path)

    valid_events = 0
    invalid_lines = 0
    event_counts: Counter[str] = Counter()

    if not path.exists() or not path.is_file():
        return {
            "valid_events": valid_events,
            "invalid_lines": invalid_lines,
            "event_counts": dict(event_counts),
            "error": f"File not found: {file_path}",
        }

    with path.open("r", encoding="utf-8") as f:
        for line in f:
            stripped = line.strip()
            if not stripped:
                invalid_lines += 1
                continue

            try:
                data = json.loads(stripped)
            except json.JSONDecodeError:
                invalid_lines += 1
                continue

            if not isinstance(data, dict):
                invalid_lines += 1
                continue

            event_type = data.get("event_type")
            if not isinstance(event_type, str) or not event_type.strip():
                invalid_lines += 1
                continue

            valid_events += 1
            event_counts[event_type] += 1

    return {
        "valid_events": valid_events,
        "invalid_lines": invalid_lines,
        "event_counts": dict(event_counts),
        "error": None,
    }


def main() -> None:
    parser = argparse.ArgumentParser(description="Summarize a JSONL events file.")
    parser.add_argument("file_path", help="Path to the JSONL file")
    args = parser.parse_args()

    summary = summarize_events(args.file_path)

    if summary["error"]:
        print(f"Warning: {summary['error']}", file=sys.stderr)

    print(f"Valid events: {summary['valid_events']}")
    print(f"Invalid lines: {summary['invalid_lines']}")
    print("Events by type:")
    for event_type, count in sorted(summary["event_counts"].items()):
        print(f"  - {event_type}: {count}")


if __name__ == "__main__":
    main()