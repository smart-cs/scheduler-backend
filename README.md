# scheduler-backend

[![CircleCI](https://circleci.com/gh/Smart-CS/scheduler-backend.svg?style=shield)](https://circleci.com/gh/Smart-CS/scheduler-backend)
[![Go Report Card](https://goreportcard.com/badge/github.com/smart-cs/scheduler-backend)](https://goreportcard.com/report/github.com/smart-cs/scheduler-backend)
[![codecov](https://codecov.io/gh/Smart-CS/scheduler-backend/branch/master/graph/badge.svg)](https://codecov.io/gh/Smart-CS/scheduler-backend)

Retrieves and uses data from UBC Course Schedule.

## Developing

Make sure you have `Go`, `make`, `Docker` installed

To get the code (to your `$GOPATH`):

```
go get -d github.com/smart-cs/scheduler-backend
```

By default, your `$GOPATH` is `~/go`. If you haven't changed it, the repository should end up in `~/go/src/github.com/smart-cs/scheduler-backend`

To run the server locally:

```shell
make run
```

## Make Commands

```shell
$ make help
clean                          Clean up
deploy                         Deploy to Heroku. Requires to be logged in on Heroku Registry.
deps                           Download dependencies
generate-apidocs               Generates API docs from docs/api.yml. Requires Spectacle.
help                           List targets & descriptions
run-docker                     Build Docker image and run it interactively locally
run                            Build and run locally on port 8080 by default or $PORT if set
test-coverage                  Run tests with coverage
test                           Run tests
```
