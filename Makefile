GOOSE := $(shell which goose)
DOCKER_COMPOSE := $(shell which docker-compose)
ifndef DOCKER_COMPOSE
		$(error "docker compose is not installed")
endif
ifndef GOOSE
	$(error "goose is not installed")
endif

.PHONY: start_server
start_server:
	$(DOCKER_COMPOSE) up -d
	$(GOOSE) -dir migrations postgres "host=localhost port=5432 user=postgres dbname=supplementapp password=postgres sslmode=disable" up

.PHONY: stop_server
stop_server:
	$(DOCKER_COMPOSE) down --remove-orphans