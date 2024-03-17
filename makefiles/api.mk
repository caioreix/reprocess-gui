IMAGE_NAME=rgui/api
CPATH=cmd/api/.local.env

local:
	@go run ./cmd/api/... -cpath $(CPATH)

watch-local:
	@reflex -r "\.go|.env$$" -s -- sh -c "go run ./cmd/api/... -cpath $(CPATH)"

compose-start: docker-build
	@cd ./cmd/api && $(MAKE) start

compose-stop:
	@cd ./cmd/api && $(MAKE) stop

docker-build:
	@docker build --no-cache --pull -f ./build/api/Dockerfile -t $(IMAGE_NAME) .
