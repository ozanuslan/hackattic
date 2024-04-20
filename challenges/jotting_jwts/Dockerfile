FROM docker.io/library/golang:1.22.1-bullseye AS builder

WORKDIR /build
COPY . .
RUN go get -d -v
RUN go build -o /go/bin/jotting_jwts

FROM docker.io/library/debian:bullseye-slim

COPY --from=builder /go/bin/jotting_jwts /usr/local/bin/jotting_jwts

ENTRYPOINT [ "/usr/local/bin/jotting_jwts" ]
