## default: TODO
default: 

# TODO: Update readme with docker as requirement (add link to docs) and how to set up secrets
# TODO: Serve docs: https://swagger.io/docs/open-source-tools/swagger-ui/usage/installation/?sbsearch=docker

## dev: TODO
dev:
	docker compose up

## help: print this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
