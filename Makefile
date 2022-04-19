
build: 
	go build -o ./orbit

example:	
	make build
	./scripts/link_examples.sh

license:
	python3 ./scripts/license.py write

test:
	go test `go list ./... | grep -v examples`

integration:
	make example
	python3 ./scripts/integration.py