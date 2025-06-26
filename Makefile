dev:
	air
sandbox-linux:
	docker compose -f ./infra/linux/docker-compose.yml up
sandbox-windows:
	docker compose -f ./infra/windows/docker-compose.yml up
teardown:
	docker compose -f ./infra/docker-compose.yml down
