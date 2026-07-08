FROM golang:1.24 AS builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/postgres-server .

FROM alpine:3.22

RUN apk add --no-cache ca-certificates && \
    addgroup -S app && adduser -S -G app -h /app app && \
    mkdir -p /app/data && chown -R app:app /app

WORKDIR /app
COPY --from=builder /out/postgres-server /app/postgres-server

USER app

CMD ["/app/postgres-server"]
