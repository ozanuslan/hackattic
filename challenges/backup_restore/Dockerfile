FROM docker.io/library/golang:1.22.1-bullseye AS builder

WORKDIR /build
COPY . .
RUN go get -d -v
RUN go build -o /go/bin/backup_restore

FROM docker.io/library/debian:bullseye-slim

RUN apt-get update && apt-get install -y curl gnupg
RUN curl -sSL "https://www.postgresql.org/media/keys/ACCC4CF8.asc" | apt-key add -
RUN sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt/ bullseye-pgdg main" >> /etc/apt/sources.list.d/pgdg.list'
RUN apt-get update && apt-get install -y postgresql-10

COPY --from=builder /go/bin/backup_restore /usr/local/bin/backup_restore

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

USER postgres

ENTRYPOINT [ "./entrypoint.sh" ]