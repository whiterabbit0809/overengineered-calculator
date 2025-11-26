# Dockerfile

# Build stage
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Run stage
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/server .
COPY web ./web

ENV PORT=8080

EXPOSE 8080

CMD ["./server"]
