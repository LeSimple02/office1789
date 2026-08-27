#!/usr/bin/env python3
import os
from pathlib import Path
import subprocess

# Paths
repo_dir = Path(__file__).parent.parent.resolve()
docker_dir = repo_dir / "docker"
env_file = docker_dir / ".env"

print("==========================================")
print("🧹 CLEANUP DOCKER ENV")
print("==========================================")

# Remove docker/.env if it exists
if env_file.exists():
    env_file.unlink()
    print(f"Removed {env_file}")
else:
    print(f"No {env_file} to remove.")

# Bring down docker compose (with -v to remove volumes)
os.chdir(docker_dir)
print("Bringing down docker compose (with volumes)...")
subprocess.run(["docker", "compose", "down", "-v"])
print("Done.")
