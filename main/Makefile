.PHONY: run build test sqlc-gen migrate-up migrate-down migrate-create docker-up docker-down docker-down-v docker-logs

run:
	go run ./cmd/api/...

build:
	go build -o bin/server ./cmd/api/...

test:
	go test -race ./...

sqlc-gen:
	sqlc generate

migrate-up:
	cd ../db-migrater && DATABASE_URL=$${DATABASE_URL:-postgres://postgres:postgres@localhost:5432/app?sslmode=disable} go run . up

migrate-down:
	cd ../db-migrater && DATABASE_URL=$${DATABASE_URL:-postgres://postgres:postgres@localhost:5432/app?sslmode=disable} go run . down

migrate-create:
	@read -p "Migration name: " name; \
	touch ../db-migrater/migrations/$$(printf "%06d" $$(($$(ls ../db-migrater/migrations/*.up.sql 2>/dev/null | wc -l) + 1)))_$${name}.up.sql; \
	touch ../db-migrater/migrations/$$(printf "%06d" $$(($$(ls ../db-migrater/migrations/*.up.sql 2>/dev/null | wc -l))))_$${name}.down.sql

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

docker-down-v:
	docker compose down -v

docker-logs:
	docker compose logs -f app
