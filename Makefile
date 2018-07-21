mocks:\
	mocks/entry_store.go

mocks/entry_store.go:
	moq -pkg=mocks -out=mocks/entry_store.go . entryStore

test: mocks
	go test . -race -p=1
