.PHONY: up data-races down restart logs logs-backend logs-frontend ps stop clean create_mig init_mig exec_mig rollback_mig status_mig in re

DOCKER_COMPOSE_DEV = docker-config/docker-compose.dev.yml
DOCKER_COMPOSE_PROD = docker-config/docker-compose.prod.yml
DOCKER_COMPOSE = docker-config/docker-compose.yml
ENV_PROD = docker-config/.env.prod

up:
	./docker-config/update_ip.sh
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_DEV) up -d
	make print_posgres_ip

data-races:
	./docker-config/update_ip.sh
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_DEV) up -d
	make print_posgres_ip

prod:
	./docker-config/update_ip.sh
	docker compose --env-file $(ENV_PROD) -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_PROD) up -d

down:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_DEV) down

stop:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_DEV) stop

fclean:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_DEV) down -v --rmi all

restart:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_DEV) restart
	make print_posgres_ip

clean-images: clean
	@echo "Cleaning images..."
	docker image prune -af

clean-caches: clean clean-images
	@echo "Cleaning Go build cache..."
	go clean -cache -modcache
	@echo "Cleaning NPM cache..."
	npm cache clean --force

re: fclean
	make up

# --- Logs section ---
# Here you can call the service you want to see the logs of
# It is very important in when you are coding, to see the errors/bugs
# to use it, you have to : 
# make logs-golang or make logs-react, ect... 
logs-%:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_DEV) logs -f  $*

ps:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_DEV) ps

print_posgres_ip:
	docker inspect postgres_container | grep IPAddress

in-%:
	$(DOCKER) exec -it $* sh

# --- Database Migrations ---
# Migrations = version control for your database schema.
# Migrations allow us to directly modify the database, based on ORM Go functions
#   make create_mig       → Create a new migration file
#   make exec_mig          → Apply pending migrations
#   make rollback_mig      → Undo the last migration
#   make status_mig        → Show applied/pending migrations
create_mig:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_DEV) exec golang go run cmd/bun/main.go create_go ${NAME}

init_mig:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_DEV) exec golang go run cmd/bun/main.go init

exec_mig:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_DEV) exec golang go run cmd/bun/main.go migrate

rollback_mig:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_DEV) exec golang go run cmd/bun/main.go rollback

status_mig:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_DEV) exec golang go run cmd/bun/main.go status