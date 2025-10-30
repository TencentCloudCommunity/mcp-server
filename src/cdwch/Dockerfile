FROM golang:1.24 as builder

ENV GO111MODULE=on

WORKDIR /app

COPY . .

RUN go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go

FROM alpine

COPY --from=builder /app/main /app/main

CMD ["/app/main"]