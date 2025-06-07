.PHONY:
start:
	@echo "Starting Docker-Postgres"
	docker compose up -d
	@echo "Starting backend server"
	air

.PHONY: stop
stop:
	@echo "Stopping Docker-Postgres"
	docker compose down

