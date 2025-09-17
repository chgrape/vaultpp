build:
	go build ./cmd/api/main.go

run: build
	./main

