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
	"io/ioutil"
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

