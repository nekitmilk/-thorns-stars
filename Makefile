.PHONY: up up-build down logs clean purge register-host wait-for-center status help

# Variables
DOCKER_COMPOSE = docker-compose -f deployments/docker-compose.yml

# Wait for monitoring center to be ready
wait-for-center:
	@echo "Waiting for Monitoring Center to start..."
	@for i in 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20; do \
		if curl -s http://localhost:8080/api/hosts >/dev/null; then \
			echo "Monitoring Center is ready!"; \
			exit 0; \
		fi; \
		echo "Still waiting for Monitoring Center... ($$i/20)"; \
		sleep 3; \
	done; \
	echo "Monitoring Center did not start in time"; \
	exit 1
# Register host in monitoring center
register-host: wait-for-center
	@echo "Registering host in Monitoring Center..."
	@chmod +x scripts/register-host.sh
	@./scripts/register-host.sh

# Start all services (without agent first)
start-infrastructure:
	@echo "Starting infrastructure services (Postgres, MongoDB, Monitoring Center)..."
	$(DOCKER_COMPOSE) up -d postgres mongodb monitoring-center

# Start agent only
start-agent:
	@echo "Starting monitoring agent..."
	$(DOCKER_COMPOSE) up -d monitoring-agent

# Start all services with proper order
up: start-infrastructure register-host start-agent
	@echo "All services started successfully!"

# Build and start with proper order
up-build: 
	@echo "Building and starting infrastructure..."
	$(DOCKER_COMPOSE) up -d --build postgres mongodb monitoring-center
	@make register-host
	@echo "Building and starting agent..."
	$(DOCKER_COMPOSE) up -d --build monitoring-agent
	@echo "All services built and started successfully!"

# Stop services
down:
	@echo "Stopping services..."
	$(DOCKER_COMPOSE) down

# View logs
logs:
	$(DOCKER_COMPOSE) logs -f

# Clean up
clean: down
	@echo "Cleaning up..."

# Full purge (including volumes)
purge: down
	@echo "Purging everything including volumes..."
	$(DOCKER_COMPOSE) down -v --remove-orphans

# Build individual components
build-agent:
	docker build -t monitoring-agent -f agent/Dockerfile agent/

build-center:
	docker build -t monitoring-center -f monitoring-center/Dockerfile monitoring-center/

# Status check
status:
	$(DOCKER_COMPOSE) ps

# Help
help:
	@echo "Available commands:"
	@echo "  make up           - Start all services with proper order"
	@echo "  make up-build     - Build and start with proper order"
	@echo "  make down         - Stop services"
	@echo "  make logs         - View logs"
	@echo "  make purge        - Stop and remove everything including volumes"
	@echo "  make register-host - Register host only (after infrastructure is up)"
	@echo "  make status       - Show container status"

.DEFAULT_GOAL := help