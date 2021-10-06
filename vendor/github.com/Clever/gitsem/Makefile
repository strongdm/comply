include golang.mk
.DEFAULT_GOAL := test # override default goal set in library makefile

SHELL := /bin/bash
PKG := github.com/Clever/gitsem
PKGS := $(shell go list ./... | grep -v /vendor)
EXECUTABLE := gitsem
VERSION := $(shell cat VERSION)
BUILDS := \
	build/$(EXECUTABLE)-v$(VERSION)-darwin-amd64 \
	build/$(EXECUTABLE)-v$(VERSION)-linux-amd64 \
	build/$(EXECUTABLE)-v$(VERSION)-windows-amd64
COMPRESSED_BUILDS := $(BUILDS:%=%.tar.gz)
RELEASE_ARTIFACTS := $(COMPRESSED_BUILDS:build/%=release/%)

.PHONY: test golint build vendor

$(eval $(call golang-version-check,1.13))

test: $(PKGS)

$(PKGS): golang-test-all-strict-deps
	$(call golang-test-all-strict,$@)


run:
	@go run main.go

build:
	go build -o bin/$(EXECUTABLE) $(PKG)

build/$(EXECUTABLE)-v$(VERSION)-darwin-amd64:
	GOARCH=amd64 GOOS=darwin go build -o "$@/$(EXECUTABLE)"
build/$(EXECUTABLE)-v$(VERSION)-linux-amd64:
	GOARCH=amd64 GOOS=linux go build -o "$@/$(EXECUTABLE)"
build/$(EXECUTABLE)-v$(VERSION)-windows-amd64:
	GOARCH=amd64 GOOS=windows go build -o "$@/$(EXECUTABLE).exe"
build: $(BUILDS)
%.tar.gz: %
	tar -C `dirname $<` -zcvf "$<.tar.gz" `basename $<`
$(RELEASE_ARTIFACTS): release/% : build/%
	mkdir -p release
	cp $< $@
release: $(RELEASE_ARTIFACTS)

clean:
	rm -rf build release


install_deps:
	go mod vendor
