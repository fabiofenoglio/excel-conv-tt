build:
	go build -o ./releases/converter .
	GOOS=windows GOARCH=amd64 go build -o ./releases/converter-win-amd64.exe .