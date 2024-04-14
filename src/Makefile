run:
	@go run ./main.go --config=env/.env

build:
	@docker build -t openapi .

docker_build:
	@docker build -t openapi .

docker_up:
	@docker compose up

docker_down:
	@docker compose down

server_tests:
	@go test -v ./tests/server_tests

stress_tests:
	@make test_stress_post
	@make test_stress_get_user_banner

test_stress_post:
	@go test -v ./tests/stress_tests/post_test.go

test_stress_get_user_banner:
	@go test -v ./tests/stress_tests/get_user_banner_test.go

test_e2e_get_user_banner:
	@go test -v ./tests/server_tests/get_user_banner_e2e_test.go

test_e2e_get_many_banners:
	@go test -v ./tests/server_tests/get_many_banners_e2e_test.go

test_e2e_post:
	@go test -v ./tests/server_tests/post_e2e_test.go

test_e2e_delete:
	@go test -v ./tests/server_tests/delete_e2e_test.go

test_e2e_patch:
	@go test -v ./tests/server_tests/patch_e2e_test.go

check:
	@go vet -vettool=$(which staticcheck -f) ./...
