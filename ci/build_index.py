#!/usr/bin/env python3
"""Build marketplace.json from packages/.

Reads all packages/<name>/metadata.json files, looks up the latest git tag
for each package (<name>@<semver>), computes artifact hashes, and writes
the index to marketplace.json at the repo root.

Usage:
  python ci/build_index.py
"""

import hashlib
import json
import os
import subprocess
import sys
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent
PACKAGES_DIR = REPO_ROOT / "packages"
INDEX_PATH = REPO_ROOT / "marketplace.json"
DEFAULT_VERSION = "0.0.1"

ARTIFACT_FILENAMES = {
    "skill": "SKILL.md",
    "rule": "RULE.md",
    "knowledge": "KNOWLEDGE.md",
}


def get_latest_tag(name):
    """Get the latest version tag for a package (<name>@<semver>)."""
    try:
        result = subprocess.run(
            ["git", "tag", "--list", f"{name}@*", "--sort=-v:refname"],
            capture_output=True,
            text=True,
            cwd=REPO_ROOT,
        )
        tags = result.stdout.strip().splitlines()
        if tags:
            # Extract version from tag like "my-package@1.2.3"
            return tags[0].split("@", 1)[1]
    except (subprocess.CalledProcessError, FileNotFoundError, IndexError):
        pass
    return DEFAULT_VERSION


def hash_file(path):
    """Compute SHA256 hash of a file."""
    sha = hashlib.sha256()
    with open(path, "rb") as f:
        for chunk in iter(lambda: f.read(8192), b""):
            sha.update(chunk)
    return f"sha256:{sha.hexdigest()}"


def list_package_files(pkg_dir):
    """List all files in a package directory, relative to the package root."""
    files = []
    for root, dirs, filenames in os.walk(pkg_dir):
        # Skip hidden directories
        dirs[:] = [d for d in dirs if not d.startswith(".")]
        for fname in filenames:
            if fname.startswith("."):
                continue
            rel = os.path.relpath(os.path.join(root, fname), pkg_dir)
            files.append(rel)
    return sorted(files)


def build_entry(pkg_dir):
    """Build an index entry for a single package."""
    metadata_path = pkg_dir / "metadata.json"
    if not metadata_path.exists():
        return None

    with open(metadata_path) as f:
        metadata = json.load(f)

    name = metadata["name"]
    pkg_type = metadata["type"]

    # Find and hash the artifact
    artifact_name = ARTIFACT_FILENAMES.get(pkg_type)
    artifact_path = pkg_dir / artifact_name if artifact_name else None
    artifact_hash = hash_file(artifact_path) if artifact_path and artifact_path.exists() else None

    version = get_latest_tag(name)

    entry = {
        "name": name,
        "type": pkg_type,
        "version": version,
        "description": metadata["description"],
        "tags": metadata["tags"],
        "author": metadata["author"],
        "license": metadata["license"],
        "path": f"packages/{name}",
        "files": list_package_files(pkg_dir),
    }

    if artifact_hash:
        entry["artifact_hash"] = artifact_hash

    if metadata.get("notes"):
        entry["notes"] = metadata["notes"]

    if metadata.get("install_target"):
        entry["install_target"] = metadata["install_target"]

    if metadata.get("upstream"):
        entry["upstream"] = metadata["upstream"]

    return entry


def main():
    if not PACKAGES_DIR.exists():
        print("No packages/ directory found.")
        sys.exit(1)

    packages = sorted(
        d for d in PACKAGES_DIR.iterdir()
        if d.is_dir() and not d.name.startswith(".")
    )

    entries = []
    for pkg_dir in packages:
        entry = build_entry(pkg_dir)
        if entry:
            entries.append(entry)
            print(f"  indexed: {entry['name']}@{entry['version']}")

    index = {"packages": entries}

    with open(INDEX_PATH, "w") as f:
        json.dump(index, f, indent=2)
        f.write("\n")

    print(f"\nWrote {len(entries)} package(s) to {INDEX_PATH}")


if __name__ == "__main__":
    main()
