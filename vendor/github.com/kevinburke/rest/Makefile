SHELL = /bin/bash -o pipefail

BUMP_VERSION := $(GOPATH)/bin/bump_version
STATICCHECK := $(GOPATH)/bin/staticcheck

deps:
	go get -u ./...

test: lint
	go test ./...

lint: | $(STATICCHECK)
	go vet ./...
	$(STATICCHECK) ./...

$(STATICCHECK):
	go get honnef.co/go/tools/cmd/staticcheck

race-test: lint
	go test -race ./...

$(BUMP_VERSION):
	go get github.com/kevinburke/bump_version

release: race-test | $(BUMP_VERSION)
	$(BUMP_VERSION) minor client.go
