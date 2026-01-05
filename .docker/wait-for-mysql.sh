#!/usr/bin/env bash
set -euo pipefail

# Usage: wait-for-mysql.sh [host] [port] [timeout]
# Defaults: host=mysql port=3306 timeout=60

HOST=${1:-mysql}
PORT=${2:-3306}
TIMEOUT=${3:-60}

echo "Waiting for MySQL at ${HOST}:${PORT} (timeout ${TIMEOUT}s)..."

start_ts=$(date +%s)
while true; do
  # Try mysqladmin if available
  if command -v mysqladmin >/dev/null 2>&1; then
    if mysqladmin ping -h "${HOST}" -P "${PORT}" --silent; then
      echo "MySQL is available"
      exit 0
    fi
  else
    # Fallback to TCP check using /dev/tcp (requires bash)
    if (echo > /dev/tcp/${HOST}/${PORT}) >/dev/null 2>&1; then
      echo "MySQL TCP port is open"
      exit 0
    fi
  fi

  now=$(date +%s)
  elapsed=$((now - start_ts))
  if [ "${elapsed}" -ge "${TIMEOUT}" ]; then
    echo "Timeout after ${TIMEOUT}s waiting for ${HOST}:${PORT}"
    exit 1
  fi
  sleep 1
done
