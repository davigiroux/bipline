.PHONY: schema generate test probe

schema:
	bash tools/fetch-schema.sh

generate:
	go generate ./internal/buffer/

test:
	go test ./...

probe:
	go run ./cmd/probe/
