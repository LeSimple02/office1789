#!/usr/bin/env python3
import subprocess
import sys
from pathlib import Path

VOLUMES = [
    "docker_docker_pgdata",
    "docker_docker_roundcube_db",
    "docker_docker_synapse_db",
    "docker_docker_synapse_data",
    "docker_docker_maildata",
    "docker_docker_mailstate",
    "docker_docker_maillogs",
]


def run(cmd, cwd=None):
    print("[CMD]", " ".join(cmd))
    return subprocess.run(cmd, cwd=cwd)


def main():
    repo = Path(__file__).resolve().parent.parent
    docker_dir = repo / "docker"

    print("[CLEAN] Stopping stack...")
    run(["docker", "compose", "down", "--remove-orphans"], cwd=docker_dir)

    print("[CLEAN] Removing named volumes (local only)...")
    run(["docker", "volume", "rm"] + VOLUMES)

    print("[NEXT] Redeploy local stack:")
    print("       python scripts/local_deploy.py")
    print("       cd docker && docker compose up -d")


if __name__ == "__main__":
    sys.exit(main())
