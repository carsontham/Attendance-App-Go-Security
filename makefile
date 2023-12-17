build:
	@go build -o bin/attendance-app ./cmd/attendance-app

run: build
	@./bin/attendance-app

# test:
# 	@go test -v ./...