#!/usr/bin/env python3
"""Validate skill registry packages.

Two validation layers:
  Layer 1 — JSON Schema validation of metadata.json
  Layer 2 — Semantic rules (name matches dir, artifact exists, etc.)

Usage:
  python ci/validate.py packages/<name>    # validate one package
  python ci/validate.py                    # validate all packages
"""

import json
import os
import sys
from pathlib import Path

try:
    import jsonschema
except ImportError:
    print("ERROR: jsonschema not installed. Run: pip install jsonschema", file=sys.stderr)
    sys.exit(1)

REPO_ROOT = Path(__file__).resolve().parent.parent
SCHEMA_PATH = REPO_ROOT / "ci" / "schema" / "metadata.schema.json"

ARTIFACT_FILENAMES = {
    "skill": "SKILL.md",
    "rule": "RULE.md",
    "knowledge": "KNOWLEDGE.md",
}


def load_schema():
    with open(SCHEMA_PATH) as f:
        return json.load(f)


def validate_schema(metadata, schema):
    """Layer 1: JSON Schema validation. Returns list of error strings."""
    errors = []
    validator = jsonschema.Draft7Validator(schema)
    for error in sorted(validator.iter_errors(metadata), key=lambda e: list(e.path)):
        path = ".".join(str(p) for p in error.absolute_path) or "(root)"
        errors.append(f"  schema: {path}: {error.message}")
    return errors


def validate_semantic(pkg_dir, metadata):
    """Layer 2: Semantic validation. Returns list of error strings."""
    errors = []
    pkg_name = pkg_dir.name

    # name must match directory
    if metadata.get("name") != pkg_name:
        errors.append(
            f"  semantic: 'name' is '{metadata.get('name')}' but directory is '{pkg_name}'"
        )

    # artifact file must exist and match type
    pkg_type = metadata.get("type", "")
    expected_artifact = ARTIFACT_FILENAMES.get(pkg_type)
    if expected_artifact is None:
        errors.append(f"  semantic: unknown type '{pkg_type}'")
    else:
        artifact_path = pkg_dir / expected_artifact
        if not artifact_path.exists():
            errors.append(
                f"  semantic: expected artifact '{expected_artifact}' not found"
            )

    # tags must contain no empty strings (schema handles minItems)
    for i, tag in enumerate(metadata.get("tags", [])):
        if not tag.strip():
            errors.append(f"  semantic: tags[{i}] is empty or whitespace")

    # upstream validation
    upstream = metadata.get("upstream")
    if upstream:
        url = upstream.get("url", "")
        if not url.startswith(("http://", "https://")):
            errors.append(
                f"  semantic: upstream.url '{url}' is not a valid HTTP(S) URL"
            )

    # README.md must exist (optional for upstream-tracked packages)
    if not upstream and not (pkg_dir / "README.md").exists():
        errors.append("  semantic: README.md not found")

    return errors


def validate_package(pkg_dir, schema):
    """Validate a single package. Returns list of error strings."""
    errors = []
    metadata_path = pkg_dir / "metadata.json"

    if not metadata_path.exists():
        return [f"  metadata.json not found in {pkg_dir}"]

    try:
        with open(metadata_path) as f:
            metadata = json.load(f)
    except json.JSONDecodeError as e:
        return [f"  metadata.json is invalid JSON: {e}"]

    errors.extend(validate_schema(metadata, schema))
    errors.extend(validate_semantic(pkg_dir, metadata))

    return errors


def find_packages():
    """Find all package directories under packages/."""
    packages_dir = REPO_ROOT / "packages"
    if not packages_dir.exists():
        return []
    return sorted(
        d for d in packages_dir.iterdir()
        if d.is_dir() and not d.name.startswith(".")
    )


def main():
    schema = load_schema()
    failed = False

    if len(sys.argv) > 1:
        # Validate specific package(s)
        targets = [Path(arg).resolve() for arg in sys.argv[1:]]
    else:
        # Validate all packages
        targets = find_packages()

    if not targets:
        print("No packages to validate.")
        return

    for pkg_dir in targets:
        if not pkg_dir.is_dir():
            print(f"FAIL {pkg_dir}: not a directory")
            failed = True
            continue

        errors = validate_package(pkg_dir, schema)
        if errors:
            print(f"FAIL {pkg_dir.name}")
            for e in errors:
                print(e)
            failed = True
        else:
            print(f"OK   {pkg_dir.name}")

    if failed:
        sys.exit(1)


if __name__ == "__main__":
    main()
