FROM golang:alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum from backend root
COPY go.mod go.sum ./
RUN go mod download

# Copy all backend source
COPY . .

# Build the specific service
RUN go build -o /interpolation-service ./interpolation/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /interpolation-service .
# Config will be mounted via volume or copied
EXPOSE 8081
CMD ["./interpolation-service"]
