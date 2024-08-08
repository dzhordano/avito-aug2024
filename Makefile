.PHONY: all

build:
	go mod download && go build -o ./.bin/app ./cmd/app/main.go

run: build
	docker compose up --remove-orphans app

lint:
	golangci-lint run

run-no-docker: build
	make migrate

	go run ./cmd/app/main.go

init-db:
	docker run --name=bootcampdb -e POSTGRES_PASSWORD=qwerty -p 5555:5432 -d postgres

migrate:
	migrate -database "postgresql://postgres:qwerty@localhost:5555/postgres?sslmode=disable" -source "file://migrations" up

init-db-test:
	docker run --name=bootcampdb_test -e POSTGRES_PASSWORD=qwerty -p 5556:5432 -d postgres

export TEST_CONTAINER_NAME=bootcampdb_test
test.integration:
	GIN_MODE=release go test -v ./tests/
	docker stop $$TEST_CONTAINER_NAME

test:
	go test --short -coverprofile=coverage.out -v ./...
	make test.coverage

test.coverage:
	go tool cover -func=cover.out | grep "total"

swag:
	swag init -o ./docs -g ./cmd/app/main.go