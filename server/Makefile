GOPATH=$(shell go env GOPATH)

LINTER_VERSION=v1.56.2
OAPI_CODEGEN_VERSION=v2.1.0 # TODO: keep an eye on the new version that uses go1.22

API_SWAGGER_DIR=api/swagger
API_SWAGGER_PACKAGE=ecomap

## default: run clean, tidy, vendor, generate-oapi, lint and build
default: clean tidy vendor generate-oapi lint build

## clean: clean the vendor, dist and api generated directories
clean:
	rm -rf vendor/
	rm -rf dist/
	rm -rf api/swagger/ecomap/

## tidy: add missing and remove unused modules
tidy:
	go mod tidy

## vendor: make vendored copy of dependencies
vendor:
	go mod vendor

## generate-oapi: install oapi-codegen and generate the golang types, client and server based on the swagger specification in api/swagger
generate-oapi:
	command -v oapi-codegen || go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION)
	mkdir -p api/swagger/$(API_SWAGGER_PACKAGE)
	oapi-codegen -package $(API_SWAGGER_PACKAGE) -generate types $(API_SWAGGER_DIR)/$(API_SWAGGER_PACKAGE).yml > $(API_SWAGGER_DIR)/$(API_SWAGGER_PACKAGE)/models.gen.go
	oapi-codegen -package $(API_SWAGGER_PACKAGE) -generate client $(API_SWAGGER_DIR)/$(API_SWAGGER_PACKAGE).yml > $(API_SWAGGER_DIR)/$(API_SWAGGER_PACKAGE)/client.gen.go
	oapi-codegen -package $(API_SWAGGER_PACKAGE) -generate std-http $(API_SWAGGER_DIR)/$(API_SWAGGER_PACKAGE).yml > $(API_SWAGGER_DIR)/$(API_SWAGGER_PACKAGE)/server.gen.go

## lint: install golangci-lint and analyse the source code with the configuration in .golangci.yml
lint:
	command -v golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin $(LINTER_VERSION)
	golangci-lint run --timeout=5m

## build: build servre to the dist directory
build:
	go build -o dist/server ./cmd

## help: print this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
