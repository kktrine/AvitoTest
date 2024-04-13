
.PHONY: run_service
run:
	@go run ./main.go --config=env/.env

.PHONY: compose
build:
	@docker build -t openapi .

docker_run:
	@docker compose up

compose_down:
	@docker compose down


.PHONY: staticcheck
staticcheck:
	@go vet -vettool=$(which staticcheck -f) ./...
