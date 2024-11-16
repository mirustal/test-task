# сборка
FROM golang:1.22-alpine AS builder
ENV CGO_ENABLED=0 TZ=Europe/Moscow

RUN apk add --no-cache git tzdata

# копирование зависимостей
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

# сборка
COPY . .
RUN go build -o /app/main cmd/bank/main.go

# Финальный образ
FROM alpine:3.20
ENV TZ=Europe/Moscow

RUN apk add --no-cache tzdata && \ 
    rm -rf /var/cache/apk/*

WORKDIR /app

COPY --from=builder /app/main .
COPY .env .
COPY config config

CMD ["./main"]