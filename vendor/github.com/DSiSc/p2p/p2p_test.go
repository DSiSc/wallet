package p2p

import (
	"errors"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/monkey"
	"github.com/DSiSc/p2p/common"
	"github.com/DSiSc/p2p/config"
	"github.com/DSiSc/p2p/message"
	"github.com/stretchr/testify/assert"
	"net"
	"reflect"
	"sync"
	"testing"
	"time"
)

func mockConfig() *config.P2PConfig {
	return &config.P2PConfig{
		AddrBookFilePath: "",
		ListenAddress:    "tcp://0.0.0.0:8080",
		MaxConnOutBound:  60,
		MaxConnInBound:   20,
		PersistentPeers:  "",
	}
}

func mockPeer(serverAddr, addr *common.NetAddress, outBound, persistent bool, msgChan chan<- *InternalMsg, conn net.Conn) *Peer {
	peer := newPeer(serverAddr, addr, outBound, persistent, msgChan, conn)
	monkey.PatchInstanceMethod(reflect.TypeOf(peer), "Start", func(peer *Peer) error {
		return nil
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(peer), "Stop", func(peer *Peer) {
	})
	return peer
}

func TestNewP2P(t *testing.T) {
	assert := assert.New(t)
	conf := mockConfig()
	p2p, err := NewP2P(conf, &eventCenter{})
	assert.Nil(err)
	assert.NotNil(p2p)
}

func TestP2P_Start(t *testing.T) {
	defer monkey.UnpatchAll()
	assert := assert.New(t)
	conf := mockConfig()
	p2p, err := NewP2P(conf, &eventCenter{})
	assert.Nil(err)

	// mock listen
	monkey.Patch(net.Listen, func(network, address string) (net.Listener, error) {
		return newTestListener(), nil
	})
	err = p2p.Start()
	assert.Nil(err)
	p2p.Stop()
}

func TestP2P_Stop(t *testing.T) {
	defer monkey.UnpatchAll()
	assert := assert.New(t)
	conf := mockConfig()
	p2p, err := NewP2P(conf, &eventCenter{})
	assert.Nil(err)

	// mock listen
	monkey.Patch(net.Listen, func(network, address string) (net.Listener, error) {
		return newTestListener(), nil
	})
	err = p2p.Start()
	assert.Nil(err)
	p2p.Stop()
	select {
	case <-p2p.quitChan:
	default:
		assert.Error(errors.New("failed To stop the peer."))
	}
}

func TestP2P_BroadCast(t *testing.T) {
	defer monkey.UnpatchAll()
	assert := assert.New(t)
	conf := mockConfig()
	conf.PersistentPeers = "tcp://192.168.1.1:8080"
	p2p, err := NewP2P(conf, &eventCenter{})
	assert.Nil(err)
	msg := &message.PingMsg{
		State: 1,
	}

	//mock peer
	serverAddr, _ := common.ParseNetAddress(conf.ListenAddress)
	addr, _ := common.ParseNetAddress(conf.PersistentPeers)
	peer := mockPeer(serverAddr, addr, true, false, p2p.internalChan, nil)
	monkey.Patch(NewOutboundPeer, func(serverAddr, addr *common.NetAddress, persistent bool, msgChan chan<- *InternalMsg) *Peer {
		return peer
	})

	err = p2p.Start()
	assert.Nil(err)

	timer := time.NewTicker(time.Second)
OUT:
	for {
		select {
		case <-timer.C:
			if len(p2p.GetPeers()) > 0 {
				break OUT
			}
		}
	}
	p2p.BroadCast(msg)
	// read message From peer's send channel
	timeoutTricker := time.NewTicker(5 * time.Second)
	var wg sync.WaitGroup
	for _, peer := range p2p.GetPeers() {
		wg.Add(1)
		go func(p *Peer) {
			for {
				select {
				case pmsg := <-p.sendChan:
					switch pmsg.Payload.(type) {
					case *message.PingMsg:
						assert.Equal(msg, pmsg.Payload)
						wg.Done()
						return
					}
				case <-timeoutTricker.C:
					assert.Nil(errors.New("read sent message failed"))
				}
			}
		}(peer)
	}
	wg.Wait()
	peer.Stop()
}

func TestP2P_SendMsg(t *testing.T) {
	defer monkey.UnpatchAll()
	assert := assert.New(t)
	conf := mockConfig()
	conf.PersistentPeers = "tcp://192.168.1.1:8080"
	p2p, err := NewP2P(conf, &eventCenter{})
	assert.Nil(err)
	//mock peer
	serverAddr, _ := common.ParseNetAddress(conf.ListenAddress)
	addr, _ := common.ParseNetAddress(conf.PersistentPeers)
	mockPeer := mockPeer(serverAddr, addr, true, false, p2p.internalChan, nil)
	monkey.Patch(NewOutboundPeer, func(serverAddr, addr *common.NetAddress, persistent bool, msgChan chan<- *InternalMsg) *Peer {
		return mockPeer
	})

	// mock listen
	monkey.Patch(net.Listen, func(network, address string) (net.Listener, error) {
		return newTestListener(), nil
	})
	err = p2p.Start()
	assert.Nil(err)

	timeoutTricker := time.NewTicker(5 * time.Second)
	timer := time.NewTicker(time.Second)
OUT:
	for {
		select {
		case <-timer.C:
			if len(p2p.GetPeers()) > 0 {
				break OUT
			}
		case <-timeoutTricker.C:
			assert.Nil(errors.New("failed To connect persistent peer"))
			break OUT
		}
	}
	msg := &message.BlockReq{}
	peer := p2p.GetPeers()[0]
	go func() {
		err := p2p.SendMsg(peer.addr, msg)
		assert.Nil(err)
	}()
	// read message From peer's send channel
OUT1:
	for {
		select {
		case pmsg := <-peer.sendChan:
			switch pmsg.Payload.(type) {
			case *message.BlockReq:
				break OUT1
			default:
				continue
			}
		case <-timeoutTricker.C:
			assert.Nil(errors.New("read sent message failed"))
		}
	}
	p2p.Stop()
}

func TestP2P_GetOutBountPeersCount(t *testing.T) {
	defer monkey.UnpatchAll()
	assert := assert.New(t)
	conf := mockConfig()
	conf.PersistentPeers = "tcp://192.168.1.1:8080"

	p2p, err := NewP2P(conf, &eventCenter{})
	assert.Nil(err)
	assert.Equal(0, p2p.GetOutBountPeersCount())

	//mock peer
	serverAddr, _ := common.ParseNetAddress(conf.ListenAddress)
	addr, _ := common.ParseNetAddress(conf.PersistentPeers)
	peer := mockPeer(serverAddr, addr, true, false, p2p.internalChan, nil)
	monkey.Patch(NewOutboundPeer, func(serverAddr, addr *common.NetAddress, persistent bool, msgChan chan<- *InternalMsg) *Peer {
		return peer
	})

	// mock listen
	monkey.Patch(net.Listen, func(network, address string) (net.Listener, error) {
		return newTestListener(), nil
	})
	err = p2p.Start()
	assert.Nil(err)
	timer := time.NewTicker(time.Second)
OUT:
	for {
		select {
		case <-timer.C:
			if len(p2p.GetPeers()) > 0 {
				break OUT
			}
		}
	}
	assert.Equal(1, p2p.GetOutBountPeersCount())
	p2p.Stop()
}

func TestP2P_GetPeerByAddress(t *testing.T) {
	defer monkey.UnpatchAll()
	assert := assert.New(t)
	conf := mockConfig()
	conf.PersistentPeers = "tcp://192.168.1.1:8080"
	p2p, err := NewP2P(conf, &eventCenter{})
	assert.Nil(err)

	//mock peer
	serverAddr, _ := common.ParseNetAddress(conf.ListenAddress)
	addr, _ := common.ParseNetAddress(conf.PersistentPeers)
	peer := mockPeer(serverAddr, addr, true, false, p2p.internalChan, nil)
	monkey.Patch(NewOutboundPeer, func(serverAddr, addr *common.NetAddress, persistent bool, msgChan chan<- *InternalMsg) *Peer {
		return peer
	})

	// mock listen
	monkey.Patch(net.Listen, func(network, address string) (net.Listener, error) {
		return newTestListener(), nil
	})

	err = p2p.Start()
	assert.Nil(err)
	timer := time.NewTicker(time.Second)
OUT:
	for {
		select {
		case <-timer.C:
			if len(p2p.GetPeers()) > 0 {
				break OUT
			}
		}
	}
	p2p.Stop()
}

func TestP2P_GetPeers(t *testing.T) {
	defer monkey.UnpatchAll()
	assert := assert.New(t)
	conf := mockConfig()
	conf.PersistentPeers = "tcp://192.168.1.1:8080"
	p2p, err := NewP2P(conf, &eventCenter{})
	// mock peer
	serverAddr, _ := common.ParseNetAddress(conf.ListenAddress)
	addr, _ := common.ParseNetAddress(conf.PersistentPeers)
	peer := mockPeer(serverAddr, addr, true, false, p2p.internalChan, nil)
	monkey.Patch(NewInboundPeer, func(serverAddr, addr *common.NetAddress, msgChan chan<- *InternalMsg, conn net.Conn) *Peer {
		return peer
	})
	monkey.Patch(NewOutboundPeer, func(serverAddr, addr *common.NetAddress, persistent bool, msgChan chan<- *InternalMsg) *Peer {
		return peer
	})

	// mock listen
	monkey.Patch(net.Listen, func(network, address string) (net.Listener, error) {
		return newTestListener(), nil
	})

	assert.Nil(err)
	err = p2p.Start()
	assert.Nil(err)
	timer := time.NewTicker(time.Second)
OUT:
	for {
		select {
		case <-timer.C:
			if len(p2p.GetPeers()) > 0 {
				break OUT
			}
		}
	}
	assert.Equal(1, len(p2p.GetPeers()))
}

func TestP2P_Gather(t *testing.T) {
	defer monkey.UnpatchAll()
	assert := assert.New(t)
	conf := mockConfig()
	conf.PersistentPeers = "tcp://192.168.1.1:8080"
	p2p, err := NewP2P(conf, &eventCenter{})
	assert.Nil(err)

	//mock peer
	serverAddr, _ := common.ParseNetAddress(conf.ListenAddress)
	addr, _ := common.ParseNetAddress(conf.PersistentPeers)
	mockPeer := mockPeer(serverAddr, addr, true, false, p2p.internalChan, nil)
	monkey.Patch(NewOutboundPeer, func(serverAddr, addr *common.NetAddress, persistent bool, msgChan chan<- *InternalMsg) *Peer {
		return mockPeer
	})

	// mock listen
	monkey.Patch(net.Listen, func(network, address string) (net.Listener, error) {
		return newTestListener(), nil
	})
	err = p2p.Start()
	assert.Nil(err)

	time.Sleep(time.Second)
	if len(p2p.GetPeers()) <= 0 {
		assert.Nil(errors.New("failed To connect persistent peer"))
	}

	// retrieve message From send channel
	go func() {
		for {
			select {
			case msg := <-p2p.GetPeers()[0].sendChan:
				switch msg.Payload.MsgType() {
				case message.GET_BLOCK_TYPE:
					p2p.internalChan <- &InternalMsg{
						From:    mockPeer.GetAddr(),
						Payload: &message.Block{},
					}
				}
			}
		}
	}()
	p2p.Gather(func(peerState uint64) bool {
		return true
	}, &message.BlockReq{})
	timer := time.NewTicker(time.Second)
	select {
	case msg := <-p2p.MessageChan():
		if msg.Payload.MsgType() != message.BLOCK_TYPE {
			assert.Nil(errors.New("failed To gather block From p2p"))
		}
	case <-timer.C:
		assert.Nil(errors.New("failed To connect persistent peer"))
	}
	p2p.Stop()
}

type testListener struct {
	connChan chan net.Conn
}

func newTestListener() *testListener {
	return &testListener{
		connChan: make(chan net.Conn),
	}
}

func (this *testListener) Accept() (conn net.Conn, err error) {
	defer func() {
		// recover From panic if one occured.
		if recover() != nil {
			err = errors.New("listener have stopped")
		}
	}()
	conn = <-this.connChan
	return
}

func (this *testListener) Close() error {
	close(this.connChan)
	return nil
}

func (this *testListener) Addr() net.Addr {
	return nil
}

type eventCenter struct {
}

// subscriber subscribe specified eventType with eventFunc
func (*eventCenter) Subscribe(eventType types.EventType, eventFunc types.EventFunc) types.Subscriber {
	return nil
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
