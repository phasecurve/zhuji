PHONY: build, tidy, fmt, test

build:
	go build ./...

tidy:
	go mod tidy

fmt:
	go fmt ./...

test:
	go fmt ./...
	gotestsum --debug -f testname -- -timeout 8s -count=1 ./...
