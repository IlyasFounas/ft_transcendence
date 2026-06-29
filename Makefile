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