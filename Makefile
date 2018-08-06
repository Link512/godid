SHELL=/bin/bash -o pipefail

VERSION = $(shell git tag --list --sort "-v:refname" | head -n 1 | sed "s/^v\(.*\)$$/\1/")
CURRENT_MAJOR = $(shell echo $(VERSION) | cut -d '.' -f 1)
CURRENT_MINOR = $(shell echo $(VERSION) | cut -d '.' -f 2)
CURRENT_BUG = $(shell echo $(VERSION) | cut -d '.' -f 3)

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

mock_entry_store.go: types.go
	moq -out=mock_entry_store.go . entryStore

test: mocks
	GODID_TEST=1 go test . -race -p=1

cover: mocks
	GODID_TEST=1 go test -race -coverprofile=coverage.txt -covermode=atomic -p=1 .
	@sed -i -e '/.*mock_entry_store\.go.*/d' ./coverage.txt

publish:
	@if [ "$(V)" = "" ]; then echo "You shouldn't be calling this directly, use publish-[major|minor|bug]"; exit 1; fi
	git tag -a -m "Bump version" v$(V)
	git push --follow-tags

update-master:
	git checkout master
	git pull

publish-major: update-master
	@make publish V=$$(($(CURRENT_MAJOR) + 1)).0.0
publish-minor: update-master
	@make publish V=$(CURRENT_MAJOR).$$(($(CURRENT_MINOR) + 1)).0
publish-bug: update-master
	@make publish V=$(CURRENT_MAJOR).$(CURRENT_MINOR).$$(($(CURRENT_BUG) + 1))

$(INSTALLED): $(DEPS)
	go install ./did/
	@touch $(INSTALLED)

$(VG_DEPS): Gopkg.toml Gopkg.lock
	vg ensure -- -v
	@touch $(VG_DEPS)
