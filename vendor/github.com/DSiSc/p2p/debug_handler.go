package p2p

import (
	"fmt"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/p2p/common"
	"github.com/DSiSc/p2p/message"
	stCommon "github.com/DSiSc/p2p/tools/common"
	"github.com/DSiSc/p2p/tools/statistics/client"
	"strconv"
	"time"
)

// DebugHandler is a handler used to trace p2p message and report the p2p topo info.
type DebugHandler struct {
	p2p              *P2P
	center           types.EventCenter
	subs             map[types.EventType]types.Subscriber
	quitChan         chan interface{}
	statisticsClient *client.StatisticsClient
}

// NewDebugHandler create a new NewDebugHandler instance
func NewDebugHandler(p2p *P2P, center types.EventCenter, debugServer string) *DebugHandler {
	return &DebugHandler{
		p2p:              p2p,
		center:           center,
		quitChan:         make(chan interface{}),
		subs:             make(map[types.EventType]types.Subscriber),
		statisticsClient: client.NewStatisticsClient(debugServer),
	}
}

// Start start p2p debug handler
func (this *DebugHandler) Start() {
	this.subs[types.EventRecvNewMsg] = this.center.Subscribe(types.EventRecvNewMsg, this.RecvMsgEventSubscriber)
	this.subs[types.EventBroadCastMsg] = this.center.Subscribe(types.EventBroadCastMsg, this.BroadCastMsgEventSubscriber)
	go this.ReportNeighborHandler()
}

// Stop stop p2p debug handler
func (this *DebugHandler) Stop() {
	close(this.quitChan)
	for event, sub := range this.subs {
		this.center.UnSubscribe(event, sub)
	}
}

// RecvMsgEventSubscriber is the type of types.EventFunc, used to subscribe the p2p related event
func (this *DebugHandler) RecvMsgEventSubscriber(msg interface{}) {
	imsg := msg.(*InternalMsg)
	switch imsg.Payload.(type) {
	case *message.Block:
	case *message.BlockReq:
	case *message.Transaction:
	case *message.TraceMsg:
	default:
		return
	}
	if this.p2p.addrManager.OurAddresses()[0].Port == imsg.To.Port {
		this.reportMessage(this.p2p.config.DebugAddr+":"+strconv.Itoa(int(imsg.To.Port)), imsg, false)
	}
}

// EventSubscriber is the type of types.EventFunc, used to subscribe the p2p related event
func (this *DebugHandler) BroadCastMsgEventSubscriber(msg interface{}) {
	imsg := msg.(*InternalMsg)
	switch imsg.Payload.(type) {
	case *message.Block:
	case *message.BlockReq:
	case *message.Transaction:
	case *message.TraceMsg:
	default:
		return
	}
	if this.p2p.addrManager.OurAddresses()[0].Port == imsg.From.Port {
		this.reportMessage(this.p2p.config.DebugAddr+":"+strconv.Itoa(int(imsg.From.Port)), imsg, true)
	}
}

// ReportNeighborHandler report peer neighbor's handler
func (this *DebugHandler) ReportNeighborHandler() {
	// send trace message periodically
	timer := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-timer.C:
			this.reportNeighbors(this.p2p.config.DebugAddr+":"+strconv.Itoa(int(this.p2p.addrManager.OurAddresses()[0].Port)), this.p2p.GetPeers())
		case <-this.quitChan:
			return
		}
	}
}

// report message to trace server.
func (this *DebugHandler) reportMessage(localAddr string, msg *InternalMsg, isSend bool) {
	cmsg := &stCommon.ReportMsg{
		ReportPeer: localAddr,
	}
	if isSend {
		cmsg.From = localAddr
	} else {
		cmsg.From = addrString(msg.From)
	}
	this.statisticsClient.ReportMsg(fmt.Sprintf("%x", msg.Payload.MsgId()), cmsg)
}

// report peer's neighbor info
func (this *DebugHandler) reportNeighbors(localAddr string, peers []*Peer) {
	neighbors := make([]*stCommon.Neighbor, 0)
	for _, peer := range peers {
		neighbor := &stCommon.Neighbor{
			Address:  addrString(peer.GetAddr()),
			OutBound: peer.IsOutBound(),
		}
		neighbors = append(neighbors, neighbor)
	}
	this.statisticsClient.ReportNeighbors(localAddr, neighbors)
}

// format NetAddress to string
func addrString(addr *common.NetAddress) string {
	return addr.IP + ":" + strconv.Itoa(int(addr.Port))
}
