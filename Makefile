.PHONY: generate test run

generate:
	templ generate ./...

test:
	go test ./...

run:
	go run ./cmd/example
