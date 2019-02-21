package common

import (
	"github.com/DSiSc/craft/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

var mockHash = types.Hash{
	0x1d, 0xcf, 0x7, 0xba, 0xfc, 0x42, 0xb0, 0x8d, 0xfd, 0x23, 0x9c, 0x45, 0xa4, 0xb9, 0x38, 0xd,
	0x8d, 0xfe, 0x5d, 0x6f, 0xa7, 0xdb, 0xd5, 0x50, 0xc9, 0x25, 0xb1, 0xb3, 0x4, 0xdc, 0xc5, 0x1c,
}

var mockHash1 = types.Hash{
	0x2d, 0xcf, 0x7, 0xba, 0xfc, 0x42, 0xb0, 0x8d, 0xfd, 0x23, 0x9c, 0x45, 0xa4, 0xb9, 0x38, 0xd,
	0x8d, 0xfe, 0x5d, 0x6f, 0xa7, 0xdb, 0xd5, 0x50, 0xc9, 0x25, 0xb1, 0xb3, 0x4, 0xdc, 0xc5, 0x1c,
}

func TestNewRingBuffer(t *testing.T) {
	assert := assert.New(t)
	ring := NewRingBuffer(1)
	assert.NotNil(ring)
}

func TestRingBuffer_AddElement(t *testing.T) {
	assert := assert.New(t)
	ring := NewRingBuffer(1)
	assert.NotNil(ring)
	ring.AddElement(mockHash, struct{}{})
	assert.True(ring.Exist(mockHash))
}

func TestRingBuffer_AddElement1(t *testing.T) {
	assert := assert.New(t)
	ring := NewRingBuffer(1)
	assert.NotNil(ring)
	ring.AddElement(mockHash, struct{}{})
	ring.AddElement(mockHash1, struct{}{})
	assert.False(ring.Exist(mockHash))
	assert.True(ring.Exist(mockHash1))
}

func TestRingBuffer_Exist(t *testing.T) {
	assert := assert.New(t)
	ring := NewRingBuffer(1)
	assert.NotNil(ring)
	ring.AddElement(mockHash, struct{}{})
	assert.True(ring.Exist(mockHash))
	assert.False(ring.Exist(mockHash1))
}
