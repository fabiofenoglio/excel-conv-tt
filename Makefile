build:
	go build -o ./dist/excel-converter
run:
	go run main.go
lint:
	golangci-lint run
test:
	echo "no tests available"
build-release:
	make lint
	goreleaser release --snapshot --rm-dist
push-release:
	make lint
	goreleaser release --rm-dist
