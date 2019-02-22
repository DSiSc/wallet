package accounts

import (
	"github.com/ethereum/go-ethereum/accounts"
	"testing"
)


func TestAuthNeededError(t *testing.T) {
	ErrLocked  := accounts.NewAuthNeededError("password or unlock")
	ErrLocked.Error()
}

