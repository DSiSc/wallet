package utils

import (
	"flag"
	"fmt"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/wallet/accounts"
	"github.com/DSiSc/wallet/accounts/keystore"
	"github.com/DSiSc/wallet/common"
	local "github.com/DSiSc/wallet/core/types"
	web3cmn "github.com/alanchchen/web3go/common"
	"github.com/alanchchen/web3go/provider"
	"github.com/alanchchen/web3go/rpc"
	"github.com/alanchchen/web3go/web3"
	"io/ioutil"
	"os"
	"path/filepath"
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
	////construct payload
	//from := fmt.Sprintf("0x%x", *(tx.Data.From))
	//to := from
	//gas := "0x" + strconv.FormatInt(int64(tx.Data.GasLimit),16)
	//gasprice := "0x" + tx.Data.Price.String()
	//value := "0x" + tx.Data.Amount.String()
	//nonce := "0x" + strconv.FormatInt(int64(tx.Data.AccountNonce),16)
	//data := "" + string(tx.Data.Payload)
	//payload := fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_sendTransaction", "id": "0", "params": [{"from": "%s", "to": "%s","gas": "%s","gasPrice": "%s", "value": "%s","nonce": "%s","data": "%s"}]}`,
	//	from, to, gas, gasprice, value, nonce, data)
	//
	////new and send request
	//req, _ := http.NewRequest("POST", "http://127.0.0.1:47768/", strings.NewReader(payload))
	//client := http.Client{}
	//resp, err := client.Do(req)
	//if err != nil {
	//	log.Error("http send request failed: %v", err)
	//	return "", err
	//}
	//
	////resolve response
	//blob, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Error("read response msg has failed: %v", err)
	//	return "", err
	//}
	//recv := new(rpctypes.RPCResponse)
	//json.Unmarshal(blob, recv)
	//
	////var result cmn.Hash
	////json.Unmarshal(recv.Result, &result)
	////txHash := fmt.Sprintf("0x%x\n", result)
	//return "", nil

	hostname := flag.String("hostname", "127.0.0.1", "The ethereum client RPC host")
	port := flag.String("port", "47768", "The ethereum client RPC port")
	verbose := flag.Bool("verbose", true, "Print verbose messages")

	if *verbose {
		fmt.Printf("Connect to %s:%s\n", *hostname, *port)
	}

	provider := provider.NewHTTPProvider(*hostname+":"+*port, rpc.GetDefaultMethod())
	web3 := web3.NewWeb3(provider)

	//from := fmt.Sprintf("0x%x", *(tx.Data.From))
	//to := fmt.Sprintf("0x%x", *(tx.Data.Recipient))
	//gas := big.NewInt(int64(tx.Data.GasLimit))


	//req := &web3cmn.TransactionRequest{
	//	From:     web3cmn.NewAddress(web3cmn.HexToBytes(from)),
	//	To:       web3cmn.NewAddress(web3cmn.HexToBytes(to)),
	//	Gas:      gas,
	//	GasPrice: tx.Data.Price,
	//	Value:    tx.Data.Amount,
	//	Data:     tx.Data.Payload,
	//}
	req := &web3cmn.TransactionRequest{
		From:     "0xb60e8dd61c5d32be8058bb8eb970870f07233155",
		To:       "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
		Nonce:    "0x1",
		Gas:      "0x0",
		GasPrice: "0x0",
		Value:    "0x0",
		Data:     "",
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

