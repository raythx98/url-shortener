#Build Stage
FROM golang:1.24.0-bookworm AS builder

RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . /app
RUN go build -o main cmd/api/main.go

# Run Stage
FROM debian:bookworm

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs

EXPOSE 5051
CMD ["./main"]
