SHELL=/bin/bash -o pipefail

VERSION = $(shell git tag --list --sort "-v:refname" | head -n 1 | sed "s/^v\(.*\)$$/\1/")
CURRENT_MAJOR = $(shell echo $(VERSION) | cut -d '.' -f 1)
CURRENT_MINOR = $(shell echo $(VERSION) | cut -d '.' -f 2)
CURRENT_BUG = $(shell echo $(VERSION) | cut -d '.' -f 3)

GOBIN ?= $(GOPATH)/bin
INSTALLED = $(GOBIN)/did

all: $(INSTALLED)

mocks:
	go generate

test: mocks
	GODID_TEST=1 go test . -race -p=1

cover: mocks
	GODID_TEST=1 go test -race -coverprofile=coverage.txt -covermode=atomic -p=1 .
	@sed -i.bak -e '/.*mock_entry_store\.go.*/d' ./coverage.txt
	@sed -i.bak -e '/.*config\.go.*/d' ./coverage.txt
	@rm coverage.txt.bak

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
