.DEFAULT_GOAL := comply
GO_SOURCES := $(shell find . -name '*.go')
THEME_SOURCES := $(shell find themes)

assets: $(THEME_SOURCES)
	go-bindata-assetfs -pkg theme -prefix themes themes/...
	mv bindata_assetfs.go internal/theme/themes_bindata.go

comply: assets $(GO_SOURCES)
	go build github.com/strongdm/comply/cmd/comply

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

