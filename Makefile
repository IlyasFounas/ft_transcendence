.PHONY: up data-races down restart logs ps stop fclean re prod prod-down prod-stop prod-fclean prod-ps create_mig init_mig exec_mig rollback_mig status_mig

DOCKER_COMPOSE      = docker-config/docker-compose.yml
DOCKER_COMPOSE_DEV  = docker-config/docker-compose.dev.yml
DOCKER_COMPOSE_PROD = docker-config/docker-compose.prod.yml
ENV_DEV             = docker-config/.env
ENV_PROD            = docker-config/.env.prod

DC_DEV  = docker compose --env-file $(ENV_DEV) -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_DEV)
DC_PROD = docker compose --env-file $(ENV_PROD) -f $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_PROD)

# --- Dev ---
up:
	$(DC_DEV) up -d
	$(MAKE) print_posgres_ip

data-races:
	$(DC_DEV) up -d
	$(MAKE) print_posgres_ip

down:
	$(DC_DEV) down

stop:
	$(DC_DEV) stop

fclean:
	$(DC_DEV) down -v --rmi all

restart:
	$(DC_DEV) up -d --force-recreate
	$(MAKE) print_posgres_ip

re: fclean
	$(MAKE) up

ps:
	$(DC_DEV) ps

logs-%:
	$(DC_DEV) logs -f $*

# --- Prod ---
prod:
	./docker-config/update_ip.sh $(ENV_PROD)
	$(DC_PROD) up -d --build

prod-down:
	$(DC_PROD) down

prod-stop:
	$(DC_PROD) stop

prod-fclean:
	$(DC_PROD) down -v --rmi all

prod-ps:
	$(DC_PROD) ps

prod-logs-%:
	$(DC_PROD) logs -f $*

# --- Utils ---
clean-images:
	@echo "Cleaning images..."
	docker image prune -af

clean-caches: clean-images
	@echo "Cleaning Go build cache..."
	go clean -cache -modcache
	@echo "Cleaning NPM cache..."
	npm cache clean --force

print_posgres_ip:
	docker inspect postgres_container | grep IPAddress

in-%:
	docker exec -it $* sh

# --- Database Migrations ---
create_mig:
	$(DC_DEV) exec golang go run cmd/bun/main.go create_go ${NAME}

init_mig:
	$(DC_DEV) exec golang go run cmd/bun/main.go init

exec_mig:
	$(DC_DEV) exec golang go run cmd/bun/main.go migrate

rollback_mig:
	$(DC_DEV) exec golang go run cmd/bun/main.go rollback

status_mig:
	$(DC_DEV) exec golang go run cmd/bun/main.go status