.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./db/migrations ${name}
.PHONY: db/migrations/up
db/migrations/up: 
	@echo 'Running up migrations...'
	migrate -path ./db/migrations -database postgresql://user:password@localhost:5432/quest?sslmode=disable up

run:
	go run ./cmd/main.go
