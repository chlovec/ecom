build:
	@go build -o bin/ecom cmd/main.go

test:
	@go test -coverprofile=coverage.out -v ./...

run: build
	@./bin/ecom