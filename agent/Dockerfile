FROM golang:1.24.5 AS builder

WORKDIR /app

# Копируем файлы модулей и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарный файл
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o monitoring-agent ./cmd/main.go

# Финальный этап
FROM alpine:latest

WORKDIR /app

# Устанавливаем зависимости для сбора метрик
RUN apk add --no-cache \
    ca-certificates \
    && update-ca-certificates

# Копируем бинарник из builder этапа
COPY --from=builder /app/monitoring-agent .

# Создаем пользователя для безопасности
RUN adduser -D -u 1000 agentuser && \
    chown agentuser:agentuser /app/monitoring-agent

USER agentuser

# Объемы для доступа к системным метрикам
VOLUME ["/host"]

# Переменные окружения
# HOST_ID
ENV MONITORING_CENTER_URL="http://localhost:8080"
ENV POLLING_INTERVAL="5m"
ENV REQUEST_TIMEOUT="30s"

# Запускаем агент
CMD ["./monitoring-agent"]