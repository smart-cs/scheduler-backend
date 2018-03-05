# scheduler-backend
Retrieves and uses data from UBC Course Schedule.

## Developing
Make sure you have `Go`, `make`, `Docker` installed

Get the code (to your `$GOPATH`):
```
go get -d github.com/smart-cs/scheduler-backend
```

Run the server locally:
```shell
make run
```

## Make Commands
```shell
$ make help
clean                          Clean up
deploy                         Deploy to Heroku. Requires to be logged in on Heroku Registry.
deps                           Download dependencies
generate-apidocs               Generates API docs from docs/api.yaml. Requires Spectacle.
help                           List targets & descriptions
run-docker                     Build Docker image and run it interactively locally
run                            Build and run locally on port 8080 by default or $PORT if set
test-coverage                  Run tests with coverage
test                           Run tests
```
