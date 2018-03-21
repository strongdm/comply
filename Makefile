.DEFAULT_GOAL := comply
GO_SOURCES := $(shell find . -name '*.go')
THEME_SOURCES := $(shell find themes)

assets: $(THEME_SOURCES)
	go-bindata-assetfs -pkg theme -ignore "\\.git" -prefix themes themes/...
	mv bindata_assetfs.go internal/theme/themes_bindata.go

comply: assets $(GO_SOURCES)
	go build github.com/strongdm/comply/cmd/comply

install: assets $(GO_SOURCES)
	go install github.com/strongdm/comply/cmd/comply