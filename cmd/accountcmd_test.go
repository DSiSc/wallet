package main

import (
	"github.com/DSiSc/wallet/utils"
	"github.com/cespare/cp"
	"path/filepath"
	"runtime"
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

func TestAccountListEmpty(t *testing.T) {
	geth := runGeth(t, "account", "list")
	geth.ExpectExit()
}

func TestAccountList(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	//keystore := filepath.Join(datadir, "testdata")
	geth := runGeth(t, "account", "list", "--datadir", datadir)
	defer geth.ExpectExit()
	if runtime.GOOS == "windows" {
		geth.Expect(`
Account #0: {b69569609605b15ff631c3e85de107d862c6f134} keystore://{{.Datadir}}\keystore\UTC--2019-02-15T12-05-05.684713000Z--b69569609605b15ff631c3e85de107d862c6f134
Account #1: {f466859ead1932d743d622cb74fc058882e8648a} keystore://{{.Datadir}}\keystore\UTC--2019-02-18T04-32-03.186218000Z--94cdad6a9c62e418608f8ef5814821e74db3e331
Account #2: {289d485d9771714cce91d3393d764e1311907acc} keystore://{{.Datadir}}\keystore\zzz
`)
	} else {
		geth.Expect(`
Account #0: {b69569609605b15ff631c3e85de107d862c6f134} keystore://{{.Datadir}}/keystore/UTC--2019-02-15T12-05-05.684713000Z--b69569609605b15ff631c3e85de107d862c6f134
Account #1: {94cdad6a9c62e418608f8ef5814821e74db3e331} keystore://{{.Datadir}}/keystore/UTC--2019-02-18T04-32-03.186218000Z--94cdad6a9c62e418608f8ef5814821e74db3e331
Account #2: {f466859ead1932d743d622cb74fc058882e8648a} keystore://{{.Datadir}}/keystore/aaa
`)
	}
}

func TestAccountNew(t *testing.T) {
	geth := runGeth(t, "account", "new", "--lightkdf")
	defer geth.ExpectExit()
	geth.Expect(`
Your new account is locked with a password. Please give a password. Do not forget this password.
Passphrase: {{.InputLine "foobar"}}
Repeat passphrase: {{.InputLine "foobar"}}
`)
	geth.ExpectRegexp(`Address: \{[0-9a-f]{40}\}\n`)
}

func TestAccountNewBadRepeat(t *testing.T) {
	geth := runGeth(t, "account", "new", "--lightkdf")
	defer geth.ExpectExit()
	geth.Expect(`
Your new account is locked with a password. Please give a password. Do not forget this password.
Passphrase: {{.InputLine "something"}}
Repeat passphrase: {{.InputLine "something else"}}
Fatal: Passphrases do not match
`)
}

func TestAccountUpdate(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	geth := runGeth(t, "account", "update",
		"--datadir", datadir, "--lightkdf",
		"943bc733365b2a54488779916dd119c474f6e352")
	defer geth.ExpectExit()
	geth.Expect(`
Unlocking account 943bc733365b2a54488779916dd119c474f6e352 | Attempt 1/3
Passphrase: {{.InputLine "foobar"}}
Please give a new password. Do not forget this password.
Passphrase: {{.InputLine "foobar2"}}
Repeat passphrase: {{.InputLine "foobar2"}}
`)
}

func TestAccountImport(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	keystore := filepath.Join(datadir, "keystore")
	prifile := filepath.Join(keystore, "pri")
	geth := runGeth(t, "account", "import", "--datadir", datadir, "--keystore", keystore, prifile)
	defer geth.ExpectExit()
	geth.Expect(`
Your new account is locked with a password. Please give a password. Do not forget this password.
Passphrase: {{.InputLine "foobar"}}
Repeat passphrase: {{.InputLine "foobar"}}
`)
	geth.ExpectRegexp(`Address: \{[0-9a-f]{40}\}\n`)
}

func TestMakeAccountManager(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	keystoreDir := filepath.Join(datadir, "keystore")
	_, _, err := utils.MakeAccountManager(keystoreDir)
	if err != nil {
		utils.Fatalf("Could not make account manager: %v", err)
	}
}
