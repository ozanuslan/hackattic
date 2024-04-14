#! /usr/bin/env bash

# This script verifies a challenge solution by sending the output to the Hackattic API.

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
challenge_out="$challenge_dir/challenge.out"

if [ ! -d "$challenge_dir" ]; then
    echo "Challenge directory not found: $challenge_dir" >&2
    exit 1
fi

if [ ! -f "$challenge_out" ]; then
    echo "Challenge output file not found: $challenge_out" >&2
    exit 1
fi

. "$env_path"

challenge_solve_url="https://hackattic.com/challenges/$challenge/solve?access_token=$ACCESS_TOKEN"

if ! command -v curl &>/dev/null; then
    echo "curl command not found" >&2
    exit 1
fi

response=$(curl -s -w "%{http_code}" -X POST -d @"$challenge_out" -H "Content-Type: application/json" "$challenge_solve_url")
response_code="${response: -3}"
response_body="${response:0:${#response}-3}"

echo "HTTP $response_code" >&2
if command -v jq &>/dev/null; then
    echo "$response_body" | jq
else
    echo "$response_body"
fi
