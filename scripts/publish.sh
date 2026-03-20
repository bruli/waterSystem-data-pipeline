#!/bin/sh

SERVER="${NATS_SERVER_URL:-nats://nats:4222}"
STREAM_NAME="${STREAM_NAME:-EVENTS}"
TYPE="${1:-weather}"

NOW=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
SEED=$(date +%s)

TEMP=$(awk -v min=15 -v max=35 -v seed="$SEED" 'BEGIN{srand(seed); printf "%.2f", min+rand()*(max-min)}')
SECONDS_VAL=$(awk -v min=5 -v max=60 -v seed="$((SEED+1))" 'BEGIN{srand(seed); print int(min+rand()*(max-min+1))}')
RAIN_FLAG=$(awk -v seed="$((SEED+2))" 'BEGIN{srand(seed); if (rand() < 0.2) print "true"; else print "false"}')

case "$TYPE" in
  weather)
    PAYLOAD="{\"temperature\": $TEMP, \"is_raining\": $RAIN_FLAG, \"executed_at\": \"$NOW\"}"
    echo "$PAYLOAD"
    nats --server "$SERVER" pub "terrace_weather" "$PAYLOAD"
    ;;
  log)
    PAYLOAD="{\"seconds\": $SECONDS_VAL, \"zone\": \"zone_testing\", \"executed_at\": \"$NOW\"}"
    echo "$PAYLOAD"
    nats --server "$SERVER" pub "execution_logs" "$PAYLOAD"
    ;;
  clean)
    echo "Cleaning stream $STREAM_NAME..."
    nats --server "$SERVER" stream rm "$STREAM_NAME" -f
    ;;
  purge)
    echo "Purging stream $STREAM_NAME..."
    nats --server "$SERVER" stream purge "$STREAM_NAME"
    ;;
  *)
    echo "Use: $0 [weather|log|clean|purge]"
    exit 1
    ;;
esac