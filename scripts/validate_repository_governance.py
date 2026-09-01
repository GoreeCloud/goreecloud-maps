#!/usr/bin/env python3
"""Fail-closed validation for mandatory GoreeCloud Maps root governance records."""

from __future__ import annotations

from pathlib import Path
import sys

ROOT = Path(__file__).resolve().parents[1]

REQUIRED_RECORDS = {
    "README.md": "# GoreeCloud Maps",
    "SPECIFICATIONS.md": "# GoreeCloud Maps Specifications",
    "FEATURES.md": "# GoreeCloud Maps Features",
    "BENEFITS.md": "# GoreeCloud Maps Benefits",
    "COMPETITIVE-OBJECTIVES.md": "# GoreeCloud Maps Competitive Objectives",
    "BRANDING.md": "# GoreeCloud Maps Branding",
}

LICENSE_MARKERS = (
    "GNU AFFERO GENERAL PUBLIC LICENSE",
    "Version 3, 19 November 2007",
)


def read_required_file(relative: str, errors: list[str]) -> str | None:
    path = ROOT / relative
    if not path.is_file() or path.is_symlink():
        errors.append(f"required root file is missing, not a regular file, or is a symlink: {relative}")
        return None
    try:
        return path.read_text(encoding="utf-8")
    except (OSError, UnicodeError) as exc:
        errors.append(f"required root file is not readable UTF-8: {relative}: {exc.__class__.__name__}")
        return None


def main() -> int:
    errors: list[str] = []

    for relative, expected_heading in REQUIRED_RECORDS.items():
        text = read_required_file(relative, errors)
        if text is None:
            continue
        lines = text.splitlines()
        first_line = lines[0] if lines else ""
        if first_line != expected_heading:
            errors.append(
                f"unexpected governance identity heading in {relative}: expected {expected_heading!r}"
            )
        if len(text.strip()) < len(expected_heading) + 80:
            errors.append(f"governance record is unexpectedly skeletal: {relative}")

    license_text = read_required_file("LICENSE", errors)
    if license_text is not None:
        for marker in LICENSE_MARKERS:
            if marker not in license_text:
                errors.append(f"LICENSE does not contain expected AGPL-3.0 marker: {marker!r}")

    if errors:
        print("GoreeCloud Maps repository governance validation failed:")
        for error in errors:
            print(f"  - {error}")
        return 1

    print(
        "GoreeCloud Maps repository governance validation passed: six mandatory root records and "
        "the explicit AGPL-3.0 license material are present and structurally valid."
    )
    return 0


if __name__ == "__main__":
    sys.exit(main())
