build:
	gb build

test:
	gb test -v

dist:
	GOOS=linux gb build
	mv bin/microfile-linux-amd64 bin/microfile.linux

clean:
	rm bin/* && rmdir bin
	rm pkg/* && rmdir pkg
	rm -rf vendor/pkg
