include .env
export

DB_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5432/$(POSTGRES_DB)?sslmode=disable

up:
	docker-compose up
down:
	docker-compose down


migrate-up:
	migrate -path migrations/ -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations/ -database "$(DB_URL)" down

migrate-version:
	migrate -path migrations/ -database "$(DB_URL)" version