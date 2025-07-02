install:
	go install github.com/air-verse/air@v1.61.7
	@echo "\033[0;32mAir installed successfully. Use 'make dev' to start the development server.\033[0m"
	go mod tidy
	@echo "\033[0;32mGo modules installed successfully!\033[0m"
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
test:
	go test ./tests/...