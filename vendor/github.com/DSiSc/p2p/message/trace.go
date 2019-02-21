package message

import (
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/p2p/common"
)

// DebugP2P is a test message, used to trace the message route in p2p.
type TraceMsg struct {
	ID     types.Hash           `json:"id"`
	Routes []*common.NetAddress `json:"routes"`
}

func (this *TraceMsg) MsgId() types.Hash {
	return this.ID
}

func (this *TraceMsg) MsgType() MessageType {
	return TRACE_TYPE
}

func (this *TraceMsg) ResponseMsgType() MessageType {
	return NIL
}
