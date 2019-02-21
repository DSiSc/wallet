package syncer

import (
	"errors"
	"github.com/DSiSc/blockchain"
	bcfg "github.com/DSiSc/blockchain/config"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/monkey"
	"github.com/DSiSc/p2p"
	pcfg "github.com/DSiSc/p2p/config"
	"github.com/DSiSc/p2p/message"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	cg := bcfg.BlockChainConfig{
		PluginName: blockchain.PLUGIN_MEMDB,
	}
	blockchain.InitBlockChain(cg, &eventCenter{})
	bc, _ := blockchain.NewLatestStateBlockChain()
	monkey.PatchInstanceMethod(reflect.TypeOf(bc), "GetCurrentBlock", func(this *blockchain.BlockChain) *types.Block {
		return &types.Block{
			Header: &types.Header{
				Height: 2,
			},
		}
	})
	monkey.Patch(blockchain.NewLatestStateBlockChain, func() (*blockchain.BlockChain, error) {
		return bc, nil
	})
	m.Run()
}

func mockP2PConfig() *pcfg.P2PConfig {
	return &pcfg.P2PConfig{}
}

func mockBlockMsg() *p2p.InternalMsg {
	block := &types.Block{
		Header: &types.Header{
			Height: 3}}
	bMsg := &message.Block{
		Block: block,
	}
	iMsg := &p2p.InternalMsg{
		Payload: bMsg,
	}
	return iMsg
}

func TestNewBlockSyncer(t *testing.T) {
	assert := assert.New(t)
	p, _ := p2p.NewP2P(mockP2PConfig(), &eventCenter{})
	sendChan := make(chan interface{})
	bs, err := NewBlockSyncer(p, sendChan, &eventCenter{})
	assert.Nil(err)
	assert.NotNil(bs)
}

func TestBlockSyncer_Start(t *testing.T) {
	defer monkey.UnpatchAll()
	assert := assert.New(t)

	pMsgChan := make(chan *p2p.InternalMsg)
	p, _ := p2p.NewP2P(mockP2PConfig(), &eventCenter{})
	monkey.PatchInstanceMethod(reflect.TypeOf(p), "Gather", func(this *p2p.P2P, filter p2p.PeerFilter, msg message.Message) error {
		pMsgChan <- mockBlockMsg()
		return nil
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(p), "MessageChan", func(*p2p.P2P) <-chan *p2p.InternalMsg {
		return pMsgChan
	})

	sendChan := make(chan interface{})

	evc := &eventCenter{}

	bs, err := NewBlockSyncer(p, sendChan, evc)
	assert.Nil(err)
	assert.NotNil(bs)

	err = bs.Start()
	assert.Nil(err)

	select {
	case msg := <-sendChan:
		switch msg.(type) {
		case *types.Block:
			bmsg := msg.(*types.Block)
			assert.Equal(uint64(3), bmsg.Header.Height)
		default:
			assert.Nil(errors.New("failed to gather block from p2p"))
		}
	}
	bs.Stop()
}

func TestBlockSyncer_Stop(t *testing.T) {
	defer monkey.UnpatchAll()
	assert := assert.New(t)
	p, _ := p2p.NewP2P(mockP2PConfig(), &eventCenter{})
	monkey.PatchInstanceMethod(reflect.TypeOf(p), "Gather", func(this *p2p.P2P, filter p2p.PeerFilter, msg message.Message) error {
		return nil
	})
	sendChan := make(chan interface{})
	evc := &eventCenter{}
	bs, err := NewBlockSyncer(p, sendChan, evc)
	assert.Nil(err)
	assert.NotNil(bs)

	err = bs.Start()
	assert.Nil(err)
	time.Sleep(time.Second)
	bs.Stop()
	select {
	case <-bs.quitChan:
	default:
		assert.Nil(errors.New("failed to stop syncer"))
	}
	bs.lock.Lock()
	assert.Equal(0, len(bs.subscribers))
	bs.lock.Unlock()
}

type eventCenter struct {
	Subscribers map[types.EventType]map[types.Subscriber]types.EventFunc
}

// subscriber subscribe specified eventType with eventFunc
func (*eventCenter) Subscribe(eventType types.EventType, eventFunc types.EventFunc) types.Subscriber {
	return make(chan interface{})
}

// subscriber unsubscribe specified eventType
func (*eventCenter) UnSubscribe(eventType types.EventType, subscriber types.Subscriber) (err error) {
	return nil
}

// notify subscriber of eventType
func (*eventCenter) Notify(eventType types.EventType, value interface{}) (err error) {
	return nil
}

// notify specified eventFunc
func (*eventCenter) NotifySubscriber(eventFunc types.EventFunc, value interface{}) {

}

// notify subscriber traversing all events
func (*eventCenter) NotifyAll() (errs []error) {
	return nil
}

// unsubscrible all event
func (*eventCenter) UnSubscribeAll() {
}
