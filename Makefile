SHELL=/bin/bash -o pipefail

ifndef VIRTUALGO
    $(error No virtualgo workspace is active)
endif

GOBIN ?= $(GOPATH)/bin
INSTALLED = $(GOBIN)/did
VG_DEPS = $(VIRTUALGO_PATH)/last-ensure
GO_FILES = $(shell find . -name "*.go" | grep -v "_test.go$$" | xargs)
DEPS = $(VG_DEPS) $(GO_FILES)

all: $(INSTALLED)

mocks:\
	mock_entry_store.go

mock_entry_store.go:
	moq -out=mock_entry_store.go . entryStore

test: mocks
	GODID_TEST=1 go test . -race -p=1

$(INSTALLED): $(DEPS)
	go install ./did/
	@touch $(INSTALLED)

$(VG_DEPS): Gopkg.toml Gopkg.lock
	vg ensure -- -v
	@touch $(VG_DEPS)
