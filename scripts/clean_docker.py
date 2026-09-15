import subprocess
import sys
from pathlib import Path


def run(cmd, cwd: Path) -> int:
    print(f"[CMD] {' '.join(cmd)}")
    return subprocess.run(cmd, cwd=str(cwd)).returncode


def main() -> None:
    script_dir = Path(__file__).resolve().parent
    project_root = script_dir.parent
    docker_dir = project_root / "docker"

    if not docker_dir.is_dir():
        print(f"[ERR] Docker directory not found at {docker_dir}")
        sys.exit(1)

    print("[INFO] Cleaning docker-compose stack: containers, volumes and local build images...")
    # --rmi local supprime les images construites par ce compose (les builds)
    code = run(["docker", "compose", "down", "-v", "--rmi", "local"], cwd=docker_dir)
    if code != 0:
        sys.exit(code)

    print("[DONE] All containers, volumes and local build images have been removed.")


if __name__ == "__main__":
    main()

