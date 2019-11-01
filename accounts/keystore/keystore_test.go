package keystore

import (
	"crypto/ecdsa"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/crypto-suite/common"
	"github.com/DSiSc/monkey"
	"github.com/DSiSc/wallet/accounts"
	ctypes "github.com/DSiSc/wallet/core/types"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/big"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

var testSigData = make([]byte, 32)

func TestKeyStore(t *testing.T) {
	dir, ks := tmpKeyStore(t, true)
	defer os.RemoveAll(dir)

	a, err := ks.NewAccount("foo")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(a.URL.Path, dir) {
		t.Errorf("account file %s doesn't have dir prefix", a.URL)
	}
	stat, err := os.Stat(a.URL.Path)
	if err != nil {
		t.Fatalf("account file %s doesn't exist (%v)", a.URL, err)
	}
	if runtime.GOOS != "windows" && stat.Mode() != 0600 {
		t.Fatalf("account file has wrong mode: got %o, want %o", stat.Mode(), 0600)
	}
	if !ks.HasAddress(a.Address) {
		t.Errorf("HasAccount(%x) should've returned true", a.Address)
	}
	if err := ks.Update(a, "foo", "bar"); err != nil {
		t.Errorf("Update error: %v", err)
	}
	if err := ks.Delete(a, "bar"); err != nil {
		t.Errorf("Delete error: %v", err)
	}
	if common.FileExist(a.URL.Path) {
		t.Errorf("account file %s should be gone after Delete", a.URL)
	}
	if ks.HasAddress(a.Address) {
		t.Errorf("HasAccount(%x) should've returned true after Delete", a.Address)
	}
}

func TestSign(t *testing.T) {
	dir, ks := tmpKeyStore(t, true)
	defer os.RemoveAll(dir)

	pass := "" // not used but required by API
	a1, err := ks.NewAccount(pass)
	if err != nil {
		t.Fatal(err)
	}
	if err := ks.Unlock(a1, ""); err != nil {
		t.Fatal(err)
	}
	if _, err := ks.SignHash(accounts.Account{Address: a1.Address}, testSigData); err != nil {
		t.Fatal(err)
	}
}

func TestSignTx(t *testing.T) {
	dir, ks := tmpKeyStore(t, true)
	defer os.RemoveAll(dir)

	monkey.Patch(ctypes.SignTx, func(tx *types.Transaction, s ctypes.Signer, prv *ecdsa.PrivateKey) (*types.Transaction, error) {
		return &types.Transaction{}, nil
	})

	pass := ""
	a1, err := ks.NewAccount(pass)
	if err != nil {
		t.Fatal(err)
	}
	if err := ks.Unlock(a1, ""); err != nil {
		t.Fatal(err)
	}

	if _, err := ks.SignTx(accounts.Account{Address: a1.Address}, &types.Transaction{}, big.NewInt(10)); err != nil {
		t.Fatal(err)
	}

	monkey.Unpatch(ctypes.SignTx)
}

func TestSignWithPassphrase(t *testing.T) {
	dir, ks := tmpKeyStore(t, true)
	defer os.RemoveAll(dir)

	pass := "passwd"
	acc, err := ks.NewAccount(pass)
	if err != nil {
		t.Fatal(err)
	}

	if _, unlocked := ks.unlocked[acc.Address]; unlocked {
		t.Fatal("expected account to be locked")
	}

	_, err = ks.SignHashWithPassphrase(acc, pass, testSigData)
	if err != nil {
		t.Fatal(err)
	}

	if _, unlocked := ks.unlocked[acc.Address]; unlocked {
		t.Fatal("expected account to be locked")
	}

	if _, err = ks.SignHashWithPassphrase(acc, "invalid passwd", testSigData); err == nil {
		t.Fatal("expected SignHashWithPassphrase to fail with invalid password")
	}
}

func TestTimedUnlock(t *testing.T) {
	dir, ks := tmpKeyStore(t, true)
	defer os.RemoveAll(dir)

	pass := "foo"
	a1, err := ks.NewAccount(pass)
	if err != nil {
		t.Fatal(err)
	}

	// Signing without passphrase fails because account is locked
	_, err = ks.SignHash(accounts.Account{Address: a1.Address}, testSigData)
	if err != ErrLocked {
		t.Fatal("Signing should've failed with ErrLocked before unlocking, got ", err)
	}

	// Signing with passphrase works
	if err = ks.TimedUnlock(a1, pass, 100*time.Millisecond); err != nil {
		t.Fatal(err)
	}

	// Signing without passphrase works because account is temp unlocked
	_, err = ks.SignHash(accounts.Account{Address: a1.Address}, testSigData)
	if err != nil {
		t.Fatal("Signing shouldn't return an error after unlocking, got ", err)
	}

	// Signing fails again after automatic locking
	time.Sleep(250 * time.Millisecond)
	_, err = ks.SignHash(accounts.Account{Address: a1.Address}, testSigData)
	if err != ErrLocked {
		t.Fatal("Signing should've failed with ErrLocked timeout expired, got ", err)
	}
}

func TestOverrideUnlock(t *testing.T) {
	dir, ks := tmpKeyStore(t, true)
	defer os.RemoveAll(dir)

	pass := "foo"
	a1, err := ks.NewAccount(pass)
	if err != nil {
		t.Fatal(err)
	}

	// Unlock indefinitely.
	if err = ks.TimedUnlock(a1, pass, 5*time.Minute); err != nil {
		t.Fatal(err)
	}

	// Signing without passphrase works because account is temp unlocked
	_, err = ks.SignHash(accounts.Account{Address: a1.Address}, testSigData)
	if err != nil {
		t.Fatal("Signing shouldn't return an error after unlocking, got ", err)
	}

	// reset unlock to a shorter period, invalidates the previous unlock
	if err = ks.TimedUnlock(a1, pass, 100*time.Millisecond); err != nil {
		t.Fatal(err)
	}

	// Signing without passphrase still works because account is temp unlocked
	_, err = ks.SignHash(accounts.Account{Address: a1.Address}, testSigData)
	if err != nil {
		t.Fatal("Signing shouldn't return an error after unlocking, got ", err)
	}

	// Signing fails again after automatic locking
	time.Sleep(250 * time.Millisecond)
	_, err = ks.SignHash(accounts.Account{Address: a1.Address}, testSigData)
	if err != ErrLocked {
		t.Fatal("Signing should've failed with ErrLocked timeout expired, got ", err)
	}
}

func TestSignRace(t *testing.T) {
	dir, ks := tmpKeyStore(t, true)
	defer os.RemoveAll(dir)

	// Create a test account.
	a1, err := ks.NewAccount("")
	if err != nil {
		t.Fatal("could not create the test account", err)
	}

	if err := ks.TimedUnlock(a1, "", 15*time.Millisecond); err != nil {
		t.Fatal("could not unlock the test account", err)
	}
	end := time.Now().Add(500 * time.Millisecond)
	for time.Now().Before(end) {
		if _, err := ks.SignHash(accounts.Account{Address: a1.Address}, testSigData); err == ErrLocked {
			return
		} else if err != nil {
			t.Errorf("Sign error: %v", err)
			return
		}
		time.Sleep(1 * time.Millisecond)
	}
	t.Errorf("Account did not lock within the timeout")
}

func TestWallets(t *testing.T) {
	dir, ks := tmpKeyStore(t, true)
	defer os.RemoveAll(dir)

	ks.Wallets()
}

func TestSignTxWithPassphrase(t *testing.T) {
	dir, ks := tmpKeyStore(t, true)
	defer os.RemoveAll(dir)

	pass := "" // not used but required by API
	a1, err := ks.NewAccount(pass)
	if err != nil {
		t.Fatal(err)
	}
	if err := ks.Unlock(a1, ""); err != nil {
		t.Fatal(err)
	}

	if _, err := ks.SignTxWithPassphrase(accounts.Account{Address: a1.Address},
		"", &types.Transaction{}, big.NewInt(10)); err != nil {
		t.Fatal(err)
	}
}

func TestUnlock(t *testing.T) {
	dir, ks := tmpKeyStore(t, true)
	defer os.RemoveAll(dir)

	pass := ""
	acc, err := ks.NewAccount(pass)
	if err != nil {
		t.Fatal(err)
	}

	//Unlock the addr
	ks.Unlock(acc, pass)
	if _, unlocked := ks.unlocked[acc.Address]; !unlocked {
		t.Fatal("expected account to be unlock")
	}

	//Lock the addr
	ks.Lock(acc.Address)
	if _, unlocked := ks.unlocked[acc.Address]; unlocked {
		t.Fatal("expected account to be locked")
	}

}

func TestExport(t *testing.T) {
	dir, ks := tmpKeyStore(t, true)
	defer os.RemoveAll(dir)

	pass := ""
	acc, err := ks.NewAccount(pass)
	assert.Equal(t, nil, err)

	jsonBytes, err1 := ks.Export(acc, pass, pass)
	assert.Equal(t, nil, err1)

	acc1, err2 := ks.Import(jsonBytes, pass, pass)
	assert.Equal(t, nil, err2)
	assert.Equal(t, acc.Address, acc1.Address)

}

func TestImportECDSA(t *testing.T) {
	dir, ks := tmpKeyStore(t, true)
	defer os.RemoveAll(dir)

	pass := ""
	acc, err := ks.NewAccount(pass)
	assert.Equal(t, nil, err)

	keyjson, err := ioutil.ReadFile(acc.URL.Path)
	assert.Equal(t, nil, err)

	//delete account
	err = ks.Delete(acc, pass)
	assert.Equal(t, nil, err)

	key, err := DecryptKey(keyjson, pass)
	assert.Equal(t, nil, err)

	//Get the privatekey
	_, err1 := ks.ImportECDSA(key.PrivateKey, "")
	assert.Equal(t, nil, err1)
}

func TestImportPreSaleKey(t *testing.T) {
	dir, ks := tmpKeyStore(t, true)
	defer os.RemoveAll(dir)

	keyJons := make([]byte, 20)

	_, err := ks.ImportPreSaleKey(keyJons, "")
	assert.Equal(t, nil, err)
}

func tmpKeyStore(t *testing.T, encrypted bool) (string, *KeyStore) {
	d, err := ioutil.TempDir("", "eth-keystore-test")
	if err != nil {
		t.Fatal(err)
	}
	new := func(kd string) *KeyStore { return nil }
	if encrypted {
		new = func(kd string) *KeyStore { return NewKeyStore(kd, veryLightScryptN, veryLightScryptP) }
	}
	return d, new(string(d))
}
