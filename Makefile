build:
	@go build -o bin/routes

run: build
	./bin/routes
