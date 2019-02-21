package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNetAddress(t *testing.T) {
	assert := assert.New(t)
	addr := NewNetAddress("tcp", "127.0.0.1", 8080)
	assert.NotNil(addr)
	assert.Equal("tcp", addr.Protocol)
	assert.Equal("127.0.0.1", addr.IP)
	assert.Equal(int32(8080), addr.Port)
}

func TestNetAddress_Equal(t *testing.T) {
	assert := assert.New(t)
	addr1 := NewNetAddress("tcp", "127.0.0.1", 8080)
	assert.NotNil(addr1)
	addr2 := NewNetAddress("tcp", "127.0.0.1", 8080)
	assert.NotNil(addr2)
	assert.Equal(addr1, addr2)
	addr3 := NewNetAddress("tcp", "127.0.0.1", 8081)
	assert.NotNil(addr3)
	assert.NotEqual(addr1, addr3)
}

func TestParseNetAddress(t *testing.T) {
	assert := assert.New(t)
	addr, err := ParseNetAddress("tcp://127.0.0.1:8080")
	assert.Nil(err)
	assert.NotNil(addr)
	assert.Equal("tcp", addr.Protocol)
	assert.Equal("127.0.0.1", addr.IP)
	assert.Equal(int32(8080), addr.Port)
}

func TestParseNetAddress1(t *testing.T) {
	assert := assert.New(t)
	addr, err := ParseNetAddress("127.0.0.1:8080")
	assert.Nil(err)
	assert.NotNil(addr)
	assert.Equal("tcp", addr.Protocol)
	assert.Equal("127.0.0.1", addr.IP)
	assert.Equal(int32(8080), addr.Port)
}
