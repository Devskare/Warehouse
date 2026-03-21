FROM golang:1.25.6-bookworm AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o main ./cmd

FROM alpine:latest

COPY --from=builder /app/main /main
COPY --from=builder /app/.env /main

EXPOSE 30004

CMD ["/main"]


