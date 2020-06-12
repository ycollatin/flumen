defaut:
	go build -o bin/gentes
fmt:
	go fmt ./
edit:
	vim -c ":b branche.go" *.go bin/data/groupes.la bin/data/lexsynt.la bin/corpus/test.txt
darwin:
	#env CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o mac/gentes
	#env OSXCROSS_NO_INCLUDE_PATH_WARNINGS=1 MACOSX_DEPLOYMENT_TARGET=10.6 CC=o64-clang CXX=o64-clang++ GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build
	env OSXCROSS_NO_INCLUDE_PATH_WARNINGS=1 MACOSX_DEPLOYMENT_TARGET=10.6 CC=o64-clang CXX=o64-clang++ GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build
w:
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc env GOOS=windows GOARCH=amd64 go build -o win/gentes.exe
dlv:
	dlv --headless debug
