#! /usr/bin/env bash

# This script runs a challenge solution and saves the output to the challenge output file.
# It can also verify the challenge output by running the verify-challenge.sh script.
# If you want to verify the challenge after you've passed, you can use the playground option.

set -euo pipefail

verify=false
playground=false
new_input=false

while [ $# -gt 0 ]; do
    case "$1" in
    --verify | -v)
        verify=true
        ;;
    --playground | -p)
        playground=true
        ;;
    --new-input | -n)
        new_input=true
        ;;
    *)
        break
        ;;
    esac
    shift
done

challenge=${1-}

if [ -z "$challenge" ]; then
    echo "Usage: $0 [flags] <challenge>" >&2
    echo "Options: --verify,-v        Verify the challenge output" >&2
    echo "         --playground,-p    Run the challenge in playground mode" >&2
    echo "         --new-input,-n     Get new challenge input" >&2
    exit 1
fi

self_dir=$(dirname "$0")
env_path="$self_dir/../.env"
challenges_dir=$(realpath "$self_dir/../challenges")
challenge_dir="$challenges_dir/$challenge"
challenge_source="$challenge_dir/main.go"
challenge_runner="$challenge_dir/run.sh"
challenge_in="$challenge_dir/challenge.in"
challenge_out="$challenge_dir/challenge.out"

. "$env_path"
export ACCESS_TOKEN

if [ ! -d "$challenge_dir" ]; then
    echo "Challenge directory not found: $challenge_dir" >&2
    exit 1
fi

if [[ ! -f "$challenge_in" && "$new_input" != true ]]; then
    echo "Challenge input file not found: $challenge_in" >&2
    exit 1
fi

if [[ "$verify" = true || "$new_input" = true ]]; then
    get_input_script="$self_dir/get-challenge-input.sh"
    "$get_input_script" "$challenge" >&2
fi

if [ ! -f "$challenge_runner" ]; then
    program_out_buf=$(cd "$challenge_dir" && go run "$challenge_source" <"$challenge_in")
else
    program_out_buf=$(cd "$challenge_dir" && bash "$challenge_runner" <"$challenge_in")
fi

echo "$program_out_buf" >"$challenge_out"
echo "Output saved to: $challenge_out" >&2

if [ "$verify" = true ]; then
    verify_script="$self_dir/verify-challenge.sh"
    if [ "$playground" = true ]; then
        "$verify_script" --playground "$challenge"
    else
        "$verify_script" "$challenge"
    fi
else
    echo "==== Output ====" >&2
    echo "$program_out_buf"
fi
