
FROM golang:1.25.6-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd

FROM alpine:latest

WORKDIR /

COPY --from=builder /app/main .

EXPOSE 30004

CMD ["./main"]

