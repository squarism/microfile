build:
	gb build

test:
	gb test -v

dist:
	@if [ -d release ]; then rm -rf release; fi
	@mkdir release
	@gb build
	@mv bin/microfile release/microfile-osx-amd64
	@GOOS=linux gb build
	@mv bin/microfile-linux-amd64 release
	@GOOS=windows gb build
	@mv bin/microfile-windows-amd64.exe release
	@echo "\nTo push a binary release and tag to github:"
	@echo "ghr -u squarism [tag] release"

clean:
	rm bin/* && rmdir bin
	rm pkg/* && rmdir pkg
	rm -rf vendor/pkg
