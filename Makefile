.DEFAULT_GOAL := comply
GO_SOURCES := $(shell find . -name '*.go')
THEME_SOURCES := $(shell find themes)

assets: $(THEME_SOURCES)
	go-bindata-assetfs -pkg theme -prefix themes themes/...
	mv bindata_assetfs.go internal/theme/themes_bindata.go

comply: assets $(GO_SOURCES)
	go build github.com/strongdm/comply/cmd/comply

dist: clean
	$(eval VERSION := $(shell git describe --tags --always --dirty="-dev"))
	$(eval LDFLAGS := -ldflags='-X "cli.Version=$(VERSION)"')
	mkdir dist
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -gcflags=-trimpath=$(GOPATH) -asmflags=-trimpath=$(GOPATH) $(LDFLAGS) -o dist/comply-$(VERSION)-darwin-amd64 github.com/strongdm/comply/cmd/comply
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -gcflags=-trimpath=$(GOPATH) -asmflags=-trimpath=$(GOPATH) $(LDFLAGS) -o dist/comply-$(VERSION)-linux-amd64 github.com/strongdm/comply/cmd/comply
	cd dist && tar -czvf comply-$(VERSION)-darwin-amd64.tgz comply-$(VERSION)-darwin-amd64
	cd dist && tar -czvf comply-$(VERSION)-linux-amd64.tgz comply-$(VERSION)-linux-amd64
clean:
	rm -rf dist
	rm -f comply

install: assets $(GO_SOURCES)
	go install github.com/strongdm/comply/cmd/comply

export-example:
	cp example/narratives/* themes/comply-soc2/narratives
	cp example/procedures/* themes/comply-soc2/procedures
	cp example/policies/* themes/comply-soc2/policies
	cp example/standards/* themes/comply-soc2/standards
	cp example/templates/* themes/comply-soc2/templates

docker:
	cd build && docker build -t strongdm/pandoc .
	docker tag jagregory/pandoc:latest strongdm/pandoc:latest
	docker push strongdm/pandoc

cleanse:
	git checkout --orphan newbranch
	git add -A
	git commit -m "Initial commit"
	git branch -D master
	git branch -m master
	git push -f origin master
	git gc --aggressive --prune=all

release: dist release-deps
	$(eval VERSION := $(shell git describe --tags --always --dirty="-dev"))
	github-release release \
	--security-token $$GH_LOGIN \
	--user strongdm \
	--repo comply \
	--tag $(VERSION) \
	--name $(VERSION)

	github-release upload \
	--security-token $$GH_LOGIN \
	--user strongdm \
	--repo comply \
	--tag $(VERSION) \
	--name comply-$(VERSION)-darwin-amd64.tgz \
	--file dist/comply-$(VERSION)-darwin-amd64.tgz

	github-release upload \
	--security-token $$GH_LOGIN \
	--user strongdm \
	--repo comply \
	--tag $(VERSION) \
	--name comply-$(VERSION)-linux-amd64.tgz \
	--file dist/comply-$(VERSION)-linux-amd64.tgz

patch-release: patch release

patch: clean gitsem
	gitsem -m "increment patch for release" patch
	git push
	git push origin --tags

release-deps: gitsem gh-release

gitsem:
	go get -u github.com/Clever/gitsem

gh-release:
	go get -u github.com/aktau/github-release