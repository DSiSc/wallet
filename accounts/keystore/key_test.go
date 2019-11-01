package keystore

import (
	crand "crypto/rand"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKey_NewKeyForDirectICAP(t *testing.T) {
	key := NewKeyForDirectICAP(crand.Reader)

	assert.NotNil(t, key.PrivateKey)
	assert.NotNil(t, key.Address)
	assert.NotNil(t, key.Id)
}

func TestKey_MarshalJSON(t *testing.T) {
	key := NewKeyForDirectICAP(crand.Reader)
	_, err := key.MarshalJSON()

	assert.Equal(t, nil, err)
}

func TestKey_UnmarshalJSON(t *testing.T) {
	key := NewKeyForDirectICAP(crand.Reader)
	j, err := key.MarshalJSON()
	err = key.UnmarshalJSON(j)

	assert.Equal(t, nil, err)
}
