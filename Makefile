# Run tests
test:
	go test -v -covermode=count -coverprofile=coverage.out -tags test ./... 

# Run the main executable sourcing .env
run-source-env:
	go build -o bin/main main.go && source .env && bin/main

# Run the main executable with current env variables
run:
	go build -o bin/main main.go && bin/main

# Installs dependencies
install:
	go get ./...

