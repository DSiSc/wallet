package syncer

import (
	"github.com/DSiSc/blockchain"
	"github.com/DSiSc/craft/log"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/p2p"
	"github.com/DSiSc/p2p/message"
	"github.com/DSiSc/syncer/common"
	"sync"
	"time"
)

// BlockSyncer block synchronize program
type BlockSyncer struct {
	blockSyncChan chan interface{}
	p2p           p2p.P2PAPI
	blockChain    *blockchain.BlockChain
	eventCenter   types.EventCenter
	sendChan      chan<- interface{}
	stallChan     chan types.Hash
	quitChan      chan interface{}
	subscribers   map[types.EventType]types.Subscriber
	lock          sync.RWMutex
	pendingBlocks map[types.Hash]interface{}
}

// NewBlockSyncer create block syncer instance.
func NewBlockSyncer(p2p p2p.P2PAPI, sendChan chan<- interface{}, eventCenter types.EventCenter) (*BlockSyncer, error) {
	blockChain, err := blockchain.NewLatestStateBlockChain()
	if err != nil {
		return nil, err
	}
	return &BlockSyncer{
		blockSyncChan: make(chan interface{}),
		p2p:           p2p,
		blockChain:    blockChain,
		sendChan:      sendChan,
		stallChan:     make(chan types.Hash),
		eventCenter:   eventCenter,
		subscribers:   make(map[types.EventType]types.Subscriber),
		quitChan:      make(chan interface{}),
	}, nil
}

// Start star block syncer
func (syncer *BlockSyncer) Start() error {
	go syncer.reqHandler()
	go syncer.recvHandler()

	syncer.subscribers[types.EventBlockCommitFailed] = syncer.eventCenter.Subscribe(types.EventBlockCommitFailed, syncer.GatherNewBlockFunc)
	syncer.subscribers[types.EventBlockVerifyFailed] = syncer.eventCenter.Subscribe(types.EventBlockVerifyFailed, syncer.GatherNewBlockFunc)
	syncer.subscribers[types.EventBlockCommitted] = syncer.eventCenter.Subscribe(types.EventBlockCommitted, syncer.GatherNewBlockFunc)
	syncer.subscribers[types.EventAddPeer] = syncer.eventCenter.Subscribe(types.EventAddPeer, syncer.GatherNewBlockFunc)
	return nil
}

// Stop stop block syncer
func (syncer *BlockSyncer) Stop() {
	syncer.lock.RLock()
	defer syncer.lock.RUnlock()
	close(syncer.quitChan)
	for eventType, subscriber := range syncer.subscribers {
		delete(syncer.subscribers, eventType)
		syncer.eventCenter.UnSubscribe(eventType, subscriber)
	}
}

// GatherNewBlockFunc gather new block from p2p network.
func (syncer *BlockSyncer) GatherNewBlockFunc(msg interface{}) {
	syncer.blockSyncChan <- msg
}

// send block sync request to gather the newest block from p2p
func (syncer *BlockSyncer) reqHandler() {
	timer := time.NewTicker(60 * time.Second)
	for {
		currentBlock := syncer.blockChain.GetCurrentBlock()
		hashStop := common.HeaderHash(currentBlock.Header)
		log.Debug("current block is %x, gather next block from p2p", hashStop)
		syncer.p2p.Gather(func(peerState uint64) bool {
			//TODO choose all peer as the candidate, so we can gather block more efficiently
			return true
		}, &message.BlockHeaderReq{
			Len:      1,
			HashStop: hashStop,
		})
		select {
		case <-syncer.blockSyncChan:
			continue
		case <-timer.C:
			continue
		case <-syncer.quitChan:
			return
		}
	}
}

// handle the block relative message from p2p network.
func (syncer *BlockSyncer) recvHandler() {
	msgChan := syncer.p2p.MessageChan()
	for {
		select {
		case msg := <-msgChan:
			switch msg.Payload.(type) {
			case *message.Block:
				bmsg := msg.Payload.(*message.Block)
				syncer.sendChan <- bmsg.Block
			case *message.BlockReq:
				brmsg := msg.Payload.(*message.BlockReq)
				block, err := syncer.blockChain.GetBlockByHash(brmsg.HeaderHash)
				if err != nil {
					return
				}
				bmsg := &message.Block{
					Block: block,
				}
				err = syncer.p2p.SendMsg(msg.From, bmsg)
				if err != nil {
					log.Error("failed to send message to peer %s, as: %v", msg.From.ToString(), err)
				}
			case *message.BlockHeaders:
				currentBlock := syncer.blockChain.GetCurrentBlock()
				bhmsg := msg.Payload.(*message.BlockHeaders)
				for i := 0; i < len(bhmsg.Headers); i++ {
					if bhmsg.Headers[i].Height <= currentBlock.Header.Height {
						continue
					}
					brmsg := &message.BlockReq{
						HeaderHash: common.HeaderHash(bhmsg.Headers[i]),
					}
					err := syncer.p2p.SendMsg(msg.From, brmsg)
					if err != nil {
						log.Error("failed to send message to peer %s, as: %v", msg.From.ToString(), err)
						return
					}
				}
			case *message.BlockHeaderReq:
				brmsg := msg.Payload.(*message.BlockHeaderReq)
				if brmsg.Len <= 0 {
					brmsg.Len = message.MAX_BLOCK_HEADER_NUM
				}
				blockHeaders := make([]*types.Header, 0)
				blockStop, err := syncer.blockChain.GetBlockByHash(brmsg.HashStop)
				if err != nil {
					log.Warn("have no block with Hash %x in local database", brmsg.HashStop)
				} else {
					for i := 1; i <= int(brmsg.Len); i++ {
						block, err := syncer.blockChain.GetBlockByHeight(blockStop.Header.Height + uint64(i))
						if err != nil {
							log.Warn("failed to get block with height %d, as:%v", blockStop.Header.Height+uint64(i), err)
							break
						}
						blockHeaders = append(blockHeaders, block.Header)
					}
				}
				bMsg := &message.BlockHeaders{
					Headers: blockHeaders,
				}
				err = syncer.p2p.SendMsg(msg.From, bMsg)
				if err != nil {
					log.Error("failed to send message to peer %s, as: %v", msg.From.ToString(), err)
				}
			}
		case <-syncer.quitChan:
			return
		}
	}
}
