#!/usr/bin/env python3
"""Fail-closed validation for mandatory GoreeCloud Maps repository governance."""

from __future__ import annotations

from pathlib import Path
import sys

ROOT = Path(__file__).resolve().parents[1]

REQUIRED_RECORDS = {
    "README.md": "# GoreeCloud Maps",
    "PROJECT-SPECIFICATIONS.md": "# GoreeCloud Maps — Project Specifications",
    "PROJECT-RECORD.md": "# GoreeCloud Maps — Project Record",
    "SPECIFICATIONS.md": "# GoreeCloud Maps Specifications",
    "FEATURES.md": "# GoreeCloud Maps Features",
    "IMPLEMENTED-FEATURES.md": "# GoreeCloud Maps — Implemented Features",
    "PLANNED-FEATURES.md": "# GoreeCloud Maps — Planned Features",
    "CHANGELOGS.md": "# GoreeCloud Maps — Changelogs",
    "BENEFITS.md": "# GoreeCloud Maps Benefits",
    "COMPETITIVE-OBJECTIVES.md": "# GoreeCloud Maps Competitive Objectives",
    "BRANDING.md": "# GoreeCloud Maps Branding",
    "USER-MANUAL.md": "# GoreeCloud Maps — User Manual",
    "PRIVACY POLICY.md": "# GoreeCloud Maps — Privacy Policy",
    "NOTES.md": "# GoreeCloud Maps — Notes",
    "SECURITY.md": "# GoreeCloud Maps — Security",
}

REQUIRED_REGULAR_FILES = (
    ".gitignore",
    ".editorconfig",
    "goreecloud.platform.yaml",
)

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

    for relative in REQUIRED_REGULAR_FILES:
        read_required_file(relative, errors)

    license_text = read_required_file("LICENSE", errors)
    if license_text is not None:
        for marker in LICENSE_MARKERS:
            if marker not in license_text:
                errors.append(f"LICENSE does not contain expected AGPL-3.0 marker: {marker!r}")

    if (ROOT / "FEATURE-ROADMAP.md").exists():
        errors.append(
            "retired FEATURE-ROADMAP.md exists; feature authority must use "
            "IMPLEMENTED-FEATURES.md and PLANNED-FEATURES.md"
        )

    for relative in ("README.md", "SPECIFICATIONS.md"):
        text = (ROOT / relative).read_text(encoding="utf-8") if (ROOT / relative).is_file() else ""
        if "GoreeCloud/Projects/Project Specification — Maps" in text:
            errors.append(
                f"{relative} still declares the retired Drive project specification as current authority"
            )

    if errors:
        print("GoreeCloud Maps repository governance validation failed:")
        for error in errors:
            print(f"  - {error}")
        return 1

    print(
        "GoreeCloud Maps repository governance validation passed: canonical project, feature, "
        "change-history, user/privacy/security, platform, editor, branding, and license records "
        "are present; the retired feature-roadmap file and stale Drive authority are absent."
    )
    return 0


if __name__ == "__main__":
    sys.exit(main())
