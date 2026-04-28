FROM golang:1.21-alpine AS builder

# Install build dependencies for sqlite3 (CGO)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Enable CGO for sqlite3
ENV CGO_ENABLED=1
RUN go build -o /history-service ./history/main.go

FROM alpine:latest
RUN apk add --no-cache libc6-compat
WORKDIR /app
COPY --from=builder /history-service .
EXPOSE 8082
CMD ["./history-service"]
