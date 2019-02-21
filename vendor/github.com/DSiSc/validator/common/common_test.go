package common

import (
	"bytes"
	"github.com/DSiSc/craft/types"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

func MockTransaction() *types.Transaction {
	to := &types.Address{
		0xb2, 0x6f, 0x2b, 0x34, 0x2a, 0xab, 0x24, 0xbc, 0xf6, 0x3e,
		0xa2, 0x18, 0xc6, 0xa9, 0x27, 0x4d, 0x30, 0xab, 0x9a, 0x15,
	}
	from := &types.Address{
		0x12, 0x6f, 0x2b, 0x34, 0x2a, 0xab, 0x24, 0xbc, 0xf6, 0x3e,
		0xa2, 0x18, 0xc6, 0xa9, 0x27, 0x4d, 0x30, 0xab, 0x9a, 0x15,
	}

	data := []byte{
		0xb2, 0x6f, 0x2b, 0x34, 0x2a, 0xab, 0x24, 0xbc, 0xf6, 0x3e,
	}

	tx := &types.Transaction{
		Data: types.TxData{
			AccountNonce: 0,
			Price:        new(big.Int).SetUint64(10),
			GasLimit:     100,
			Recipient:    to,
			From:         from,
			Amount:       new(big.Int).SetUint64(50),
			Payload:      data,
		},
	}
	return tx
}

var MockHash = types.Hash{
	0xb5, 0x42, 0x6d, 0xed, 0xca, 0x42, 0xd6, 0xa9, 0x1, 0x6d, 0x78, 0xd, 0x93, 0x66, 0x71, 0x81,
	0x15, 0x1, 0xf7, 0x18, 0x25, 0x6c, 0x2f, 0xd3, 0xa4, 0xb6, 0x36, 0x2f, 0x5d, 0xa6, 0xfb, 0xef,
}

var MockBlockHash = types.Hash{
	0xaf, 0x4e, 0x5b, 0xa3, 0x16, 0x97, 0x74, 0x6a, 0x26, 0x9d, 0x9b, 0x9e, 0xf1, 0x9d, 0xa8, 0xb3,
	0xf9, 0x32, 0x68, 0x16, 0xf4, 0x73, 0xd4, 0xb3, 0x6a, 0xaf, 0x2d, 0x6d, 0xfa, 0x82, 0xd9, 0x89,
}

var MockHeaderHash = types.Hash{
	0xcc, 0x88, 0x1c, 0x28, 0x30, 0x38, 0x50, 0x46, 0x2c, 0xcb, 0xae, 0xe5, 0xa4, 0x88, 0x85, 0x75,
	0xdf, 0xae, 0xd7, 0xd3, 0x39, 0x17, 0x9a, 0xfc, 0x9c, 0x4, 0x5e, 0xcd, 0x98, 0x8a, 0x39, 0xdd,
}

func TestTxHash(t *testing.T) {
	hash := TxHash(MockTransaction())
	assert.Equal(t, MockHash, hash)
	hash = TxHash(MockTransaction())
	assert.Equal(t, MockHash, hash)
}

func MockBlock() *types.Block {
	return &types.Block{
		Header: &types.Header{
			ChainID:       1,
			PrevBlockHash: MockHash,
			StateRoot:     MockHash,
			TxRoot:        MockHash,
			ReceiptsRoot:  MockHash,
			Height:        1,
			Timestamp:     uint64(time.Date(2018, time.August, 28, 0, 0, 0, 0, time.UTC).Unix()),
		},
		Transactions: make([]*types.Transaction, 0),
	}
}

func TestBlockHash(t *testing.T) {
	assert := assert.New(t)
	block := MockBlock()
	var tmp types.Hash

	header := block.Header
	assert.True(bytes.Equal(tmp[:], header.MixDigest[:]))
	assert.True(bytes.Equal(tmp[:], block.HeaderHash[:]))

	digest := HeaderDigest(header)
	exceptDigest := types.Hash{
		0xf9, 0x3, 0xec, 0x7d, 0x3, 0xe, 0xcb, 0xe, 0xbb, 0x77, 0x1a, 0xd9, 0xc0, 0xfa, 0xa8, 0x5e,
		0xde, 0xe7, 0x84, 0xd3, 0x9d, 0xc1, 0xec, 0xf4, 0xbf, 0x90, 0xbb, 0x2a, 0x51, 0x30, 0x97, 0x97,
	}
	assert.Equal(exceptDigest, digest)

	assert.True(bytes.Equal(tmp[:], block.HeaderHash[:]))
	headerHash := HeaderHash(block)
	exceptHeaderHash := types.Hash{
		0xf9, 0x3, 0xec, 0x7d, 0x3, 0xe, 0xcb, 0xe, 0xbb, 0x77, 0x1a, 0xd9, 0xc0, 0xfa, 0xa8, 0x5e,
		0xde, 0xe7, 0x84, 0xd3, 0x9d, 0xc1, 0xec, 0xf4, 0xbf, 0x90, 0xbb, 0x2a, 0x51, 0x30, 0x97, 0x97,
	}
	assert.Equal(exceptHeaderHash, headerHash)
	block.HeaderHash = HeaderHash(block)
	headerHash = HeaderHash(block)
	assert.Equal(exceptHeaderHash, headerHash)
}

func TestNewRefAddress(t *testing.T) {
	mock := types.Address{
		0xb2, 0x6f, 0x2b, 0x34, 0x2a, 0xab, 0x24, 0xbc, 0xf6, 0x3e,
		0xa2, 0x18, 0xc6, 0xa9, 0x27, 0x4d, 0x30, 0xab, 0x9a, 0x15,
	}
	refAddr := NewRefAddress(mock)
	assert.NotNil(t, refAddr)
	addr := refAddr.Address()
	assert.Equal(t, mock, addr)
}

func TestHashToByte(t *testing.T) {
	bytes := HashToByte(MockBlockHash)
	var exceptByte = []byte{
		0xaf, 0x4e, 0x5b, 0xa3, 0x16, 0x97, 0x74, 0x6a, 0x26, 0x9d, 0x9b, 0x9e, 0xf1, 0x9d, 0xa8, 0xb3,
		0xf9, 0x32, 0x68, 0x16, 0xf4, 0x73, 0xd4, 0xb3, 0x6a, 0xaf, 0x2d, 0x6d, 0xfa, 0x82, 0xd9, 0x89,
	}
	assert.Equal(t, exceptByte, bytes)
}

func TestByteToHash(t *testing.T) {
	var byteSrc = []byte{
		0xaf, 0x4e, 0x5b, 0xa3, 0x16, 0x97, 0x74, 0x6a, 0x26, 0x9d, 0x9b, 0x9e, 0xf1, 0x9d, 0xa8, 0xb3,
		0xf9, 0x32, 0x68, 0x16, 0xf4, 0x73, 0xd4, 0xb3, 0x6a, 0xaf, 0x2d, 0x6d, 0xfa, 0x82, 0xd9, 0x89,
	}
	hash := ByteToHash(byteSrc)
	assert.Equal(t, MockBlockHash, hash)
}
