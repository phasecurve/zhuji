PHONY: tidy, fmt, test

tidy:
	go mod tidy

fmt:
	go fmt ./...

test:
	go fmt ./...
	gotestsum --debug -f testname -- -count=1 ./...
