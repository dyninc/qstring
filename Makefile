.PHONY: all style test

all: test

style:
	go vet

test: style
	go test -cover
