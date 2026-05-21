FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o auth-app ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/auth-app .
COPY --from=builder /app/migrations ./migrations

EXPOSE 50051

CMD ["./auth-app"]