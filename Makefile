.PHONY: fmt test

fmt:
	gci write .
	gofumpt -l -w .

test:
	go test -v ./...
