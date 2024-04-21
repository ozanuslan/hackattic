#! /usr/bin/env bash

export challenge_name="jotting_jwts"
export img_name="hackattic_jotting_jwts"
export container_name="jotting_jwts"

export ngrok_img_name="ngrok/ngrok:latest"
export ngrok_container_name="ngrok"

function challenge_exit() {
    docker kill "$ngrok_container_name" >/dev/null 2>&1 && echo "Killed $ngrok_container_name container" >&2 || echo "Failed to kill $ngrok_container_name container" >&2
    docker kill "$container_name" >/dev/null 2>&1 && echo "Killed $container_name container" >&2 || echo "Failed to kill $container_name container" >&2
}
