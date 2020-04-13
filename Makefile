# Run the main executable sourcing .env
run-source-env:
	go build -o bin/main main.go && source .env && bin/main

# Run the main executable with current env variables
run:
	go build -o bin/main main.go && bin/main
