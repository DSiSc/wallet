package p2p

import (
	"github.com/DSiSc/blockchain"
	"sync/atomic"
)

var localState atomic.Value

func init() {
	localState.Store(uint64(0))
}

// LocalState get local current state
func LocalState() uint64 {
	bc, err := blockchain.NewLatestStateBlockChain()
	if err != nil {
		return localState.Load().(uint64)
	}
	currentHeight := bc.GetCurrentBlockHeight()
	localState.Store(currentHeight)
	return currentHeight
}
