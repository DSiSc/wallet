package p2p

import (
	"errors"
	"github.com/DSiSc/monkey"
	"github.com/DSiSc/p2p/message"
	"github.com/stretchr/testify/assert"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestNewPeerConn(t *testing.T) {
	assert := assert.New(t)
	conn := newTestConn()
	recvChan := make(chan message.Message)
	peerConn := NewPeerConn(conn, recvChan)
	assert.NotNil(peerConn)
}

func TestPeerConn_Start(t *testing.T) {
	defer monkey.UnpatchAll()
	assert := assert.New(t)
	conn := newTestConn()
	monkey.PatchInstanceMethod(reflect.TypeOf(conn), "Read", func(c *testConn, b []byte) (n int, err error) {
		msg := &message.PingMsg{
			State: 1,
		}
		msgByte, _ := message.EncodeMessage(msg)
		copy(b, msgByte)
		return len(msgByte), nil
	})
	recvChan := make(chan message.Message)
	peerConn := NewPeerConn(conn, recvChan)
	assert.NotNil(peerConn)
	peerConn.Start()

	msg := &message.PingMsg{
		State: 1,
	}
	timer := time.NewTicker(time.Second)
	select {
	case m := <-recvChan:
		assert.Equal(msg, m)
	case <-timer.C:
		assert.Nil(errors.New("read message From connection time out"))
	}
	peerConn.Stop()
}

func TestPeerConn_Stop(t *testing.T) {
	assert := assert.New(t)
	conn := newTestConn()
	recvChan := make(chan message.Message)
	peerConn := NewPeerConn(conn, recvChan)
	assert.NotNil(peerConn)
	peerConn.Start()
	peerConn.Stop()
	timer := time.NewTicker(time.Second)
	select {
	case <-peerConn.quitChan:
	case <-timer.C:
		assert.Nil(errors.New("failed To stop peer connection"))
	}
}

func TestPeerConn_SendMessage(t *testing.T) {
	defer monkey.UnpatchAll()
	assert := assert.New(t)
	testChan := make(chan []byte)
	conn := newTestConn()
	monkey.PatchInstanceMethod(reflect.TypeOf(conn), "Read", func(c *testConn, b []byte) (n int, err error) {
		msgByte, _ := message.EncodeMessage(&message.PingMsg{
			State: 1,
		})
		copy(b, msgByte)
		return len(msgByte), nil
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(conn), "Write", func(this *testConn, b []byte) (n int, err error) {
		go func(bs []byte) {
			testChan <- b
		}(b)
		return len(b), nil
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(conn), "RemoteAddr", func(this *testConn) net.Addr {
		_, addr, _ := net.ParseCIDR("192.168.1.1/24")
		return addr
	})
	recvChan := make(chan message.Message)
	peerConn := NewPeerConn(conn, recvChan)
	assert.NotNil(peerConn)
	peerConn.Start()
	msg := &message.PingMsg{
		State: 1,
	}
	msgByte1, _ := message.EncodeMessage(msg)

	err := peerConn.SendMessage(msg)
	assert.Nil(err)

	timer := time.NewTicker(time.Second)
	select {
	case msgByte2 := <-testChan:
		assert.Equal(msgByte1, msgByte2)
	case <-timer.C:
		assert.Nil(errors.New("failed To stop peer connection"))
	}

	peerConn.Stop()
}
