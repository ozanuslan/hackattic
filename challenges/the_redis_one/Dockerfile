FROM docker.io/library/golang:1.22.1 AS builder

WORKDIR /build
COPY . .
RUN go get -d -v
RUN go build -o /go/bin/the_redis_one

FROM docker.io/library/debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    redis-server \
    jq \
    curl \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

RUN curl -sL https://golang.org/dl/go1.22.1.linux-$(dpkg --print-architecture).tar.gz | tar -C /usr/local -xzf - && \
    export PATH=$PATH:/usr/local/go/bin && \
    go env -w GOBIN=/usr/local/bin && \
    go install github.com/hdt3213/rdb@latest

COPY --from=builder /go/bin/the_redis_one /usr/local/bin/the_redis_one

COPY ./entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]