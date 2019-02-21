package p2p

import (
	"github.com/DSiSc/p2p/common"
	"github.com/DSiSc/p2p/message"
)

type P2PAPI interface {
	// Start start p2p service
	Start() error

	// Stop stop p2p service
	Stop()

	// BroadCast broad cast message To all neighbor peers
	BroadCast(msg message.Message)

	// SendMsg send message to a peer
	SendMsg(peerAddr *common.NetAddress, msg message.Message) error

	// Gather gather newest data From p2p network
	Gather(peerFilter PeerFilter, reqMsg message.Message) error

	// MessageChan get p2p's message channel, (Messages sent To the server will eventually be placed in the message channel)
	MessageChan() <-chan *InternalMsg
}
