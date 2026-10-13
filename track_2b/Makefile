.PHONY: run stop clean

COMPOSE = docker compose -f compose.run.yaml

# Builds and starts PostgreSQL + LaporSiap in Docker from a clean checkout.
# Open http://localhost:8787 once the app reports ready.
run:
	$(COMPOSE) up --build

stop:
	$(COMPOSE) down

# Also removes the demo database volume.
clean:
	$(COMPOSE) down --volumes
