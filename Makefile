build:
	@go build -o bin/ecom cmd/main.go

test:
	@go test -coverprofile=coverage.out -v ./...

run: build
	@./bin/ecom

migration:
	@go run cmd/migrate/create/main.go $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

%:
	@:
