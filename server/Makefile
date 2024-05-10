GOPATH=$(shell go env GOPATH)

BUILD_GIT_HASH = $(shell git describe --always)
BUILD_TIMESTAMP = $(shell TZ="GMT" LC_TIME="en_US.utf8" date)

LINTER_VERSION=v1.56.2
OAPI_CODEGEN_VERSION=80f0b97 # TODO: keep an eye on the new version that uses go1.22

API_SWAGGER_SPEC_FILE=api/swagger/ecomap.yml
API_SWAGGER_GEN_DIR=api/ecomap
API_SWAGGER_GEN_PACKAGE=ecomap

## default: run clean, generate, tidy, vendor, lint, test and build
default: clean generate tidy vendor lint test build

## dev: run clean, generate, tidy, vendor and build
dev: clean generate tidy vendor build

## clean: clean the vendor, dist and api generated directories
clean:
	rm -rf vendor/
	rm -rf dist/
	rm -rf $(API_SWAGGER_GEN_DIR)/

## tidy: add missing and remove unused modules
tidy:
	go mod tidy

## vendor: make vendored copy of dependencies
vendor:
	go mod vendor

## generate-oapi: install oapi-codegen and generate the golang types, client and server based on the swagger specification in api/swagger
generate-oapi:
	go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION)
	mkdir -p $(API_SWAGGER_GEN_DIR)
	oapi-codegen -package $(API_SWAGGER_GEN_PACKAGE) -generate types $(API_SWAGGER_SPEC_FILE) > $(API_SWAGGER_GEN_DIR)/models.gen.go
	oapi-codegen -package $(API_SWAGGER_GEN_PACKAGE) -generate client $(API_SWAGGER_SPEC_FILE) > $(API_SWAGGER_GEN_DIR)/client.gen.go
	oapi-codegen -package $(API_SWAGGER_GEN_PACKAGE) -generate std-http $(API_SWAGGER_SPEC_FILE) > $(API_SWAGGER_GEN_DIR)/server.gen.go

## generate: run generate-oapi
generate: generate-oapi

## lint: install golangci-lint and analyze the source code with the configuration in .golangci.yml
lint:
	command -v golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin $(LINTER_VERSION)
	golangci-lint run --timeout=5m

## test-unit: run unit tests
test-unit:
	go test ./...

## test-integration: run integration tests
test-integration:
	go test -tags=integration ./...

## test: run unit and integration tests
test: test-unit test-integration

BUILD_FLAGS = " \
-X 'main.BuildGitHash=$(BUILD_GIT_HASH)' \
-X 'main.BuildTimestamp=$(BUILD_TIMESTAMP)'"

## build: build server to the dist directory
build:
	go build -ldflags $(BUILD_FLAGS) -o dist/server ./cmd/server

## help: print this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
