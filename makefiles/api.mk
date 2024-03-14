IMAGE_NAME=rgui/api

local:
	go run ./cmd/api/...

compose-start: docker-build
	@cd ./cmd/api && $(MAKE) start

compose-stop:
	@cd ./cmd/api && $(MAKE) stop

docker-build:
	@docker build --no-cache --pull -f ./build/api/Dockerfile -t $(IMAGE_NAME) .
