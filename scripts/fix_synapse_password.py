import os
import re
import subprocess
import sys
from pathlib import Path


def load_env(env_path: Path) -> dict:
    data: dict[str, str] = {}
    if not env_path.is_file():
        print(f"[ERR] {env_path} not found")
        sys.exit(1)
    for line in env_path.read_text(encoding="utf-8").splitlines():
        line = line.strip()
        if not line or line.startswith("#"):
            continue
        if "=" not in line:
            continue
        k, v = line.split("=", 1)
        data[k.strip()] = v.strip()
    return data


def run(cmd: list[str], cwd: Path) -> int:
    print(f"[CMD] {' '.join(cmd)}")
    return subprocess.run(cmd, cwd=str(cwd)).returncode


def fix_homeserver_yaml(docker_dir: Path, user: str, password: str, db_name: str) -> None:
    """
    Regenerate docker/synapse/conf/homeserver.yaml with a known good
    configuration, plugging in the DB credentials from docker/.env.

    This completely overwrites the file to avoid any lingering YAML
    corruption from previous edits.
    """
    yaml_path = docker_dir / "synapse" / "conf" / "homeserver.yaml"

    content = f"""server_name: "office1789.com"
report_stats: false

listeners:
  - port: 8008
    tls: false
    type: http
    x_forwarded: false
    bind_addresses: ['0.0.0.0']
    resources:
      - names: [client]
        compress: false

media_store_path: "/data/media_store"
signing_key_path: "/tmp/homeserver.signing.key"

enable_registration: false
registration_shared_secret: "OFFICE1789_REGISTRATION_SECRET"

# Configuration du système de mot de passe
password_config:
  enabled: true
  localdb_enabled: true
  pepper: ""

enable_set_displayname: true
enable_set_avatar_url: true
enable_3pid_changes: false

# PostgreSQL database configuration
database:
  name: psycopg2
  allow_unsafe_locale: true
  args:
    user: {user}
    password: {password}
    database: {db_name}
    host: postgres_synapse
    port: 5432
    cp_min: 5
    cp_max: 10

# TURN / VoIP configuration
turn_uris:
  - "turn:localhost:3478?transport=udp"
  - "turn:localhost:3478?transport=tcp"
turn_shared_secret: "MY_TURN_SECRET_1789"
turn_user_lifetime: 1h

# Rate limiting configuration - Augmenté pour SSO Office1789
rc_login:
  address:
    per_second: 10
    burst_count: 20
  account:
    per_second: 5
    burst_count: 10
  failed_attempts:
    per_second: 1
    burst_count: 5
"""

    yaml_path.write_text(content, encoding="utf-8")
    print(f"[INFO] Regenerated {yaml_path} with a clean PostgreSQL config matching docker/.env")


def main() -> None:
    script_dir = Path(__file__).resolve().parent
    repo_root = script_dir.parent
    docker_dir = repo_root / "docker"
    env_file = docker_dir / ".env"

    env = load_env(env_file)

    user = env.get("SYNAPSE_DB_USER") or "synapse"
    password = env.get("SYNAPSE_DB_PASSWORD")
    db_name = env.get("SYNAPSE_DB_NAME") or "synapse"

    if not password:
        print("[ERR] SYNAPSE_DB_PASSWORD is not set in docker/.env; nothing to fix.")
        sys.exit(1)

    # Ensure homeserver.yaml database block matches docker/.env
    fix_homeserver_yaml(docker_dir, user, password, db_name)

    # Align Postgres role password
    print(
        f"[INFO] Will align Postgres user '{user}' password "
        f"with value from docker/.env for database '{db_name}'."
    )

    # Escape single quotes in password for SQL
    pw_sql = password.replace("'", "''")
    sql = f"ALTER USER {user} WITH PASSWORD '{pw_sql}';"

    # Execute inside postgres_synapse container.
    # Note: for this service, POSTGRES_USER is SYNAPSE_DB_USER, so `user` is superuser.
    cmd = [
        "docker",
        "compose",
        "exec",
        "-T",
        "postgres_synapse",
        "psql",
        "-U",
        user,
        "-d",
        db_name,
        "-c",
        sql,
    ]

    code = run(cmd, cwd=docker_dir)
    if code != 0:
        print("[ERR] Failed to run ALTER USER in postgres_synapse.")
        sys.exit(code)

    print("[OK] Synapse DB user password updated. You can now restart synapse:")
    print("     cd docker && docker compose up -d synapse")


if __name__ == "__main__":
    main()
