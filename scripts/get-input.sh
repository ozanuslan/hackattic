#! /usr/bin/env bash

set -euo pipefail

challenge=${1-}

if [ -z "$challenge" ]; then
    echo "Usage: $0 <challenge>" >&2
    exit 1
fi

self_dir=$(dirname "$0")
env_path="$self_dir/../.env"
challenges_dir=$(realpath "$self_dir/../challenges")
challenge_dir="$challenges_dir/$challenge"
challenge_in="$challenge_dir/challenge.in"

if [ ! -d "$challenge_dir" ]; then
    echo "Challenge directory not found: $challenge_dir" >&2
    exit 1
fi

. "$env_path"

challenge_input_url="https://hackattic.com/challenges/$challenge/problem?access_token=$ACCESS_TOKEN"

if ! command -v curl &>/dev/null; then
    echo "curl command not found" >&2
    exit 1
fi

curl -s "$challenge_input_url" >"$challenge_in"

echo "Input saved to: $challenge_in"
