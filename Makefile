mocks:\
	mock_entry_store.go

mock_entry_store.go:
	moq -out=mock_entry_store.go . entryStore

test: mocks
	go test . -race -p=1
