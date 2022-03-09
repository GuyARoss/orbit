
build: 
	go build -o ./orbit

example:	
	make build
	./scripts/link_examples.sh

license:
	python3 ./scripts/license.py