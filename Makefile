.PHONY: lint
lint:
	golangci-lint run --timeout="5m" ./...

.PHONY: test
test:
	go test -timeout "10m" -race -cover -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: build
build:
	go build ./cmd/changelog
