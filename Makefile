defaut:
	go build -o bin/gentes
fmt:
	gofmt -w .
edit:
	vim *.go bin/data/regles.la bin/data/lexsynt.la bin/corpus/test.txt
darwin:
	env OSXCROSS_NO_INCLUDE_PATH_WARNINGS=1 MACOSX_DEPLOYMENT_TARGET=10.6 CC=o64-clang CXX=o64-clang++ GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o darwin/gentes
w:
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc env GOOS=windows GOARCH=amd64 go build -o win/gentes.exe
