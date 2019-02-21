package evm

import (
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/evm-NG/util"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

// test new evm context
func TestNewEVMContext(t *testing.T) {
	assert := assert.New(t)
	msg := types.Transaction{
		Data: types.TxData{
			From:         &callerAddress,
			Recipient:    &contractAddress,
			AccountNonce: 0,
			Amount:       big.NewInt(0),
			GasLimit:     10000000000,
			Price:        big.NewInt(2),
			Payload:      nil,
		},
	}
	//(from Address, to *Address, nonce uint64, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte, checkNonce bool
	//types.NewMessage(
	//	callerAddress,
	//	&contractAddress,
	//	0,
	//	big.NewInt(0x5af3107a4000),
	//	0,
	//	big.NewInt(2),
	//	nil,
	//	false)
	header := &types.Header{
		PrevBlockHash: util.HexToHash(""),
		Height:        1,
		Timestamp:     1,
	}
	author := util.HexToAddress("0x0000000000000000000000000000000000000000")

	bc := mockPreBlockChain()
	context := NewEVMContext(msg, header, bc, author)
	assert.NotNil(context)
}

// test get hash func implemention
func TestGetHashFn(t *testing.T) {
	assert := assert.New(t)
	bc := mockPreBlockChain()
	cuBlock := bc.GetCurrentBlock()
	header := &types.Header{
		Height:        cuBlock.Header.Height + 1,
		PrevBlockHash: cuBlock.HeaderHash,
	}
	hashFunc := GetHashFn(header, bc)
	hash := hashFunc(cuBlock.Header.Height)
	assert.Equal(hash, cuBlock.HeaderHash)
}

// test can transfer function
func TestCanTransfer(t *testing.T) {
	assert := assert.New(t)
	address := util.HexToAddress("0x0000000000000000000000000000000000000000")
	bc := mockPreBlockChain()
	bc.SetBalance(address, big.NewInt(50))

	result := CanTransfer(bc, address, big.NewInt(10))
	assert.True(result)
}

// test transfer function
func TestTransfer(t *testing.T) {
	assert := assert.New(t)
	address1 := util.HexToAddress("0x0000000000000000000000000000000000000000")
	address2 := util.HexToAddress("0x0000000000000000000000000000000000000001")
	bc := mockPreBlockChain()
	bc.SetBalance(address1, big.NewInt(100))
	bc.SetBalance(address2, big.NewInt(100))

	Transfer(bc, address1, address2, big.NewInt(50))
	assert.Equal(big.NewInt(50), bc.GetBalance(address1))
	assert.Equal(big.NewInt(150), bc.GetBalance(address2))
}
