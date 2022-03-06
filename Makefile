# Set the shell to bash always
SHELL := /bin/bash

build:
	docker buildx build --platform linux/arm/v7 . -f deploy/Dockerfile -t luebken/awesome-crossplane-providers:latest
	docker push luebken/awesome-crossplane-providers

run:
	@docker run -v ${PWD}:/data --env MY_GITHUB_TOKEN=${MY_GITHUB_TOKEN} luebken/awesome-crossplane-providers

build-2:
	docker build . -f deploy/Dockerfile -t luebken/awesome-crossplane-providers:latest

run-2:
	@docker run -v ${PWD}:/data --platform linux/arm/v7 --env MY_GITHUB_TOKEN=${MY_GITHUB_TOKEN} luebken/awesome-crossplane-providers


run-local:
	go run ./cmd/axpp/main.go