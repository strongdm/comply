# This is the default Clever Golang Makefile.
# It is stored in the dev-handbook repo, github.com/Clever/dev-handbook
# Please do not alter this file directly.
GOLANG_MK_VERSION := 1.0.0

SHELL := /bin/bash
SYSTEM := $(shell uname -a | cut -d" " -f1 | tr '[:upper:]' '[:lower:]')
.PHONY: golang-test-deps golang-ensure-curl-installed

# set timezone to UTC for golang to match circle and deploys
export TZ=UTC

# go build flags for use across all commands which accept them
GO_BUILD_FLAGS := "-mod=vendor"

# if the gopath includes several directories, use only the first
GOPATH=$(shell echo $$GOPATH | cut -d: -f1)

# This block checks and confirms that the proper Go toolchain version is installed.
# It uses ^ matching in the semver sense -- you can be ahead by a minor
# version, but not a major version (patch is ignored).
# arg1: golang version
define golang-version-check
_ := $(if  \
		$(shell  \
			expr >/dev/null  \
				`go version | cut -d" " -f3 | cut -c3- | cut -d. -f2 | sed -E 's/beta[0-9]+//'`  \
				\>= `echo $(1) | cut -d. -f2`  \
				\&  \
				`go version | cut -d" " -f3 | cut -c3- | cut -d. -f1`  \
				= `echo $(1) | cut -d. -f1`  \
			&& echo 1),  \
		@echo "",  \
		$(error must be running Go version ^$(1) - you are running $(shell go version | cut -d" " -f3 | cut -c3-)))
endef

# FGT is a utility that exits with 1 whenever any stderr/stdout output is recieved.
# We pin its version since its a simple tool that does its job as-is;
# so we're defended against it breaking or changing in the future.
FGT := $(GOPATH)/bin/fgt
$(FGT):
	go get github.com/GeertJohan/fgt@262f7b11eec07dc7b147c44641236f3212fee89d

golang-ensure-curl-installed:
	@command -v curl >/dev/null 2>&1 || { echo >&2 "curl not installed. Please install curl."; exit 1; }

# Golint is a tool for linting Golang code for common errors.
# We pin its version because an update could add a new lint check which would make
# previously passing tests start failing without changing our code.
GOLINT := $(GOPATH)/bin/golint
$(GOLINT):
	go get golang.org/x/lint/golint@738671d3881b9731cc63024d5d88cf28db875626

# golang-fmt-deps requires the FGT tool for checking output
golang-fmt-deps: $(FGT)

# golang-fmt checks that all golang files in the pkg are formatted correctly.
# arg1: pkg path
define golang-fmt
@echo "FORMATTING $(1)..."
@PKG_PATH=$$(go list -f '{{.Dir}}' $(1)); $(FGT) gofmt -l=true $${PKG_PATH}/*.go
endef

# golang-lint-deps requires the golint tool for golang linting.
golang-lint-deps: $(GOLINT)

# golang-lint calls golint on all golang files in the pkg.
# arg1: pkg path
define golang-lint
@echo "LINTING $(1)..."
@PKG_PATH=$$(go list -f '{{.Dir}}' $(1)); find $${PKG_PATH}/*.go -type f | grep -v gen_ | xargs $(GOLINT)
endef

# golang-lint-deps-strict requires the golint tool for golang linting.
golang-lint-deps-strict: $(GOLINT) $(FGT)

# golang-lint-strict calls golint on all golang files in the pkg and fails if any lint
# errors are found.
# arg1: pkg path
define golang-lint-strict
@echo "LINTING $(1)..."
@PKG_PATH=$$(go list -f '{{.Dir}}' $(1)); find $${PKG_PATH}/*.go -type f | grep -v gen_ | xargs $(FGT) $(GOLINT)
endef

# golang-test-deps is here for consistency
golang-test-deps:

# golang-test uses the Go toolchain to run all tests in the pkg.
# arg1: pkg path
define golang-test
@echo "TESTING $(1)..."
@go test $(GO_BUILD_FLAGS) -v $(1)
endef

# golang-test-strict-deps is here for consistency
golang-test-strict-deps:

# golang-test-strict uses the Go toolchain to run all tests in the pkg with the race flag
# arg1: pkg path
define golang-test-strict
@echo "TESTING $(1)..."
@go test -v $(GO_BUILD_FLAGS) -race $(1)
endef

# golang-vet-deps is here for consistency
golang-vet-deps:

# golang-vet uses the Go toolchain to vet all the pkg for common mistakes.
# arg1: pkg path
define golang-vet
@echo "VETTING $(1)..."
@go vet $(GO_BUILD_FLAGS) $(1)
endef

# golang-test-all-deps installs all dependencies needed for different test cases.
golang-test-all-deps: golang-fmt-deps golang-lint-deps golang-test-deps golang-vet-deps

# golang-test-all calls fmt, lint, vet and test on the specified pkg.
# arg1: pkg path
define golang-test-all
$(call golang-fmt,$(1))
$(call golang-lint,$(1))
$(call golang-vet,$(1))
$(call golang-test,$(1))
endef

# golang-test-all-strict-deps: installs all dependencies needed for different test cases.
golang-test-all-strict-deps: golang-fmt-deps golang-lint-deps-strict golang-test-strict-deps golang-vet-deps

# golang-test-all-strict calls fmt, lint, vet and test on the specified pkg with strict
# requirements that no errors are thrown while linting.
# arg1: pkg path
define golang-test-all-strict
$(call golang-fmt,$(1))
$(call golang-lint-strict,$(1))
$(call golang-vet,$(1))
$(call golang-test-strict,$(1))
endef

# golang-build: builds a golang binary. ensures CGO build is done during CI. This is needed to make a binary that works with a Docker alpine image.
# arg1: pkg path
# arg2: executable name
define golang-build
@echo "BUILDING..."
@if [ -z "$$CI" ]; then \
	go build $(GO_BUILD_FLAGS) -o bin/$(2) $(1); \
else \
	echo "-> Building CGO binary"; \
	CGO_ENABLED=0 go build $(GO_BUILD_FLAGS) -installsuffix cgo -o bin/$(2) $(1); \
fi;
endef

# golang-update-makefile downloads latest version of golang.mk
golang-update-makefile:
	@wget https://raw.githubusercontent.com/Clever/dev-handbook/master/make/golang-v1.mk -O /tmp/golang.mk 2>/dev/null
	@if ! grep -q $(GOLANG_MK_VERSION) /tmp/golang.mk; then cp /tmp/golang.mk golang.mk && echo "golang.mk updated"; else echo "golang.mk is up-to-date"; fi
