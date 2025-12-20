.PHONY: build run test clean

build:
	go build -o myshell app/*.go

run: build
	./myshell

test:
	go test ./app/...

clean:
	rm -f myshell /tmp/codecrafters-build-shell-go

install: build
	cp myshell /usr/local/bin/