FRONT_END_BINARY=frontEnd
BROKER_BINARY=broker
AUTH_BINARY=auth
LOGGER_BINARY=logger

## up: starts all containers in the background without forcing build
up:
	docker-compose down
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_all
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build docker images for all services
build_all: build_broker build_auth build_logger build_mail build_listener

build_broker:
	cd broker && docker build -t microservice-broker .
	@echo "Done!"

build_auth:
	cd auth && docker build -t microservice-auth .
	@echo "Done!"

build_logger:
	cd logger && docker build -t microservice-logger .
	@echo "Done!"

build_mail:
	cd mail && docker build -t microservice-mail .
	@echo "Done!"

build_listener:
	cd listener && docker build -t microservice-listener .
	@echo "Done!"


## build_front: builds the frone end binary
build_front:
	@echo "Building front end binary..."
	cd ../front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"

## start: starts the front end
start: build_front
	@echo "Starting front end"
	cd ../front-end && ./${FRONT_END_BINARY} &

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"