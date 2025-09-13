#!/usr/bin/env sh
# wait-for-db.sh
# Simple wait-for-postgres script useful in containers

set -e

host="${DB_HOST:-db}"
port="${DB_PORT:-5432}"
user="${DB_USER:-postgres}"
password="${DB_PASSWORD:-postgres}"
db="${DB_NAME:-mergemoney}"
timeout=${DB_TIMEOUT:-30}

echo "Waiting for postgres at ${host}:${port} (db=${db})..."

export PGPASSWORD="$password"

tries=0
while ! pg_isready -h "$host" -p "$port" -U "$user" -d "$db" >/dev/null 2>&1; do
  tries=$((tries + 1))
  if [ "$tries" -ge "$timeout" ]; then
    echo "Timed out waiting for Postgres after ${timeout} seconds"
    exit 1
  fi
  sleep 1
done

echo "Postgres is available"

exec "$@"
