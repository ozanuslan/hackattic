#! /usr/bin/env bash

set -eu

cat <<CONFIG >/etc/postgresql/10/main/pg_hba.conf
# TYPE  DATABASE        USER            ADDRESS             METHOD
# "local" is for Unix domain socket connections only
local   all             all                                 trust

# IPv4 local connections:
host    all             all             0.0.0.0/32          md5
host    all             all             127.0.0.1/32        trust
host    all             all             all                 trust

# IPv6 local connections:
host    all             all             ::1/128             trust

# Allow replication connections from localhost, by a user with the
# replication privilege.
#local   replication     postgres                                peer
#host    replication     postgres        127.0.0.1/32            ident
#host    replication     postgres        ::1/128                 ident
CONFIG

service postgresql restart >&2

/usr/local/bin/backup_restore decode </dev/stdin | psql >&2

export PGUSER=postgres
export PGHOST=localhost
/usr/local/bin/backup_restore get-ssns
