
build: 
	go build -o ./orbit

example:	
	make build
	./scripts/link_examples.sh

license:
	python3 ./scripts/license.py write

gotest:
	go test `go list ./... | grep -v examples`

test:
	make gotest

	echo 'running integration tests'
	make example
	pytest ./scripts/test

