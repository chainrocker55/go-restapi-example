test:
	go mod tidy
	go test -v ./...
run:
	go run main.go