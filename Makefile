.PHONY: run infra-up infra-down sqlc test fmt

run:
	go run ./apps/api

infra-up:
	docker compose up -d

infra-down:
	docker compose down

sqlc:
	sqlc generate

test:
	go test ./...

fmt:
	gofmt -w apps internal
