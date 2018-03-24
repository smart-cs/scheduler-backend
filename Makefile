HEROKU_APP_NAME="schedulecreator"

help: ## List targets & descriptions
	@cat Makefile* | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build-linux-binary: # Build the Go binary for Linux
	GO_ENABLED=0 GOOS=linux go build .

deploy: generate-apidocs build-linux-binary ## Deploy to Heroku. Requires to be logged in on Heroku Registry.
	docker build --rm -f Dockerfile -t registry.heroku.com/$(HEROKU_APP_NAME)/web .
	docker push registry.heroku.com/$(HEROKU_APP_NAME)/web
	make clean

run-docker: generate-apidocs build-linux-binary ## Build Docker image and run it interactively locally
	docker build --rm -f Dockerfile -t scheduler-backend:latest .
	docker run --rm -it -p 8080:8080 scheduler-backend:latest

run: ## Build and run locally on port 8080 by default or $PORT if set
	go run main.go

generate-apidocs: ## Generates API docs from docs/api.yaml. Requires Spectacle.
	spectacle apidocs/api.yaml --target-dir static

deps: ## Download dependencies
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

test: ## Run tests
	go test ./...

test-coverage: ## Run tests with coverage
	go test -cover ./...

clean: ## Clean up
	rm -f scheduler-backend
	rm -rf static

.PHONY: help build-linux-binary deploy run-docker run generate-apidocs deps test test-coverage clean
