# Этап 1: Сборка (Builder)
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o catalog-app ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/catalog-app .

EXPOSE 50053

CMD ["./catalog-app"]