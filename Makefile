dev:
	air
sandbox:
	docker compose -f ./infrastructure/docker-compose.yml up
teardown:
	docker compose -f ./infrastructure/docker-compose.yml down
