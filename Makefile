## default: TODO docs
default: 
	echo "TODO"

# TODO: Update readme with docker as requirement (add link to docs) and how to set up secrets
# TODO: add migrations test docker
# TODO: use secrets in github

## dev: TODO docs
dev:
	echo "TODO"

## help: print this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
