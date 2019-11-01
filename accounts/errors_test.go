package accounts

import (
	"fmt"
	"testing"
)

func TestAuthNeededError(t *testing.T) {
	ErrLocked := NewAuthNeededError("password or unlock")
	str := ErrLocked.Error()
	fmt.Print(str)
}
