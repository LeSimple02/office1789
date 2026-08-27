#!/usr/bin/env python3
import crypt
import sys
from pathlib import Path
from typing import Optional


def load_env_value(env_path: Path, key: str) -> Optional[str]:
    if not env_path.exists():
        return None
    for line in env_path.read_text(encoding="utf-8").splitlines():
        line = line.strip()
        if not line or line.startswith("#") or "=" not in line:
            continue
        k, v = line.split("=", 1)
        if k.strip() == key:
            return v.strip()
    return None


def main() -> int:
    if len(sys.argv) != 3:
        print("Usage: add_admin_account.py <username_without_at> <password>")
        return 1

    username = sys.argv[1].strip()
    password = sys.argv[2].strip()

    if not username or "@" in username or not password:
        print("Invalid username or password. Pass only the local part (e.g. 'admin').")
        return 1

    repo_root = Path(__file__).resolve().parent.parent
    env_path = repo_root / "docker" / ".env"
    domain_base = load_env_value(env_path, "DOMAIN_BASE") or "office1789.com"
    email = f"{username}@{domain_base}"

    accounts_path = repo_root / "docker" / "config" / "postfix-accounts.cf"
    accounts_path.parent.mkdir(parents=True, exist_ok=True)

    try:
        hashed = crypt.crypt(password, crypt.mksalt(crypt.METHOD_SHA512))
    except Exception as e:
        print(f"Failed to hash password: {e}")
        return 1

    lines: list[str] = []
    if accounts_path.exists():
        for line in accounts_path.read_text(encoding="utf-8").splitlines():
            if line.split("|", 1)[0].strip() == email:
                continue
            if line.strip():
                lines.append(line.strip())
    lines.append(f"{email}|{hashed}")
    accounts_path.write_text("\n".join(lines) + "\n", encoding="utf-8")
    print(f"[OK] Synced admin mailbox {email} in {accounts_path}")
    print("Now restart mailserver: cd docker && docker compose restart mailserver")
    return 0


if __name__ == "__main__":
    sys.exit(main())
