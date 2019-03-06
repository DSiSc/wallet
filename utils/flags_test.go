package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewApp(t *testing.T) {
	gitCommit := "wallet"
	// The app that holds all commands and flags.
	app1 := NewApp(gitCommit, "")

	assert.NotNil(t, app1)
}
