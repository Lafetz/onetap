.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./db/migrations ${name}
.PHONY: db/migrations/up
db/migrations/up: 
	@echo 'Running up migrations...'
	migrate -path ./db/migrations -database ${url} up
run:
	@echo "Loading environment variables from .env file"
	@set -o allexport; source ./load_env.sh; set +o allexport; \
	echo "Running Go application"; \
	go run ./cmd/main.go
lint:
	golangci-lint run
test:
	go test  ./... 
coverage:
	go test  -coverprofile=coverage.out ./... ;
	go tool cover -func=coverage.out
docs:
	
build:
	go build -o ./bin/web ./cmd/web