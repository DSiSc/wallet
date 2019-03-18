package utils

import (
	"github.com/DSiSc/wallet/accounts/keystore"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestNewApp(t *testing.T) {
	gitCommit := "wallet"
	// The app that holds all commands and flags.
	app1 := NewApp(gitCommit, "")

	assert.NotNil(t, app1)
}

func TestMakeAddress(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	ks := filepath.Join(datadir, "keystore")
	manager, _, _ := MakeAccountManager(ks)
	keystore := manager.Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
	_, err := MakeAddress(keystore, "2")

	assert.Equal(t, nil, err)
}
