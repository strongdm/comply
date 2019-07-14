.PHONY: all clean checks test build image dependencies

SRCS = $(shell git ls-files '*.go' | grep -v '^vendor/')

MAIN_DIRECTORY := ./go-bindata/
BIN_OUTPUT := dist/go-bindata

TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse HEAD)
VERSION := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))

default: clean test build

clean:
	rm -rf dist/ builds/ cover.out

build: clean
	@echo Version: $(VERSION)
	go build -v -ldflags '-X "main.AppVersion=${VERSION}"' -o ${BIN_OUTPUT} ${MAIN_DIRECTORY}

dependencies:
	dep ensure -v

test: clean
	go test -v -cover ./...

testdata:
	make -C testdata

fmt:
	gofmt -s -l -w $(SRCS)
