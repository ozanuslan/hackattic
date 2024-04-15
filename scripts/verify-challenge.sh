#! /usr/bin/env bash

# This script verifies a challenge solution by sending the output to the Hackattic API.
# You can use the playground option to verify the challenge after you've passed.

set -euo pipefail

challenge=${1-}

playground=false
if [ "$challenge" = "--playground" ] || [ "$challenge" = "-p" ]; then
    playground=true
    challenge=${2-}
fi

if [ -z "$challenge" ] || [ "$challenge" = "--help" ] || [ "$challenge" = "-h" ]; then
    echo "Usage: $0 [flags] <challenge>" >&2
    echo "Options: --playground,-p    Verify the challenge in playground mode" >&2
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

if [ "$playground" = true ]; then
    challenge_solve_url="$challenge_solve_url&playground=1"
fi

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
