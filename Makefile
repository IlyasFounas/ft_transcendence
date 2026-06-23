.PHONY: up down restart logs ps

up:
	mkdir -p docker-config/data/postgres
	docker compose -f docker-config/docker-compose.yml up -d

down:
	docker compose -f docker-config/docker-compose.yml down

restart:
	docker compose -f docker-config/docker-compose.yml restart

logs:
	docker compose -f docker-config/docker-compose.yml logs -f

ps:
	docker compose -f docker-config/docker-compose.yml ps