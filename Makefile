run:
	go run cmd/app/main.go

migrate-up:
	migrate -path internal/database/migrations -database "sqlite://data/gustav.db" up

migrate-down:
	migrate -path internal/database/migrations -database "sqlite://data/gustav.db" down