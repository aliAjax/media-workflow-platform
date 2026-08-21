#!/bin/sh
set -eu
base=${MEDIA_BASE_URL:-http://127.0.0.1:8084}
curl -fsS "$base/healthz" >/dev/null
curl -fsS -X POST "$base/api/v1/assets" -H 'content-type: application/json' -d '{"tenant_id":"demo","name":"sample.mp4","content_type":"video/mp4","size":128,"checksum":"abc"}' | grep -q 'tenant_id'
echo "media smoke ok"
