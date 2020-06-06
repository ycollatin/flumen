defaut:
	go build -o bin/gentes
fmt:
	go fmt ./
ed:
	nvim -c ":b branche.go" *.go bin/data/groupes.la bin/data/lexsynt.la bin/corpus/test.txt
edit:
	vim -c ":b branche.go" *.go bin/data/groupes.la bin/data/lexsynt.la bin/corpus/test.txt
darwin:
	env GOOS=darwin GOARCH=amd64 go build -o mac/gentes
w:
	env GOOS=windows GOARCH=amd64 go build -o win/gentes.exe
dlv:
	dlv --headless debug
