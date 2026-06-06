-include .env
export

.PHONY: schema generate test probe test-integration

schema:
	bash tools/fetch-schema.sh

generate:
	go generate ./internal/buffer/

test:
	go test ./...

probe:
	go run ./cmd/probe/

test-integration:
	go test -tags integration -v ./internal/generator/
