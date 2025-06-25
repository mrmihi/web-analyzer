dev:
	air
sandbox:
	docker compose -f ./infra/docker-compose.yml up
teardown:
	docker compose -f ./infra/docker-compose.yml down
