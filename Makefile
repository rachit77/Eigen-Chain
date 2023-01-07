build:	
	go build -o ./bin/eigen-chain

run: build
	./bin/eigen-chain

test:
	go test ./...
