.DEFAULT_GOAL := comply
GO_SOURCES := $(shell find . -name '*.go')
THEME_SOURCES := $(shell find themes)

assets: $(THEME_SOURCES)
	@go get github.com/jteeuwen/go-bindata/...
	@go get github.com/elazarl/go-bindata-assetfs/...
	@go install github.com/elazarl/go-bindata-assetfs
	go-bindata-assetfs -o bindata.go -pkg theme -prefix themes themes/...
	mv bindata.go internal/theme/themes_bindata.go

comply: assets $(GO_SOURCES)
	@# $(eval VERSION := $(shell git describe --tags --always --dirty="-dev"))
	@# $(eval LDFLAGS := -ldflags='-X "github.com/strongdm/comply/internal/cli.Version=$(VERSION)"')
	go build $(LDFLAGS) github.com/strongdm/comply

dist: clean
	$(eval VERSION := $(shell git describe --tags --always --dirty="-dev"))
	$(eval LDFLAGS := -ldflags='-X "github.com/strongdm/comply/internal/cli.Version=$(VERSION)"')
	mkdir dist
	echo $(VERSION)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -gcflags=-trimpath=$(GOPATH) -asmflags=-trimpath=$(GOPATH) -ldflags '-extldflags "-static"' $(LDFLAGS) -o dist/comply-$(VERSION)-darwin-amd64 .
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -gcflags=-trimpath=$(GOPATH) -asmflags=-trimpath=$(GOPATH) -ldflags '-extldflags "-static"' $(LDFLAGS) -o dist/comply-$(VERSION)-linux-amd64 .
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -gcflags=-trimpath=$(GOPATH) -asmflags=-trimpath=$(GOPATH) -ldflags '-extldflags "-static"' $(LDFLAGS) -o dist/comply-$(VERSION)-windows-amd64.exe .
	cd dist && tar -czvf comply-$(VERSION)-darwin-amd64.tgz comply-$(VERSION)-darwin-amd64
	cd dist && tar -czvf comply-$(VERSION)-linux-amd64.tgz comply-$(VERSION)-linux-amd64
	cd dist && zip comply-$(VERSION)-windows-amd64.zip comply-$(VERSION)-windows-amd64.exe

brew: clean $(GO_SOURCES)
	$(eval VERSION := $(shell cat version))
	$(eval LDFLAGS := -ldflags='-X "github.com/strongdm/comply/internal/cli.Version=$(VERSION)"')
	mkdir bin
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -gcflags=-trimpath=$(GOPATH) -asmflags=-trimpath=$(GOPATH) $(LDFLAGS) -o bin/comply .

clean:
	rm -rf bin
	rm -rf dist
	rm -f comply

install: assets $(GO_SOURCES)
	go install github.com/strongdm/comply

push-assets: is-clean assets
	git commit -am "automated asset refresh (via Makefile)"
	git push

is-clean:
ifeq ($(strip $(shell git status --porcelain 2>/dev/null)),)
	# good to proceed
else
	@echo working directory must be clean to proceed
	@exit 1
endif

docker:
	cd build && docker build -t strongdm/pandoc .
	docker push strongdm/pandoc

cleanse:
	git checkout --orphan newbranch
	git add -A
	git commit -m "Initial commit"
	git branch -D master
	git branch -m master
	git push -f origin master
	git gc --aggressive --prune=all

release-env:
ifndef GH_LOGIN
	$(error GH_LOGIN must be set to a valid GitHub token)
endif
ifndef COMPLY_TAPDIR
	$(error COMPLY_TAPDIR must be set to the path of the comply homebrew tap repo)
endif

release: release-env dist release-deps
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

	@echo "Update homebrew formula with the following: "
	$(eval SHA := $(shell curl -s -L https://github.com/strongdm/comply/archive/$(VERSION).tar.gz |shasum -a 256|cut -d" " -f1))
	@echo "version $(VERSION) sha $(SHA)"
	cd $$COMPLY_TAPDIR && ./update.sh $(VERSION) $(SHA)

patch-release: release-env patch release
	$(eval VERSION := $(shell git describe --tags --always --dirty="-dev"))
	curl -X POST --data-urlencode 'payload={"channel": "#release", "username": "release", "text": "comply $(VERSION) released", "icon_emoji": ":shipit:"}' https://hooks.slack.com/services/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

minor-release: release-env minor release
	$(eval VERSION := $(shell git describe --tags --always --dirty="-dev"))
	curl -X POST --data-urlencode 'payload={"channel": "#release", "username": "release", "text": "comply $(VERSION) released", "icon_emoji": ":shipit:"}' https://hooks.slack.com/services/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

docker-release:
	docker build --build-arg COMPLY_VERSION=`cat VERSION` -t strongdm/comply .
	docker push strongdm/comply

patch: clean gitsem
	gitsem -m "increment patch for release (via Makefile)" patch
	git push
	git push origin --tags

minor: clean gitsem
	gitsem -m "increment minor for release (via Makefile)" minor
	git push
	git push origin --tags

release-deps: gitsem gh-release

gitsem:
	go get -u github.com/Clever/gitsem

gh-release:
	go get -u github.com/aktau/github-release