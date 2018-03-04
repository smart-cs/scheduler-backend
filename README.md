# schedulecreator-backend
Retrieves and uses data from UBC Course Schedule.

## Developing
Make sure you have `Go`, `make`, `Docker` installed

Get the code (to your `$GOPATH`):
```
go get -d github.com/nickwu241/schedulecreator-backend
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
run                            Build and run locally
test-coverage                  Run tests with coverage
test                           Run tests
```
