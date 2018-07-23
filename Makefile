mocks:\
	mock_entry_store.go

mock_entry_store.go:
	moq -out=mock_entry_store.go . entryStore

test: mocks
	GODID_TEST=1 go test . -race -p=1
