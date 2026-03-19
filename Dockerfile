FROM ubuntu:latest
LABEL authors="daniloo"

ENTRYPOINT ["top", "-b"]

FROM golang:1.25.6-bookworm

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o /app/exe main.go

CMD ["/app/exe"]