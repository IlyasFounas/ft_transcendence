.PHONY: up down restart logs ps

DOCKER_COMPOSE= docker-config/docker-compose.yml
DOCKER_COMPOSE_OVERRIDE= docker-config/docker-compose.override.yml

up:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_OVERRIDE) up -d

down:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_OVERRIDE) down

restart:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_OVERRIDE) restart

logs:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_OVERRIDE) logs -f

ps:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_OVERRIDE) ps