package utils

import (
	"fmt"
	"github.com/DSiSc/wallet/common"
	local "github.com/DSiSc/wallet/core/types"
	"github.com/cespare/cp"
	"io/ioutil"
	"math/big"
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

func TestSendTransaction(t *testing.T) {
	nonce := uint64(1)
	from := common.Address{
		0xb2, 0x6f, 0x2b, 0x34, 0x2a, 0xab, 0x24, 0xbc, 0xf6, 0x3e,
		0xa2, 0x18, 0xc6, 0xa9, 0x27, 0x4d, 0x30, 0xab, 0x9a, 0x15,
	}
	to := from
	amount := big.NewInt(0)
	gaslimit := uint64(0)
	gasprice := big.NewInt(1000)
	// data := nil
	tx := local.NewTransaction(nonce, to, amount, gaslimit, gasprice, nil, from)
	txHash, err := SendTransaction(tx)
	if err != nil {
		fmt.Println("send tx has failed, ", err)
		return
	}
	fmt.Println(txHash)
}
