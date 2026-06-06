COMPOSE = docker compose -f backend/Docker-Compose.yaml

.PHONY: build down

# Build (or rebuild) all service images and start the stack.
build:
	$(COMPOSE) up --build

# Stop and remove containers, networks, and named volumes (-v wipes the DB).
down:
	$(COMPOSE) down -v
