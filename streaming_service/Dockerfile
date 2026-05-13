FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o streaming-app ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/streaming-app .

RUN mkdir -p /app/audio_files

EXPOSE 50052

CMD ["./streaming-app"]