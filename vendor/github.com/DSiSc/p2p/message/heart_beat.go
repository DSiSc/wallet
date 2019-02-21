package message

import "github.com/DSiSc/craft/types"

// PingMsg ping message
type PingMsg struct {
	State uint64
}

func (this *PingMsg) MsgId() types.Hash {
	return EmptyHash
}

func (this *PingMsg) MsgType() MessageType {
	return PING_TYPE
}

func (this *PingMsg) ResponseMsgType() MessageType {
	return PONG_TYPE
}

// PongMsg pong message
type PongMsg struct {
	State uint64
}

func (this *PongMsg) MsgId() types.Hash {
	return EmptyHash
}

func (this *PongMsg) MsgType() MessageType {
	return PONG_TYPE
}

func (this *PongMsg) ResponseMsgType() MessageType {
	return NIL
}
