build:
	gb build

test:
	gb test -v

dist:
	GOOS=linux gb build
	mv bin/microfile-linux-amd64 bin/microfile.linux

