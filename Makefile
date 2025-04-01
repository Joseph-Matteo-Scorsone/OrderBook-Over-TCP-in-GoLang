build:
	@go build -o bin/tpcexchange cmd/main.go


run: build
	@./bin/tpcexchange