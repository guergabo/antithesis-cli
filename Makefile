.PHONY: lint fmt antithesis test test-cover clean

BUILD = build/$(GOOS)/$(GOARCH)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOEXE ?= $(shell go env GOEXE)
GO ?= go 
ANTITHESIS = ${BUILD}/antithesis${GOEXE}

lint: 
	golangci-lint run ./... 

fmt: 
	${GO} fmt ./... 

antithesis: 
	${GO} build -o ${ANTITHESIS} . 

test: antithesis 
	${GO} test ./... 

test-cover: antithesis 
	${GO} test -cover ./...

clean: 
	rm -rf ./build
