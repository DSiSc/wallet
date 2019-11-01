package types

import (
	"github.com/DSiSc/craft/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSignerInfo(t *testing.T) {
	assert := assert.New(t)

	dataHash := types.Hash{}
	dataHash1 := rlpHash([]byte("Hello World"))
	copy(dataHash[:], dataHash1[:])
	key, addr := DefaultTestKey()

	// no extra data
	sigData, err := SignData(dataHash, 0, key)
	assert.Nil(err)
	assert.NotNil(sigData)

	sig := SignerInfo(sigData.R, sigData.S, sigData.V, dataHash[:], 0)
	assert.Equal(sig.Signer.Load().(*types.Address)[:], addr[:])

	// extra data 1
	sigData, err = SignData(dataHash, 1, key)
	assert.Nil(err)
	assert.NotNil(sigData)

	sig = SignerInfo(sigData.R, sigData.S, sigData.V, dataHash[:], 1)
	assert.Equal(sig.Signer.Load().(*types.Address)[:], addr[:])
}
