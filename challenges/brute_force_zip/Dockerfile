FROM docker.io/library/golang:1.22.1-bullseye AS builder

WORKDIR /build
COPY . .
# RUN go get -d -v
RUN go build -o /go/bin/brute_force_zip

FROM docker.io/library/debian:bullseye-slim

RUN apt-get update && apt-get install -y fcrackzip ca-certificates unzip && rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/bin/brute_force_zip /usr/local/bin/brute_force_zip

WORKDIR /brute_force_zip

ENTRYPOINT [ "/usr/local/bin/brute_force_zip" ]