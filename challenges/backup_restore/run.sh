#! /bin/bash

set -euo pipefail

img_name="hackattic_backup_restore"

# if the image does not exist, build it
if ! docker image inspect "$img_name" &>/dev/null; then
    docker build -t "$img_name" .
fi

docker run -i --rm "$img_name" </dev/stdin
