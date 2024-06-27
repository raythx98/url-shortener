#Build Stage
FROM golang:1.22.4-bookworm AS builder
WORKDIR /app

RUN apt-get update && apt-get install

COPY go.mod go.sum ./
RUN go mod download

COPY . /app
RUN go build -o main cmd/api/main.go

# Run Stage
FROM debian:bookworm
WORKDIR /app

COPY --from=builder /app/main .

COPY .envrc .
EXPOSE 5051
CMD ["./main"]
