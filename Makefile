build:
	gb build

test:
	gb test -v

dist:
	GOOS=linux gb build
	mv bin/dropboy-linux-amd64 bin/dropboy.linux

