package utils

import (
	"fmt"
	"github.com/cespare/cp"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func tmpDatadirWithKeystore(t *testing.T) string {
	datadir := tmpdir(t)
	keystore := filepath.Join(datadir, "keystore")
	source := filepath.Join("..", "accounts", "keystore", "testdata", "keystore")
	if err := cp.CopyAll(keystore, source); err != nil {
		t.Fatal(err)
	}
	return datadir
}

func tmpdir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "geth-test")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestGetUnlockedKey(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	ks := filepath.Join(datadir, "keystore")
	ac, key, err := GetUnlockedKeyByDir("94cdad6a9c62e418608f8ef5814821e74db3e331", "", ks)
	fmt.Print(ac, key, err)
}
