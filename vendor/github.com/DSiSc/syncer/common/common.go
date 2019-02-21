package common

import (
	"encoding/json"
	"github.com/DSiSc/craft/config"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/crypto-suite/crypto/sha3"
)

// HeaderHash calculate block's hash
func HeaderHash(header *types.Header) (hash types.Hash) {
	jsonByte, _ := json.Marshal(header)
	sumByte := Sum(jsonByte)
	copy(hash[:], sumByte)
	return
}

// Sum returns the first 32 bytes of hash of the bz.
func Sum(bz []byte) []byte {
	var alg string
	if value, ok := config.GlobalConfig.Load(config.HashAlgName); ok {
		alg = value.(string)
	} else {
		alg = "SHA256"
	}
	hasher := sha3.NewHashByAlgName(alg)
	hasher.Write(bz)
	hash := hasher.Sum(nil)
	return hash[:types.HashLength]
}
