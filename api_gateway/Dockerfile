# Stage 1: Builder
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway ./cmd/main.go

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api-gateway .

EXPOSE 8080

CMD ["./api-gateway"]
