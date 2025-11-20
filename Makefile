PHONY: build, tidy, fmt, test, asm

build:
	go build ./...

tidy:
	go mod tidy

fmt:
	go fmt ./...

test:
	go fmt ./...
	gotestsum --debug -f testname -- -timeout 8s -count=1 ./...

asm:
	as -o output.o output.s
	ld -o output output.o
