#! /usr/bin/env bash

set -euo pipefail

if ! command -v docker &>/dev/null; then
    echo "docker command not found" >&2
    exit 1
fi

self_dir=$(dirname "$0")
source "$self_dir/common.sh"

docker run -d \
    --rm \
    --net=host \
    --name "$ngrok_container_name" \
    -e NGROK_AUTHTOKEN="$NGROK_AUTH_TOKEN" \
    "$ngrok_img_name" http --domain="$NGROK_DOMAIN" 1337 >/dev/null

if [[ -z "$(docker images -q "$img_name")" || -n "${CLEAN_BUILD:-}" ]]; then
    docker build -t "$img_name" . >&2
fi

echo -n "{\"app_url\":\"https://$NGROK_DOMAIN\"}"

exec $(
    docker run -i \
        --rm \
        -p 1337:1337 \
        --name "$container_name" \
        "$img_name" <&0 >/dev/null 2> >(while read -r line; do echo "[container] $line" >&2; done) &
    sleep 1
)
