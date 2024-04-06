.PHONY: make

setup:
	@go install github.com/cespare/reflex@latest
	@go install github.com/vektra/mockery/v2@latest
	@go install github.com/mgechev/revive@latest
	@go install github.com/securego/gosec/v2/cmd/gosec@latest

deps:
	@go mod tidy
	@go mod download

lint:
	@revive ./...

mocks:
	@find ./internal/apps/*/core/port -type f -name '*.go' -exec bash -c 'dir=$$(dirname "{}"); cd $$dir; mockery --dir . --outpkg $$(basename "$$dir")mock --name=".*"' \;

mocks-quiet:
	@$(MAKE) mocks > /dev/null 2>&1

test:
	@go test -race -count=1 -cover ./...

sec:
	@gosec -quiet -exclude-generated ./...

pre-commit:
	@echo Updating go.mod...
	@go mod tidy
	@echo Building mocks...
	@$(MAKE) mocks-quiet -s
	@echo Running lint...
	@$(MAKE) lint -s
	@echo Running gosec...
	@$(MAKE) sec -s
	@echo Running tests...
	@$(MAKE) test -s

local-api:
	@$(MAKE) -f makefiles/api.mk local

watch-api:
	@$(MAKE) -f makefiles/api.mk watch-local

compose-api:
	@$(MAKE) -f makefiles/api.mk compose-start

compose-api-stop:
	@$(MAKE) -f makefiles/api.mk compose-stop
