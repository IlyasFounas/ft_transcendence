.PHONY: up down restart logs ps

DOCKER_COMPOSE= docker-config/docker-compose.yml

up:
	docker compose -f $(DOCKER_COMPOSE) up -d

down:
	docker compose -f $(DOCKER_COMPOSE) down

restart:
	docker compose -f $(DOCKER_COMPOSE) restart

logs:
	docker compose -f $(DOCKER_COMPOSE) logs -f

ps:
	docker compose -f $(DOCKER_COMPOSE) ps