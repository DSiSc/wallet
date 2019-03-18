// Copyright(c) 2018 DSiSc Group. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"github.com/DSiSc/crypto-suite/rlp"
	"io"
	"math/big"

	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/wallet/common"
)

func CopyBytes(b []byte) (copiedBytes []byte) {
	if b == nil {
		return nil
	}
	copiedBytes = make([]byte, len(b))
	copy(copiedBytes, b)

	return
}

func TypeConvert(a *common.Address) *types.Address {
	var address types.Address
	copy(address[:], a[:])
	return &address
}

// New a transaction
func newTransaction(nonce uint64, to *common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte, from *common.Address) *types.Transaction {
	if len(data) > 0 {
		data = CopyBytes(data)
	}
	d := types.TxData{
		AccountNonce: nonce,
		Recipient:    TypeConvert(to),
		From:         TypeConvert(from),
		Payload:      data,
		Amount:       new(big.Int),
		GasLimit:     gasLimit,
		Price:        new(big.Int),
		V:            new(big.Int),
		R:            new(big.Int),
		S:            new(big.Int),
	}
	if amount != nil {
		d.Amount.Set(amount)
	}
	if gasPrice != nil {
		d.Price.Set(gasPrice)
	}

	return &types.Transaction{Data: d}
}

func NewTransaction(nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte, from common.Address) *types.Transaction {
	return newTransaction(nonce, &to, amount, gasLimit, gasPrice, data, &from)
}

func ChainId(tx *types.Transaction) *big.Int {
	return deriveChainId(tx.Data.V)
}

func Protected(tx *types.Transaction) bool {
	return isProtectedV(tx.Data.V)
}

func isProtectedV(V *big.Int) bool {
	if V.BitLen() <= 8 {
		v := V.Uint64()
		return v != 27 && v != 28

	}
	// anything not 27 or 28 are considered unprotected
	return true

}

// WithSignature returns a new transaction with the given signature.
// This signature needs to be formatted as described in the yellow paper (v+27).
func WithSignature(tx *types.Transaction, signer Signer, sig []byte) (*types.Transaction, error) {
	r, s, v, err := signer.SignatureValues(tx, sig)
	if err != nil {
		return nil, err

	}

	cpy := &types.Transaction{Data: tx.Data}
	cpy.Data.R, cpy.Data.S, cpy.Data.V = r, s, v
	return cpy, nil

}

// EncodeRLP implements rlp.Encoder
func EncodeRLP(tx *types.Transaction,w io.Writer) error {
	return rlp.Encode(w, tx.Data)
}

// EncodeToBytes returns the RLP encoding of val.
func EncodeToRLP(tx *types.Transaction) ([]byte, error){
	return rlp.EncodeToBytes(tx)
}
