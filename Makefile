.PHONY: all coverage style test

all: test

coverage:
	go test -v -cover -coverprofile=coverage.out

style:
	go vet

test: style
	go test -v -cover
