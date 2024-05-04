#! /usr/bin/env bash

set -eu

redis_host=127.0.0.1
redis_port=6379
redis_dir=/var/lib/redis
redis_dbfilename=dump.rdb

cat <<CONFIG >/etc/redis/redis.conf
bind $redis_host
port $redis_port
dir $redis_dir
dbfilename $redis_dbfilename
appendonly no
protected-mode no
CONFIG

stdin=$(</dev/stdin)

/usr/local/bin/the_redis_one decode < <(echo "$stdin") >"$redis_dir/$redis_dbfilename"

chown redis: /etc/redis/redis.conf "$redis_dir/$redis_dbfilename"

redis-server /etc/redis/redis.conf > >(while read -r line; do echo "[redis] $line" >&2; done) &
sleep 1 # wait for possible slow startup

expr_timestamp=$(/usr/local/bin/rdb -c json -o expr.rdb "$redis_dir/$redis_dbfilename" && cat expr.rdb | jq -r '.[] | select(.expiration) | .expiration' | head -n1 | xargs -I{} date -d {} +%s%3N && rm expr.rdb)
echo "Extracted expiry timestamp: $expr_timestamp" >&2

export REDIS_HOST=$redis_host
export REDIS_PORT=$redis_port
export REDIS_USER=redis
export REDIS_PASSWORD=
export EXPIRY_TIMESTAMP=$expr_timestamp

/usr/local/bin/the_redis_one read-redis < <(echo "$stdin")
