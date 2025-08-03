include .env

migrate:
	goose -dir db/migrations postgres "$(DB_CONN)" up

downgrade:
	goose -dir db/migrations postgres "$(DB_CONN)" down

migrate-status:
	goose -dir db/migrations postgres "$(DB_CONN)" status

migrate-create:
	goose -dir db/migrations create $(name) sql

sqlc:
	docker run --rm \
	-v "$(CURDIR):/src" \
	-w /src/db \
	sqlc/sqlc generate -f sqlc.yaml
	git add ./internal