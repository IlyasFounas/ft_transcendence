.PHONY: up down restart logs logs-backend logs-frontend ps stop clean init_mig create_mig exec_mig rollback_mig status_mig

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

clean:
	rm -rf dist/ build/ main #put the right folders/files

clean-images: clean
	@echo "Cleaning images..."
	docker image prune -af

clean-caches: clean clean-images
	@echo "Cleaning Go build cache..."
	go clean -cache -modcache
	@echo "Cleaning NPM cache..."
	npm cache clean --force

# --- Logs section ---
# Here you can call the service you want to see the logs of
# It is very important in when you are coding, to see the errors/bugs
# to use it, you have to : 
# make logs-golang or make logs-react, ect... 
logs:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_OVERRIDE) logs -f

logs-%:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_OVERRIDE) logs -f  $*

ps:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_OVERRIDE) ps

print_posgres_ip:
	docker inspect postgres_container | grep IPAddress

# --- Database Migrations ---
# Migrations = version control for your database schema.
# Migrations allow us to directly modify the database, based on ORM Go functions
#   make create_mig       → Create a new migration file
#   make exec_mig          → Apply pending migrations
#   make rollback_mig      → Undo the last migration
#   make status_mig        → Show applied/pending migrations
	
init_mig:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_OVERRIDE) exec golang go run cmd/bun/main.go init
create_mig:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_OVERRIDE) exec golang go run cmd/bun/main.go create_go ${NAME}

exec_mig:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_OVERRIDE) exec golang go run cmd/bun/main.go migrate

rollback_mig:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_OVERRIDE) exec golang go run cmd/bun/main.go rollback

status_mig:
	docker compose -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_OVERRIDE) exec golang go run cmd/bun/main.go status