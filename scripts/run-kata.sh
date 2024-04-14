#! /usr/bin/env bash

# This script runs a kata solution and compares the output with the expected output.

set -euo pipefail

kata=${1-}

if [ -z "$kata" ]; then
    echo "Usage: $0 <kata>" >&2
    exit 1
fi

self_dir=$(dirname "$0")
kata_dir=$(realpath "$self_dir/../katas")
kata_source="$kata_dir/$kata/main.go"
kata_in="$kata_dir/$kata/sample.in"
kata_out="$kata_dir/$kata/sample.out"

if ! command -v go &>/dev/null; then
    echo "go command not found" >&2
    exit 1
fi

if [ ! -f "$kata_source" ]; then
    echo "Kata source file not found: $kata_source" >&2
    exit 1
fi

if [ ! -f "$kata_in" ]; then
    echo "Kata input file not found: $kata_in" >&2
    exit 1
fi

if [ ! -f "$kata_out" ]; then
    echo "Kata output file not found: $kata_out" >&2
    exit 1
fi

program_out_buf=$(go run "$kata_source" <"$kata_in")
kata_out_buf=$(cat "$kata_out")

if [ "$program_out_buf" = "$kata_out_buf" ]; then
    echo "OK" >&2
    echo "==== Output ===="
    echo "$program_out_buf"
else
    echo "FAIL" >&2
    echo "==== Expected | Got ====" >&2
    diff -d -y --color=auto <(echo "$kata_out_buf") <(echo "$program_out_buf")
fi
