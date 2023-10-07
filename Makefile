build:
	@go build -o bin/drive-editor

run: build
	@./bin/drive-editor

test:
	@go test -v ./...