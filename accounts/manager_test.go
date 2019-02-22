package accounts

import (
	"github.com/cespare/cp"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func tmpdir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "geth-test")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func tmpDatadirWithKeystore(t *testing.T) string {
	datadir := tmpdir(t)
	keystore := filepath.Join(datadir, "keystore")
	source := filepath.Join("..", "accounts", "keystore", "testdata")
	if err := cp.CopyAll(keystore, source); err != nil {
		t.Fatal(err)
	}
	return datadir
}

func TestMakeAccountManager(t *testing.T) {
	// Assemble the account manager and supported backends
	backends := []Backend{ }

	NewManager(backends...)
}
