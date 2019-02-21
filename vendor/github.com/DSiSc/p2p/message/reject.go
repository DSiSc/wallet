package message

import "github.com/DSiSc/craft/types"

// RejectMsg reject message
type RejectMsg struct {
	Reason string
}

func (this *RejectMsg) MsgId() types.Hash {
	return EmptyHash
}

func (this *RejectMsg) MsgType() MessageType {
	return REJECT_TYPE
}

func (this *RejectMsg) ResponseMsgType() MessageType {
	return NIL
}
