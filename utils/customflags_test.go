package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDirectoryFlagFunc(t *testing.T) {
	dir := new(DirectoryString)

	value := "keystore"
	dir.Value = value
	assert.Equal(t, value, dir.String())

	value = "keystore1"
	dir.Set(value)
	assert.Equal(t, value, dir.String())

	value = "~/keystore1"
	err := dir.Set(value)
	assert.Equal(t, nil, err)

	dirFlag := DirectoryFlag {
		Name:	"datadir",
		Usage:	"data directory",
		Value:	DirectoryString{"keystore"},
	}

	assert.Equal(t, "--datadir \"keystore\"	data directory", dirFlag.String())
	assert.Equal(t, "datadir", dirFlag.GetName())

}
