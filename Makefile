defaut:
	go build -o bin/gentes
fmt:
	go fmt ./
edit:
	vim -c ":b main.go" *.go bin/data/groupes.la bin/data/lexsynt.la
darwin:
	env GOOS=darwin GOARCH=amd64 go build -o mac/gentes
w:
	env GOOS=windows GOARCH=amd64 go build -o win/gentes.exe
