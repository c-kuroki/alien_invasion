all: alien_invasion

alien_invasion:
	go build -o ./cmd/alien_invasion/alien_invasion ./cmd/alien_invasion/...

.PHONY: lint
lint:
	@echo "Linting code..."
	gofmt -s -w ./.
	golangci-lint run ./... -E gofmt
	go mod tidy
	@echo "Linting complete!"

.PHONY: test
test: clean
	@echo "Running tests..."
	go test ./... -v
	@echo "Tests complete!"

.PHONY : clean
clean:
	@echo "Cleaning env..."
	go clean -cache
	go clean -testcache
	rm -f ./cmd/alien_invasion/alien_invasion
	@echo "Cleaned env!"

