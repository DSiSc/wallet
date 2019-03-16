package utils

import (
	"flag"
	"fmt"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/wallet/accounts"
	"github.com/DSiSc/wallet/accounts/keystore"
	"github.com/DSiSc/wallet/common"
	local "github.com/DSiSc/wallet/core/types"
	web3cmn "github.com/DSiSc/web3go/common"
	"github.com/DSiSc/web3go/provider"
	"github.com/DSiSc/web3go/rpc"
	"github.com/DSiSc/web3go/web3"
	"github.com/DSiSc/validator/tools"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
)




func statusOK(code int) bool { return code >= 200 && code <= 299 }

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

func Find(address string) (accounts.Account, error) {
	keyStoreDir := keystore.KeyStoreScheme
	return FindByDir(address, keyStoreDir)
}

func FindByDir(address string, keystoreDir string) (accounts.Account, error) {
	scryptN, scryptP, keydir, err := AccountConfig(keystoreDir)
	if err != nil {
		fmt.Printf("Failed to read configuration: %v", err)
		return accounts.Account{}, err
	}

	ks := keystore.NewKeyStore(keydir, scryptN, scryptP)

	//get account by stirng address
	account, err := MakeAddress(ks, address)
	if err != nil {
		fmt.Printf("Failed to make stirng address to account: %v", err)
		return accounts.Account{}, err
	}

	return ks.Find(account)
}

func SignTx(address string, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	keyStoreDir := keystore.KeyStoreScheme
	return SignTxByDir(address, tx, chainID, keyStoreDir, nil)
}

func SignTxByPassWord(tx *types.Transaction, password string) (*types.Transaction, error) {
	keyStoreDir := keystore.KeyStoreScheme

	addr := common.Address(*(tx.Data.From))
	address := addr.Hex()

	fmt.Println("address: ", address)
	return SignTxByDir(address, tx, nil, keyStoreDir, &password)
}

func SignTxByDir(address string, tx *types.Transaction, chainID *big.Int, keystoreDir string, password *string) (*types.Transaction, error) {
	scryptN, scryptP, keydir, err := AccountConfig(keystoreDir)
	if err != nil {
		fmt.Printf("Failed to read configuration: %v\n", err)
		return nil, err
	}

	ks := keystore.NewKeyStore(keydir, scryptN, scryptP)

	//get account by stirng address
	account, err := MakeAddress(ks, address)
	if err != nil {
		fmt.Printf("Failed to make stirng address to account: %v, address: %s\n", err, address)
		return nil, err
	}

	//unlock the account
	if password != nil {
		return ks.SignTxWithPassphrase(account, *password, tx, chainID)
	}

	return ks.SignTx(account, tx, chainID)
}


func GetUnlockedKey(address string, passphrase string) (accounts.Account, *keystore.Key, error) {
	keyStoreDir := keystore.KeyStoreScheme
	return GetUnlockedKeyByDir(address, passphrase, keyStoreDir)
}

func GetUnlockedKeyByDir(address string, passphrase string, keystoreDir string) (accounts.Account, *keystore.Key, error){
	scryptN, scryptP, keydir, err := AccountConfig(keystoreDir)
	if err != nil {
		fmt.Printf("Failed to read configuration: %v", err)
		return accounts.Account{}, nil, err
	}

	ks := keystore.NewKeyStore(keydir, scryptN, scryptP)

	//get account by stirng address
	account, err := MakeAddress(ks, address)
	if err != nil {
		fmt.Printf("Failed to make stirng address to account: %v", err)
		return accounts.Account{}, nil, err
	}

	ac, key, err := ks.GetDecryptedKey(account, passphrase)

	return ac, key, err
}



//Send a signed transaction
func SendTransaction(tx *types.Transaction) (common.Hash, error) {
	//format 0x string
	from := fmt.Sprintf("0x%x", *(tx.Data.From))
	to := from
	gas := "0x" + strconv.FormatInt(int64(tx.Data.GasLimit),16)
	gasprice := "0x" + tx.Data.Price.String()
	value := "0x" + tx.Data.Amount.String()
	data := ""

	if tx.Data.Payload != nil {
		data = "0x" + string(tx.Data.Payload)
	} else {
		data = ""
	}

	hostname := flag.String("hostname", "127.0.0.1", "The ethereum client RPC host")
	port := flag.String("port", "47768", "The ethereum client RPC port")
	verbose := flag.Bool("verbose", true, "Print verbose messages")

	if *verbose {
		fmt.Printf("Connect to %s:%s\n", *hostname, *port)
	}

	provider := provider.NewHTTPProvider(*hostname+":"+*port, rpc.GetDefaultMethod())
	web3 := web3.NewWeb3(provider)

	req := &web3cmn.TransactionRequest{
		From:     from,
		To:       to,
		Gas:      gas,
		GasPrice: gasprice,
		Value:    value,
		Data:     data,
	}

	hash, err := web3.Eth.SendTransaction(req)
	return common.Hash(hash), err
}

func SendTransactionWeb3(tx *web3cmn.TransactionRequest) (common.Hash, error) {
	hostname := flag.String("hostname", "127.0.0.1", "The ethereum client RPC host")
	port := flag.String("port", "47768", "The ethereum client RPC port")
	//port := flag.String("port", "8545", "The ethereum client RPC port")
	verbose := flag.Bool("verbose", true, "Print verbose messages")

	if *verbose {
		fmt.Printf("Connect to %s:%s\n", *hostname, *port)
	}

	provider := provider.NewHTTPProvider(*hostname+":"+*port, rpc.GetDefaultMethod())
	web3 := web3.NewWeb3(provider)

	hash, err := web3.Eth.SendTransaction(tx)
	return common.Hash(hash), err
}

func SendRawTransaction(tx *types.Transaction) (common.Hash, error) {
	hostname := flag.String("hostname", "127.0.0.1", "The ethereum client RPC host")
	port := flag.String("port", "47768", "The ethereum client RPC port")
	verbose := flag.Bool("verbose", false, "Print verbose messages")

	if *verbose {
		fmt.Printf("Connect to %s:%s\n", *hostname, *port)
	}

	provider := provider.NewHTTPProvider(*hostname+":"+*port, rpc.GetDefaultMethod())
	web3 := web3.NewWeb3(provider)

	txBytes, _ := local.EncodeToRLP(tx)
	hash, err := web3.Eth.SendRawTransaction(txBytes)

	return common.Hash(hash), err
}

func SendRawTransactionWeb3(txBytesStr string) (common.Hash, error) {
	hostname := flag.String("hostname", "127.0.0.1", "The ethereum client RPC host")
	port := flag.String("port", "47768", "The ethereum client RPC port")
	verbose := flag.Bool("verbose", false, "Print verbose messages")

	if *verbose {
		fmt.Printf("Connect to %s:%s\n", *hostname, *port)
	}

	provider := provider.NewHTTPProvider(*hostname+":"+*port, rpc.GetDefaultMethod())
	web3 := web3.NewWeb3(provider)

	bytes := tools.FromHex(txBytesStr)
	hash, err := web3.Eth.SendRawTransaction(bytes)

	return common.Hash(hash), err
}

func NewAccount(keyStoreDir string, password string) (common.Address, error) {
	if keyStoreDir == "" {
		dataDir := "./"
		if keyStoreDir == "" {
			keyStoreDir =  keystore.KeyStoreScheme
		}
		keyStoreDir = filepath.Join(dataDir, keyStoreDir)
	}

	scryptN, scryptP, keydir, err := AccountConfig(keyStoreDir)

	if err != nil {
		Fatalf("Failed to read configuration: %v", err)
	}

	address, err := keystore.StoreKey(keydir, password, scryptN, scryptP)
	if err != nil {
		Fatalf("Failed to create account: %v", err)
	}
	fmt.Printf("Address: {%x}\n", address)

	return address, err
}

func Unlock(ks *keystore.KeyStore, addr string, passphrase string) error {
	address := tools.HexToAddress(addr)
	account := accounts.Account{
		Address: common.Address(address),
	}
	return ks.Unlock(account, passphrase)
}

func Lock(ks *keystore.KeyStore, addr string) error {
	address := tools.HexToAddress(addr)
	return ks.Lock(common.Address(address))
}

func ListAccounts(keyStoreDir string) (error){
	var index int
	if keyStoreDir == "" {
		dataDir := "./"
		if keyStoreDir == "" {
			keyStoreDir =  keystore.KeyStoreScheme
		}
		keyStoreDir = filepath.Join(dataDir, keyStoreDir)
	}

	manager, _, err := MakeAccountManager(keyStoreDir)
	if err != nil {
		Fatalf("Could not make account manager: %v", err)
	}

	for _, wallet := range manager.Wallets() {
		for _, account := range wallet.Accounts() {
			fmt.Printf("Account #%d: {%x} %s\n", index, account.Address, &account.URL)
			index++
		}
	}
	return nil
}