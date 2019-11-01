package types

import (
	"crypto/ecdsa"
	"encoding/json"
	ctypes "github.com/DSiSc/craft/types"
	cryp "github.com/DSiSc/crypto-suite/crypto"
	"github.com/DSiSc/wallet/common"
	"math/big"
	"sync/atomic"
)

// DataSignature common data signature struct
type DataSignature struct {
	// signer address
	Signer atomic.Value
	// Extra Data used when signing/verifying signature
	ExtraData uint64 `json:"extraData"`
	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`
}

// BlockSig return the signature info
func BlockSig(block *ctypes.Block) []*DataSignature {
	sigs := make([]*DataSignature, 0)
	for _, val := range block.Header.SigData {
		// unmarshal base sig info
		var sig DataSignature
		err := json.Unmarshal(val, sig)
		if err != nil {
			continue
		}
		recoverdSig := SignerInfo(sig.R, sig.S, sig.V, block.HeaderHash[:], sig.ExtraData)
		sigs = append(sigs, recoverdSig)
	}
	return sigs
}

// SignData signs the data using the given signer and private key
func SignData(hash ctypes.Hash, extraData uint64, prv *ecdsa.PrivateKey) (*DataSignature, error) {
	contentHash := sigHash(hash[:], extraData)
	sig, err := cryp.Sign(contentHash[:], prv)
	if err != nil {
		return nil, err
	}

	signer := NewEIP155Signer(big.NewInt(0))
	tx := &ctypes.Transaction{
		Data: ctypes.TxData{},
	}
	signedTx, err := WithSignature(tx, signer, sig)
	return &DataSignature{
		R: signedTx.Data.R,
		S: signedTx.Data.S,
		V: signedTx.Data.V,
	}, err
}

// SignerInfo return the signature info
func SignerInfo(r, s, v *big.Int, data []byte, extraData uint64) *DataSignature {
	addr, err := recoverPlain(sigHash(data, extraData), r, s, v, true)
	if err != nil {
		return nil
	}
	sig := &DataSignature{
		R: r,
		S: s,
		V: v,
	}
	sig.Signer.Store(TypeConvert(&addr))
	return sig
}

func sigHash(data []byte, extraData uint64) common.Hash {
	if extraData <= 0 {
		return common.BytesToHash(data)
	}
	return rlpHash([]interface{}{
		data,
		extraData,
	})

}
