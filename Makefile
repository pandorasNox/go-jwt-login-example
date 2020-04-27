

include ./hack/help.mk


UID:=$(shell id -u)
GID:=$(shell id -g)
PWD:=$(shell pwd)


.PHONY: setup
setup: ##@setup builds the container image(s) and starts the setup
	docker-compose --compatibility build
	docker-compose --compatibility up -d


.PHONY: build
build: ##@setup builds the container image(s)
	docker-compose --compatibility build --build-arg CMD_NAME=go-jwt

