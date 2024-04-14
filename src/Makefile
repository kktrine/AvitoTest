run:
	@go run ./main.go --config=env/.env

build:
	@docker build -t openapi .

docker_run:
	@docker compose up

compose_down:
	@docker compose down


test_post:
	@go test -v ./tests/stress_tests/post_test.go

test_stress_get_user_banner:
	@go test -v ./tests/stress_tests/get_user_banner_test.go

test_e2e_get_user_banner:
	@go test -v ./tests/server_tests/get_user_banner_e2e_test.go

test_e2e_get_many_banners:
	@go test -v ./tests/server_tests/get_many_banners_e2e_test.go


static_check:
	@go vet -vettool=$(which staticcheck -f) ./...
