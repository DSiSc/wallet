package message

import (
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/p2p/common"
)

// Transaction message
type Transaction struct {
	Tx *types.Transaction `json:"tx"`
}

func (this *Transaction) MsgId() types.Hash {
	return common.TxHash(this.Tx)
}

func (this *Transaction) MsgType() MessageType {
	return TX_TYPE
}

func (this *Transaction) ResponseMsgType() MessageType {
	return NIL
}
