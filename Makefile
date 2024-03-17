.PHONY: make

setup:
	@go install github.com/cespare/reflex@latest
	@# @curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.43.0

deps:
	@go mod tidy
	@go mod download

lint:
	@golangci-lint run ./...

test:
	@go test -v -race ./...

local-api:
	@$(MAKE) -f makefiles/api.mk local

watch-api:
	@$(MAKE) -f makefiles/api.mk watch-local

compose-api:
	@$(MAKE) -f makefiles/api.mk compose-start

compose-api-stop:
	@$(MAKE) -f makefiles/api.mk compose-stop
