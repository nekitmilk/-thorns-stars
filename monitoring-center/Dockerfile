FROM golang:1.24.5 AS builder

WORKDIR /app

# Копируем файлы модулей
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем весь код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o monitoring-center ./cmd/main.go

# Финальный образ
FROM alpine:latest

WORKDIR /app

# Копируем бинарник
COPY --from=builder /app/monitoring-center .

RUN adduser -D -u 1000 appuser && \
    chown appuser:appuser /app/monitoring-center

USER appuser

EXPOSE 8080

CMD ["./monitoring-center"]