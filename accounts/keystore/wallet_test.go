package keystore

import (
	"github.com/DSiSc/craft/types"
	"github.com/stretchr/testify/assert"
	"math/big"
	"os"
	"testing"
)

func TestKeystoreWallet_Accounts(t *testing.T) {
	dir, ks := tmpKeyStore(t, true)
	defer os.RemoveAll(dir)

	acc, err := ks.NewAccount("foo")
	assert.Equal(t, nil, err)

	ksWallet := keystoreWallet{account: acc, keystore: ks}
	accs := ksWallet.Accounts()
	assert.NotNil(t, accs)

	err = ksWallet.Close()
	assert.Equal(t, nil, err)
}

func TestKeystoreWallet_Status(t *testing.T) {
	dir, ks := tmpKeyStore(t, true)
	defer os.RemoveAll(dir)

	acc, err := ks.NewAccount("foo")
	assert.Equal(t, nil, err)

	ksWallet := keystoreWallet{account: acc, keystore: ks}
	_, err = ksWallet.Status()
	assert.Equal(t, nil, err)
}

func TestKeystoreWallet_Open(t *testing.T) {
	dir, ks := tmpKeyStore(t, true)
	defer os.RemoveAll(dir)

	acc, err := ks.NewAccount("foo")
	assert.Equal(t, nil, err)

	ksWallet := keystoreWallet{account: acc, keystore: ks}
	err = ksWallet.Open("")
	assert.Equal(t, nil, err)
}

func TestKeystoreWallet_Contains(t *testing.T) {
	dir, ks := tmpKeyStore(t, true)
	defer os.RemoveAll(dir)

	acc, err := ks.NewAccount("foo")
	assert.Equal(t, nil, err)

	ksWallet := keystoreWallet{account: acc, keystore: ks}
	flag := ksWallet.Contains(acc)
	assert.Equal(t, true, flag)
}

func TestKeystoreWallet_SignHash(t *testing.T) {
	dir, ks := tmpKeyStore(t, true)
	defer os.RemoveAll(dir)

	acc, err := ks.NewAccount("foo")
	assert.Equal(t, nil, err)
	ks.Unlock(acc, "foo")

	ksWallet := keystoreWallet{account: acc, keystore: ks}
	_, err = ksWallet.SignHash(acc, make([]byte, 32))
	assert.Equal(t, nil, err)
}

func TestKeystoreWallet_SignHashWithPassphrase(t *testing.T) {
	dir, ks := tmpKeyStore(t, true)
	defer os.RemoveAll(dir)

	acc, err := ks.NewAccount("foo")
	assert.Equal(t, nil, err)

	ksWallet := keystoreWallet{account: acc, keystore: ks}
	_, err = ksWallet.SignHashWithPassphrase(acc, "foo", make([]byte, 32))
	assert.Equal(t, nil, err)
}

func TestKeystoreWallet_SignTx(t *testing.T) {
	dir, ks := tmpKeyStore(t, true)
	defer os.RemoveAll(dir)

	acc, err := ks.NewAccount("foo")
	assert.Equal(t, nil, err)
	ks.Unlock(acc, "foo")

	ksWallet := keystoreWallet{account: acc, keystore: ks}
	tx := types.Transaction{}
	_, err = ksWallet.SignTx(acc, &tx, big.NewInt(10))
	assert.Equal(t, nil, err)
}

func TestKeystoreWallet_SignTxWithPassphrase(t *testing.T) {
	dir, ks := tmpKeyStore(t, true)
	defer os.RemoveAll(dir)

	acc, err := ks.NewAccount("foo")
	assert.Equal(t, nil, err)

	ksWallet := keystoreWallet{account: acc, keystore: ks}
	tx := types.Transaction{}
	_, err = ksWallet.SignTxWithPassphrase(acc, "foo", &tx, big.NewInt(10))
	assert.Equal(t, nil, err)
}