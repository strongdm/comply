GO_SOURCES := $(shell find . -name '*.go')

comply: $(GO_SOURCES)
	go build github.com/strongdm/comply/cmd/comply