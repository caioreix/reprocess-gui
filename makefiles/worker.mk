IMAGE_NAME=rgui/worker
CPATH=cmd/worker/.local.env

local:
	@go run ./cmd/worker/... -cpath $(CPATH)

watch-local:
	@reflex -r "\.go|.env$$" -s -- sh -c "go run ./cmd/worker/... -cpath $(CPATH)"

compose-start: docker-build
	@cd ./cmd/worker && $(MAKE) start

compose-stop:
	@cd ./cmd/worker && $(MAKE) stop

docker-build:
	@docker build --no-cache --pull -f ./build/worker/Dockerfile -t $(IMAGE_NAME) .
