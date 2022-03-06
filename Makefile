# Set the shell to bash always
SHELL := /bin/bash

build:
	docker build . -f deploy/Dockerfile -t luebken/awesome-crossplane-providers:latest

run:
	@docker run --env GITHUB_TOKEN=${GITHUB_TOKEN} luebken/awesome-crossplane-providers

.PHONY: build