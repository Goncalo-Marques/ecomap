BUILD_GIT_HASH = $(shell git describe --always --dirty)
BUILD_TIMESTAMP = $(shell TZ="GMT" LC_TIME="en_US.utf8" date)

## default: build and start docker containers 
default: up

# TODO: Update readme with docker as requirement (add link to docs) and how to set up secrets
# TODO: add migrations test docker
# TODO: use secrets in github

## dev: TODO docs
dev:
	echo "TODO"

## build: build server, web app and docker containers
build:
	make -C server
	npm --prefix web ci
	npm run --prefix web build
	mv web/dist server/dist/web
	docker compose build --build-arg BUILD_GIT_HASH="$(BUILD_GIT_HASH)" --build-arg BUILD_TIMESTAMP="$(BUILD_TIMESTAMP)"

## up: build and start docker containers
up: build
	docker compose up --no-build

## help: print this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
