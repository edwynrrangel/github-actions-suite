FROM golang:1.22.3-alpine3.19

WORKDIR /action

COPY src .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/read-yaml

ENTRYPOINT ["./bin/read-yaml"]