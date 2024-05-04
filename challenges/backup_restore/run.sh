#! /usr/bin/env bash

# This script is used to run the backup_restore challenge.
# You can set the CLEAN_BUILD environment variable to force a clean image build.

set -euo pipefail

challenge="backup_restore"
img_name="hackattic_backup_restore"

if ! command -v docker &>/dev/null; then
    echo "docker command not found" >&2
    exit 1
fi

if [[ -z "$(docker images -q "$img_name")" || -n "${CLEAN_BUILD:-}" ]]; then
    docker image prune --force --filter="label=$img_name" >&2
    docker build --label "$img_name" -t "$img_name" . >&2
fi

docker run -i --rm --name="$challenge" "$img_name" <&0 2> >(while read -r line; do echo "[container] $line" >&2; done)
