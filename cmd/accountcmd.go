package cmd

import (
	"fmt"
	"github.com/DSiSc/crypto-suite/crypto"
	"github.com/DSiSc/wallet/accounts"
	"github.com/DSiSc/wallet/accounts/keystore"
	"github.com/DSiSc/wallet/utils"
	"github.com/urfave/cli"
	"path/filepath"
)

var (

	AccountCommand = cli.Command{
		Name:     "account",
		Usage:    "Manage accounts",
		Category: "ACCOUNT COMMANDS",
		Description: `Manage accounts, list all existing accounts, import a private key into a new
account, create a new account or update an existing account.`,
		Subcommands: []cli.Command{
			{
				Name:   "list",
				Usage:  "Print summary of existing accounts",
				Action: utils.MigrateFlags(AccountList),
				Flags: []cli.Flag{
					utils.DataDirFlag,
					utils.KeyStoreDirFlag,
				},
				Description: `Print a short summary of all accounts`,
			},
			{
				Name:   "new",
				Usage:  "Create a new account",
				Action: utils.MigrateFlags(accountCreate),
				Flags: []cli.Flag{
					utils.DataDirFlag,
					utils.KeyStoreDirFlag,
					utils.PasswordFileFlag,
					utils.LightKDFFlag,
				},
				Description: `geth account new`,
			},
			{
				Name:         "update",
				Usage:        "Update an existing account",
				Action: utils.MigrateFlags(accountUpdate),
				Flags: []cli.Flag{
					utils.DataDirFlag,
					utils.KeyStoreDirFlag,
					utils.LightKDFFlag,
				},
				Description:  `Update an existing account`,
			},
			{
				Name:   "import",
				Usage:  "Import a private key into a new account",
				Action: utils.MigrateFlags(accountImport),
				Flags: []cli.Flag{
					utils.DataDirFlag,
					utils.KeyStoreDirFlag,
					utils.PasswordFileFlag,
					utils.LightKDFFlag,
				},
				ArgsUsage:   "<keyFile>",
				Description: `Imports an unencrypted private key from <keyfile> and creates a new account`,
			},
		},
	}
)

func AccountList(ctx *cli.Context) error {
	var index int

	dataDir := ctx.GlobalString(utils.DataDirFlag.Name)
	keyStoreDir := ctx.GlobalString(utils.KeyStoreDirFlag.Name)
	if keyStoreDir == "" {
		keyStoreDir =  keystore.KeyStoreScheme
	}
	keyStoreDir = filepath.Join(dataDir, keyStoreDir)

	manager, _, err := utils.MakeAccountManager(keyStoreDir)
	if err != nil {
		utils.Fatalf("Could not make account manager: %v", err)
	}

	for _, wallet := range manager.Wallets() {
		for _, account := range wallet.Accounts() {
			fmt.Printf("Account #%d: {%x} %s\n", index, account.Address, &account.URL)
			index++
		}
	}

	return nil
}

// tries unlocking the specified account a few times.
func unlockAccount(ctx *cli.Context, ks *keystore.KeyStore, address string, i int, passwords []string) (accounts.Account, string) {
	account, err := utils.MakeAddress(ks, address)
	if err != nil {
		utils.Fatalf("Could not list accounts: %v", err)
	}
	for trials := 0; trials < 3; trials++ {
		prompt := fmt.Sprintf("Unlocking account %s | Attempt %d/%d", address, trials+1, 3)
		password := getPassPhrase(prompt, false, i, passwords)
		err = ks.Unlock(account, password)
		if err == nil {
	    	//temp	log.Info("Unlocked account", "address", account.Address.Hex())
			return account, password
		}
		if err, ok := err.(*keystore.AmbiguousAddrError); ok {
			//temp			log.Info("Unlocked account", "address", account.Address.Hex())
			return ambiguousAddrRecovery(ks, err, password), password
		}
		if err != keystore.ErrDecrypt {
			// No need to prompt again if the error is not decryption-related.
			break
		}
	}
	// All trials expended to unlock account, bail out
	utils.Fatalf("Failed to unlock account %s (%v)", address, err)

	return accounts.Account{}, ""
}

func ambiguousAddrRecovery(ks *keystore.KeyStore, err *keystore.AmbiguousAddrError, auth string) accounts.Account {
	fmt.Printf("Multiple key files exist for address %x:\n", err.Addr)
	for _, a := range err.Matches {
		fmt.Println("  ", a.URL)
	}
	fmt.Println("Testing your passphrase against all of them...")
	var match *accounts.Account
	for _, a := range err.Matches {
		if err := ks.Unlock(a, auth); err == nil {
			match = &a
			break
		}
	}
	if match == nil {
		utils.Fatalf("None of the listed files could be unlocked.")
	}
	fmt.Printf("Your passphrase unlocked %s\n", match.URL)
	fmt.Println("In order to avoid this warning, you need to remove the following duplicate key files:")
	for _, a := range err.Matches {
		if a != *match {
			fmt.Println("  ", a.URL)
		}
	}
	return *match
}


// accountCreate creates a new account into the keystore defined by the CLI flags.
func accountCreate(ctx *cli.Context) error {

	dataDir := ctx.GlobalString(utils.DataDirFlag.Name)
	//get keyStoreDir from KeyStoreDirFlag, if not use the default value
	keyStoreDir := ctx.GlobalString(utils.KeyStoreDirFlag.Name)
	if keyStoreDir == "" {
		keyStoreDir =  keystore.KeyStoreScheme
	}
	keyStoreDir = filepath.Join(dataDir, keyStoreDir)

	scryptN, scryptP, keydir, err := utils.AccountConfig(keyStoreDir)

	if err != nil {
		utils.Fatalf("Failed to read configuration: %v", err)
	}

	password := getPassPhrase("Your new account is locked with a password. Please give a password. Do not forget this password.", true, 0, utils.MakePasswordList(ctx))

	address, err := keystore.StoreKey(keydir, password, scryptN, scryptP)

	if err != nil {
		utils.Fatalf("Failed to create account: %v", err)
	}
	fmt.Printf("Address: {%x}\n", address)
	return nil
}

func accountUpdate(ctx *cli.Context) error {
	if len(ctx.Args()) == 0 {
		utils.Fatalf("No accounts specified to update")
	}

	dataDir := ctx.GlobalString(utils.DataDirFlag.Name)
	keyStoreDir := ctx.GlobalString(utils.KeyStoreDirFlag.Name)
	if keyStoreDir == "" {
		keyStoreDir =  keystore.KeyStoreScheme
	}
	keyStoreDir = filepath.Join(dataDir, keyStoreDir)

	manager, _, _ := utils.MakeAccountManager(keyStoreDir)
	ks := manager.Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)

	for _, addr := range ctx.Args() {
		account, oldPassword := unlockAccount(ctx, ks, addr, 0, nil)
		newPassword := getPassPhrase("Please give a new password. Do not forget this password.", true, 0, nil)
		if err := ks.Update(account, oldPassword, newPassword); err != nil {
			utils.Fatalf("Could not update the account: %v", err)
		}
	}
	return nil
}

// getPassPhrase retrieves the password associated with an account, either fetched
// from a list of preloaded passphrases, or requested interactively from the user.
func getPassPhrase(prompt string, confirmation bool, i int, passwords []string) string {
	// If a list of passwords was supplied, retrieve from them
	if len(passwords) > 0 {
		if i < len(passwords) {
			return passwords[i]
		}
		return passwords[len(passwords)-1]
	}
	// Otherwise prompt the user for the password
	if prompt != "" {
		fmt.Println(prompt)
	}

	var password string
	fmt.Println("Passphrase: ")
	_, err := fmt.Scanln(&password)
	if err != nil {
		utils.Fatalf("Failed to read passphrase: %v", err)
	}

	if confirmation {
		fmt.Println("Repeat passphrase: ")
		var confirm string
		_, err := fmt.Scanln(&confirm)

		if err != nil {
			utils.Fatalf("Failed to read passphrase confirmation: %v", err)
		}
		if password != confirm {
			utils.Fatalf("Passphrases do not match")
		}
	}
	return password
}

func accountImport(ctx *cli.Context) error {
	keyfile := ctx.Args().First()
	if len(keyfile) == 0 {
		utils.Fatalf("keyfile must be given as argument")
	}
	key, err := crypto.LoadECDSA(keyfile)
	if err != nil {
		utils.Fatalf("Failed to load the private key: %v", err)
	}

	dataDir := ctx.GlobalString(utils.DataDirFlag.Name)
	keyStoreDir := ctx.GlobalString(utils.KeyStoreDirFlag.Name)
	if keyStoreDir == "" {
		keyStoreDir =  keystore.KeyStoreScheme
	}
	keyStoreDir = filepath.Join(dataDir, keyStoreDir)

	manager, _, _ := utils.MakeAccountManager(keyStoreDir)
	ks := manager.Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)

	passphrase := getPassPhrase("Your new account is locked with a password. Please give a password. Do not forget this password.", true, 0, utils.MakePasswordList(ctx))
	acct, err := ks.ImportECDSA(key, passphrase)
	if err != nil {
		utils.Fatalf("Could not create the account: %v", err)
	}
	fmt.Printf("Address: {%x}\n", acct.Address)
	return nil
}


