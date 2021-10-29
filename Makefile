run:
	go run cmd/trove/main.go

build:
	go build -o bin/ ./cmd/trove

dev:
	gow run cmd/trove/main.go

