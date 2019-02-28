package utils

import (
	"encoding/json"
	"fmt"
	cmn "github.com/DSiSc/apigateway/common"
	"github.com/DSiSc/apigateway/rpc/lib/types"
	"github.com/DSiSc/craft/log"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/wallet/accounts"
	"github.com/DSiSc/wallet/accounts/keystore"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
func SendTransaction(tx *types.Transaction) (string, error) {
	//construct payload
	from := fmt.Sprintf("0x%x", *(tx.Data.From))
	to := from
	gas := "0x" + strconv.FormatInt(int64(tx.Data.GasLimit),16)
	gasprice := "0x" + tx.Data.Price.String()
	value := "0x" + tx.Data.Amount.String()
	nonce := "0x" + strconv.FormatInt(int64(tx.Data.AccountNonce),16)
	data := "" + string(tx.Data.Payload)
	payload := fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_sendTransaction", "id": "0", "params": [{"from": "%s", "to": "%s","gas": "%s","gasPrice": "%s", "value": "%s","nonce": "%s","data": "%s"}]}`,
		from, to, gas, gasprice, value, nonce, data)

	//new and send request
	req, _ := http.NewRequest("POST", "http://127.0.0.1:47768/", strings.NewReader(payload))
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("http send request failed: %v", err)
		return "", err
	}

	//resolve response
	blob, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("read response msg has failed: %v", err)
		return "", err
	}
	recv := new(rpctypes.RPCResponse)
	json.Unmarshal(blob, recv)

	var result cmn.Hash
	json.Unmarshal(recv.Result, &result)
	txHash := fmt.Sprintf("0x%x\n", result)
	return txHash, nil
}
