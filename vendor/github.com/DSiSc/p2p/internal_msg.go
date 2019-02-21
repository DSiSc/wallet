package p2p

import (
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/p2p/common"
	"github.com/DSiSc/p2p/message"
)

const (
	nilError = iota
)

// internal message type
type InternalMsg struct {
	From    *common.NetAddress
	To      *common.NetAddress
	Payload message.Message
	RespTo  chan interface{}
}

// peer disconect message ping message
type peerDisconnecMsg struct {
	err error
}

func (this *peerDisconnecMsg) MsgId() types.Hash {
	return message.EmptyHash
}

func (this *peerDisconnecMsg) MsgType() message.MessageType {
	return message.DISCONNECT_TYPE
}

func (this *peerDisconnecMsg) ResponseMsgType() message.MessageType {
	return message.NIL
}
