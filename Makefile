build:
	go build -o ./dist/excel-converter
run:
	go run main.go
lint:
	golangci-lint run
test:
	go test ./...
clean:
	go mod tidy
	go fmt $(go list ./... | grep -v /vendor/)
check:
	make build
	make clean
	make lint
	make test
build-release:
	make check
	goreleaser release --snapshot --rm-dist
push-release:
	make check
	goreleaser release --rm-dist
