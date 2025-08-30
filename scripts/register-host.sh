#!/bin/bash

# Возможно имеет смысл добавить хост в миграцию и убрать этот скрипт
# Load environment variables
set -a
# source ../deployments/.env
source deployments/.env
set +a

# Check if monitoring center is available
echo "Waiting for Monitoring Center to be ready..."
until curl -s http://localhost:8080/api/hosts >/dev/null; do
    sleep 2
done

# Check if host already exists
echo "Checking if host $HOST_ID already exists..."
response=$(curl -s -o /dev/null -w "%{http_code}" \
  "http://localhost:8080/api/hosts/$HOST_ID")

if [ "$response" -eq 200 ]; then
    echo "Host $HOST_ID already exists. Skipping registration."
    exit 0
fi

# Register new host
echo "Registering host $HOST_NAME with ID $HOST_ID..."
response=$(curl -s -w "%{http_code}" -X POST http://localhost:8080/api/hosts \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"$HOST_NAME\",
    \"ip\": \"$HOST_IP\",
    \"priority\": $HOST_PRIORITY
  }")

status_code=${response: -3}
response_body=${response%???}

if [ "$status_code" -eq 201 ]; then
    echo "Host registered successfully: $response_body"
elif [ "$status_code" -eq 400 ]; then
    echo "Host already exists or invalid data: $response_body"
else
    echo "Failed to register host. Status: $status_code, Response: $response_body"
    exit 1
fi