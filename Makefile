HEROKU_APP_NAME="sheltered-taiga-32349"

help: ## List targets & descriptions
	@cat Makefile* | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build-linux-binary: ## Build the Go binary for Linux
	GO_ENABLED=0 GOOS=linux go build .

deploy: generate-apidocs build-linux-binary ## Deploy to Heroku. Requires to be logged in on Heroku Registry.
	docker build --rm -f Dockerfile -t registry.heroku.com/$(HEROKU_APP_NAME)/web .
	docker push registry.heroku.com/$(HEROKU_APP_NAME)/web

run-docker: build-linux-binary ## Build Docker image and run it interactively
	docker build --rm -f Dockerfile -t schedulecreator-backend:latest .
	docker run --rm -it -p 8080:8080 schedulecreator-backend:latest

run: ## Build and run locally
	go build .
	./schedulecreator-backend

generate-apidocs: ## Generates API docs from docs/api.yaml. Requires Spectacle.
	spectacle apidocs/api.yaml --target-dir static

clean: ## Clean up
	rm -f schedulecreator-backend
	rm -rf static

.PHONY: help build-linux-binary deploy run-docker run generate-apidocs clean
