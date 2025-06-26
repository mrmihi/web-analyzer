build:
	go build -o ./bin/server ./src
start:
	./bin/server
dev:
	air
sandbox-linux:
	docker compose -f ./infra/linux/docker-compose.yml up
sandbox-windows:
	docker compose -f ./infra/windows/docker-compose.yml up
teardown-linux:
	docker compose -f ./infra/linux/docker-compose.yml down
teardown-windows:
	docker compose -f ./infra/windows/docker-compose.yml down
