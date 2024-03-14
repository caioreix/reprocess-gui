.PHONY: build

LINUX_AMD64 = GOOS=linux GOARCH=amd64 GO111MODULE=on

build:
	@$(LINUX_AMD64) go build -a -v -o api-build main.go

start:
	@docker compose --env-file .local.env -f docker-compose.yml up -d

stop:
	@docker compose --env-file .local.env -f docker-compose.yml down
