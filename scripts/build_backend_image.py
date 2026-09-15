import argparse
import shutil
import subprocess
import sys
from datetime import datetime
from pathlib import Path


def run(cmd, cwd=None):
    print(f"[CMD] {' '.join(cmd)}")
    result = subprocess.run(cmd, cwd=cwd)
    if result.returncode != 0:
        sys.exit(result.returncode)


def main():
    parser = argparse.ArgumentParser(
        description="Build the backend Docker image locally and optionally save it to a tarball for transfer."
    )
    parser.add_argument(
        "--tag",
        default=f"office1789/backend:{datetime.now().strftime('%Y%m%d')}",
        help="Docker image tag to use (default: office1789/backend:YYYYMMDD)",
    )
    parser.add_argument(
        "--save",
        type=Path,
        help="Optional path to save the image tarball (e.g., ./backend-image.tar)",
    )
    parser.add_argument(
        "--push",
        action="store_true",
        help="Also push the image to the configured registry (requires docker login)",
    )
    args = parser.parse_args()

    repo_root = Path(__file__).resolve().parent.parent
    backend_dir = repo_root / "backend"
    dockerfile = backend_dir / "Dockerfile"

    if not shutil.which("docker"):
        print("[ERR] docker CLI not found in PATH.")
        sys.exit(1)

    if not dockerfile.exists():
        print(f"[ERR] Dockerfile not found at {dockerfile}")
        sys.exit(1)

    print(f"[INFO] Building backend image with tag: {args.tag}")
    run(["docker", "build", "-t", args.tag, "-f", str(dockerfile), str(backend_dir)])

    if args.save:
        save_path = args.save.resolve()
        save_path.parent.mkdir(parents=True, exist_ok=True)
        print(f"[INFO] Saving image to: {save_path}")
        run(["docker", "save", args.tag, "-o", str(save_path)])
        print(f"[NEXT] Transfer {save_path} to the server, then run:")
        print(f"       docker load < {save_path.name}")
        print(f"       # in docker/.env on the server set BACKEND_IMAGE={args.tag}")
        print("       docker compose -f docker/docker-compose.yml up -d backend")

    if args.push:
        print(f"[INFO] Pushing image: {args.tag}")
        run(["docker", "push", args.tag])

    print("[DONE] Backend image ready.")


if __name__ == "__main__":
    main()
