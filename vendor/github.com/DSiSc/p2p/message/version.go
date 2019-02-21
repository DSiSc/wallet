package message

import "github.com/DSiSc/craft/types"

// Version version message
type Version struct {
	Version string `json:"version"`
	PortMe  int32  `json:"port_me"`
}

func (this *Version) MsgId() types.Hash {
	return EmptyHash
}

func (this *Version) MsgType() MessageType {
	return VERSION_TYPE
}

func (this *Version) ResponseMsgType() MessageType {
	return VERACK_TYPE
}

// Version ack message
type VersionAck struct {
}

func (this *VersionAck) MsgId() types.Hash {
	return EmptyHash
}

func (this *VersionAck) MsgType() MessageType {
	return VERACK_TYPE
}

func (this *VersionAck) ResponseMsgType() MessageType {
	return NIL
}
