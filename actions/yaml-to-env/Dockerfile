FROM golang:1.22.3-alpine3.19 AS builder

WORKDIR /app

COPY src .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/read-yaml main.go

FROM alpine:3.19 AS publisher

COPY --from=builder /app/bin/read-yaml /usr/local/bin/read-yaml

ENTRYPOINT ["read-yaml"]