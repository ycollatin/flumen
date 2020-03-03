
defaut:
	go build -o bin/gentes
edit:
	vim -c ":b main.go" *.go
darwin:
	env GOOS=darwin GOARCH=amd64 go build -o mac/publicola
w:
	env GOOS=windows GOARCH=amd64 go build -o win/publicola.exe
