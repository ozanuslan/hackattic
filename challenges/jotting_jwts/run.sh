#! /usr/bin/env bash

set -euo pipefail

challenge_name="jotting_jwts"
img_name="hackattic_jotting_jwts"
container_name="jotting_jwts"

ngrok_img_name="ngrok/ngrok:latest"
ngrok_container_name="ngrok"

if ! command -v docker &>/dev/null; then
    echo "docker command not found" >&2
    exit 1
fi

function challenge_exit() {
    docker kill "$ngrok_container_name" >/dev/null 2>&1 || true
    docker kill "$container_name" >/dev/null 2>&1 || true
}
trap challenge_exit SIGTERM SIGINT ERR EXIT

docker run -d --net=host --rm --name "$ngrok_container_name" -e NGROK_AUTHTOKEN="$NGROK_AUTH_TOKEN" "$ngrok_img_name" http --domain="$NGROK_DOMAIN" 1337

if [[ -z "$(docker images -q "$img_name")" || -n "${CLEAN_BUILD:-}" ]]; then
    docker build -t "$img_name" .
fi

docker run -i -p 1337:1337 --rm --name "$container_name" "$img_name" </dev/stdin &
sleep 1

echo "{\"app_url\":\"https://$NGROK_DOMAIN\"}" >"$CHALLENGE_OUT"
"$VERIFY_SCRIPT" "$challenge_name" >&2

exit 254
