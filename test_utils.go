package godid

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

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
