package utils

import (
	"fmt"
	"github.com/DSiSc/wallet/common"
	local "github.com/DSiSc/wallet/core/types"
	"github.com/cespare/cp"
	"github.com/stretchr/testify/assert"
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

func TestStatusOK(t *testing.T) {
	result := statusOK(201)
	assert.Equal(t,true, result)

	result = statusOK(300)
	assert.Equal(t,false, result)
}
func TestFind(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	ks := filepath.Join(datadir, "keystore")
	acc, err := FindByDir("94cdad6a9c62e418608f8ef5814821e74db3e331", ks)
	assert.Equal(t, nil, err)
	assert.NotNil(t, acc)
}

func TestGetUnlockedKey(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	ks := filepath.Join(datadir, "keystore")
	_, _, err := GetUnlockedKey("94cdad6a9c62e418608f8ef5814821e74db3e331", "")
	assert.NotNil(t, err)

	_, _, err = GetUnlockedKeyByDir("94cdad6a9c62e418608f8ef5814821e74db3e331", "", ks)
	assert.Equal(t, nil, err)
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

//func TestSendRawTransaction(t *testing.T) {
//	nonce := uint64(1)
//	from := common.Address{
//		0xb2, 0x6f, 0x2b, 0x34, 0x2a, 0xab, 0x24, 0xbc, 0xf6, 0x3e,
//		0xa2, 0x18, 0xc6, 0xa9, 0x27, 0x4d, 0x30, 0xab, 0x9a, 0x15,
//	}
//	to := from
//	amount := big.NewInt(0)
//	gaslimit := uint64(0)
//	gasprice := big.NewInt(1000)
//	// data := nil
//	tx := local.NewTransaction(nonce, to, amount, gaslimit, gasprice, nil, from)
//	txHash, err := SendRawTransaction(tx)
//	if err != nil {
//		fmt.Println("send raw tx has failed, ", err)
//		return
//	}
//
//	expectHash := common.Hash{
//		105, 24, 76, 225, 150, 125, 28, 144, 68, 17, 185, 70, 162, 62, 105, 42,
//		16, 46, 238, 27, 148, 229, 81, 36, 136, 115, 27, 151, 68, 77, 195, 216,
//	}
//	//assert.Equal(t, expectHash, txHash)
//	fmt.Println(expectHash, txHash)
//}

func TestMakeAccountManager(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	ks := filepath.Join(datadir, "keystore")

	am, _, _ := MakeAccountManager(ks)

	assert.NotNil(t, am)
}

func TestAccountConfig(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	ks := filepath.Join(datadir, "keystore")

	_, _, _, err := AccountConfig(ks)
	assert.Equal(t, nil, err)
}

