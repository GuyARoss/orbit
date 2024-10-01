# Contributing to orbit

## Finding a task

Before contributing to this project, please ensure that an open issue first exists outlining that issue you plan to resolve.

## Running the project

To run the project, please first download the required tools.

- nodeJS >= v11
- golang 1.14

## Building the project

After golang is installed on your machine, you can build with `go run build`

## Running tests

### Requirements

- lighthouse CLI `npm i -g lighthouse`
- pytest `pip install pytest`

You can run the entire test suite with `make test` or the go tests with `make gotest` or integration tests with `make integrationtest`

## Committing code

This project uses the golang project standard for commit messages, you find it [here](https://go.dev/doc/contribute#commit_messages).
