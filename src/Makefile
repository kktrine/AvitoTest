.PHONY: run_service
run_service:
	@go run ./main.go --config=env/.env

.PHONY: compose
compose:
	@docker-compose up --build banner-server

.PHONY: staticcheck
staticcheck:
	@go vet -vettool=$(which staticcheck -f) ./...
