// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package evm

import (
	"math/big"

	"github.com/DSiSc/blockchain"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/evm-NG/common/math"
	"github.com/DSiSc/evm-NG/util"
)

// NewEVMContext creates a new context for use in the EVM.
func NewEVMContext(tx types.Transaction, header *types.Header, chain *blockchain.BlockChain, author types.Address) Context {
	var beneficiary types.Address
	if (types.Address{} == author) {
		// TODO: Initially we specify a zero addressWWW
		beneficiary = util.HexToAddress("0x0000000000000000000000000000000000000000")
	}

	return Context{
		CanTransfer: CanTransfer,
		Transfer:    Transfer,
		GetHash:     GetHashFn(header, chain),
		Origin:      *tx.Data.From,
		Coinbase:    beneficiary,
		BlockNumber: new(big.Int).SetUint64(header.Height),
		Time:        new(big.Int).SetUint64(header.Timestamp),
		// TODO: Initially we specify a fixed difficulty
		Difficulty: new(big.Int).SetUint64(0x20000),
		// TODO: Initially we will not specify a precise gas limit
		GasLimit: uint64(math.MaxInt64),
		GasPrice: new(big.Int).Set(tx.Data.Price),
	}
}

// GetHashFn returns a GetHashFunc which retrieves header hashes by number
func GetHashFn(ref *types.Header, chain *blockchain.BlockChain) func(n uint64) types.Hash {
	var cache map[uint64]types.Hash

	return func(n uint64) types.Hash {
		// If there's no hash cache yet, make one
		if cache == nil {
			cache = map[uint64]types.Hash{
				ref.Height - 1: ref.PrevBlockHash,
			}
		}
		// Try to fulfill the request from the cache
		if hash, ok := cache[n]; ok {
			return hash
		}
		height := ref.Height - 1
		for {
			block, err := chain.GetBlockByHeight(height)
			if nil != err || nil == block.Header {
				break
			}
			cache[height-1] = block.Header.PrevBlockHash
			if n == height-1 {
				return block.Header.PrevBlockHash
			}
			height = height - 1
		}

		return types.Hash{}
	}
}

// CanTransfer checks whether there are enough funds in the address' account to make a transfer.
// This does not take the necessary gas in to account to make the transfer valid.
func CanTransfer(db *blockchain.BlockChain, addr types.Address, amount *big.Int) bool {
	return db.GetBalance(addr).Cmp(amount) >= 0
}

// Transfer subtracts amount from sender and adds amount to recipient using the given Db
func Transfer(db *blockchain.BlockChain, sender, recipient types.Address, amount *big.Int) {
	db.SubBalance(sender, amount)
	db.AddBalance(recipient, amount)
}
