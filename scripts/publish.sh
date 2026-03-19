#!/usr/bin/env sh

SERVER="${NATS_SERVER_URL:-nats://nats:4222}"
TYPE="${1:-weather}"

NOW=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

case "$TYPE" in
  weather)
    nats --server "$SERVER" pub "terrace.weather" "{\"temperature\": 26, \"is_raining\": \"true\", \"executed_at\": \"$NOW\"}"
    ;;
  log)

    ;;
  *)
    echo "Use: $0 [weather|log]"
    exit 1
    ;;
esac