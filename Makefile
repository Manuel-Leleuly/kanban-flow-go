.PHONY: run
run:
	docker compose -f docker-compose.yaml -f docker-compose.dev.yaml up --build

.PHONY: build-prod
build-prod:
	docker compose -f docker-compose.yaml -f docker-compose.prod.yaml build

.PHONY: run-prod
run-prod:
	docker compose -f docker-compose.yaml -f docker-compose.prod.yaml up

.PHONY: stop
stop:
	docker compose down

.PHONY: stop-prod
stop-prod:
	docker compose -f docker.compose.yaml -f docker-compose.prod.yaml down

.PHONY: clean
clean:
	docker compose --file docker-compose.yml --file docker-compose.override.yml down --rmi all --volumes --remove-orphans
	docker compose --file docker-compose.yml --file docker-compose.prod.yml down --rmi all --volumes --remove-orphans