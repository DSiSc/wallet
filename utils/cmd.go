package utils

import (
	"github.com/DSiSc/wallet/accounts"
	"github.com/DSiSc/wallet/accounts/keystore"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func AccountConfig(keystoreDir string) (int, int, string, error) {
	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP

	var (
		keydir string
		err    error
	)
	switch {
	case filepath.IsAbs(keystoreDir):
		keydir = keystoreDir

	case keystoreDir != "":
		keydir, err = filepath.Abs(keystoreDir)
	}
	return scryptN, scryptP, keydir, err
}

func MakeAccountManager(keystoreDir string) (*accounts.Manager, string, error) {
	scryptN, scryptP, keydir, err := AccountConfig(keystoreDir)
	var ephemeral string
	if keydir == "" {
		// There is no datadir.
		keydir, err = ioutil.TempDir("", "wallet-keystore")
		ephemeral = keydir
	}

	if err != nil {
		return nil, "", err
	}
	if err := os.MkdirAll(keydir, 0700); err != nil {
		return nil, "", err
	}

	// Assemble the account manager and supported backends
	backends := []accounts.Backend{
		keystore.NewKeyStore(keydir, scryptN, scryptP),
	}

	return accounts.NewManager(backends...), ephemeral, nil
}


func GetUnlockedKey(address string, passphrase string) (accounts.Account, *keystore.Key, error) {
	keyStoreDir := keystore.KeyStoreScheme
	return GetUnlockedKeyByDir(address, passphrase, keyStoreDir)
}

func GetUnlockedKeyByDir(address string, passphrase string, keystoreDir string) (accounts.Account, *keystore.Key, error){
	scryptN, scryptP, keydir, err := AccountConfig(keystoreDir)
	if err != nil {
		log.Printf("Failed to read configuration: %v", err)
		return accounts.Account{}, nil, err
	}

	ks := keystore.NewKeyStore(keydir, scryptN, scryptP)

	//get account by stirng address
	account, err := MakeAddress(ks, address)
	if err != nil {
		log.Printf("Failed to make stirng address to account: %v", err)
		return accounts.Account{}, nil, err
	}

	ac, key, err := ks.GetDecryptedKey(account, passphrase)

	return ac, key, err
}
