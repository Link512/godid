package godid

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func getTestBoltStore(t *testing.T) *boltStore {
	cleanupTestBoltStore()
	var err error
	testStore, err := newBoltStore(config{StorePath: "test.db"})
	require.NoError(t, err)
	return testStore
}

func cleanupTestBoltStore() {
	os.Remove("test.db")
}
