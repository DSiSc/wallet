package accounts

import (
	"testing"
)


func TestAuthNeededError(t *testing.T) {
	ErrLocked  := NewAuthNeededError("password or unlock")
	ErrLocked.Error()
}

