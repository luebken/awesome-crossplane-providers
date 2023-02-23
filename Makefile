# Set the shell to bash always
SHELL := /bin/bash

docker-build:
	docker build . -f deploy/Dockerfile -t luebken/awesome-crossplane-providers:latest

docker-push:
	docker push luebken/awesome-crossplane-providers

docker-build-2:
	docker buildx build --platform linux/arm/v7 . -f deploy/Dockerfile -t luebken/awesome-crossplane-providers:latest

docker-run:
	@docker run -v ${PWD}:/repo --env MY_GITHUB_TOKEN=${MY_GITHUB_TOKEN} luebken/awesome-crossplane-providers

docker-run-2:
	@docker run -v ${PWD}/reports:/reports --platform linux/arm/v7 --env MY_GITHUB_TOKEN=${MY_GITHUB_TOKEN} luebken/awesome-crossplane-providers

# Searching for potential Crossplane provider repos.
# Updates providers.txt
run-local-provider-names:
	go run ./cmd/axpp/main.go provider-names

run-local-provider-stats:
	go run ./cmd/axpp/main.go provider-stats