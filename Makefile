include .env
export

DB_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5432/$(POSTGRES_DB)?sslmode=disable

up:
	docker-compose up --build -d
down:
	docker-compose down

clear:
	docker-compose down -v

migrate-up:
	migrate -path services/core/migrations/ -database "$(DB_URL)" up

migrate-down:
	migrate -path services/core/migrations/ -database "$(DB_URL)" down

migrate-version:
	migrate -path services/core/migrations/ -database "$(DB_URL)" version