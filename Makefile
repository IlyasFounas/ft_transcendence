.PHONY: up down restart logs logs-backend logs-frontend ps stop clean

DOCKER_COMPOSE = docker-config/docker-compose.yml
DOCKER_COMPOSE_OVERRIDE = docker-config/docker-compose.override.yml

up:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_OVERRIDE) up -d

down:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_OVERRIDE) down

stop:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_OVERRIDE) stop

fclean:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_OVERRIDE) down -v --rmi all

restart:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_OVERRIDE) restart

logs:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_OVERRIDE) logs -f

logs-backend:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_OVERRIDE) logs -f golang

logs-frontend:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_OVERRIDE) logs -f react

ps:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_OVERRIDE) ps

# --- Database Migrations ---
# Migrations = version control for your database schema.
# Migrations allow us to directly modify the database, based on ORM Go functions
#   make create_mig       → Create a new migration file
#   make exec_mig          → Apply pending migrations
#   make rollback_mig      → Undo the last migration
#   make status_mig        → Show applied/pending migrations
create_mig:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_OVERRIDE) exec golang go run cmd/bun/main.go db create_go

exec_mig:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_OVERRIDE) exec golang go run cmd/bun/main.go db migrate

rollback_mig:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_OVERRIDE) exec golang go run cmd/bun/main.go db rollback

status_mig:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_OVERRIDE) exec golang go run cmd/bun/main.go db status
	