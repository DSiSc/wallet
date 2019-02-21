package message

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DSiSc/craft/types"
	"io"
)

var EmptyHash = types.Hash{}

type Message interface {
	MsgId() types.Hash
	MsgType() MessageType
	ResponseMsgType() MessageType
}

//MessageType is the message type in p2p network
type MessageType uint32

const (
	NIL              = MessageType(iota) // nil message type
	VERSION_TYPE                         //peer`s information
	VERACK_TYPE                          //ack msg after version recv
	GETADDR_TYPE                         //req nbr address from peer
	ADDR_TYPE                            //nbr address
	PING_TYPE                            //ping  sync height
	PONG_TYPE                            //pong  recv nbr height
	GET_HEADERS_TYPE                     //req blk hdr
	HEADERS_TYPE                         //blk hdr
	BLOCK_TYPE                           //blk payload
	TX_TYPE                              //transaction
	GET_BLOCK_TYPE                       //req blks from peer
	NOT_FOUND_TYPE                       //peer can`t find blk according to the hash
	REJECT_TYPE
	DISCONNECT_TYPE //peer disconnect info raise by link
	TRACE_TYPE      //trace message
)

// message's header
type messageHeader struct {
	Magic   uint32
	MsgType MessageType
	Length  uint32
}

// EncodeMessage encode message to byte array.
func EncodeMessage(msg Message) ([]byte, error) {
	if msg == nil {
		return nil, errors.New("empty message content")
	}

	msgByte, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode message %v to json, as: %v", msg, err)
	}

	header, err := buildMessageHeader(msg, len(msgByte))
	if err != nil {
		return nil, err
	}

	buf, err := encodeMessageHeader(header)
	if err != nil {
		return nil, err
	}

	return append(buf, msgByte...), nil
}

// encodeMessageHeader encode message header to byte array.
func encodeMessageHeader(header *messageHeader) ([]byte, error) {
	buf := make([]byte, 12)
	binary.LittleEndian.PutUint32(buf, header.Magic)
	binary.LittleEndian.PutUint32(buf[4:], uint32(header.MsgType))
	binary.LittleEndian.PutUint32(buf[8:], header.Length)
	return buf, nil
}

// ReadMessage read message
func ReadMessage(reader io.Reader) (Message, error) {
	header, err := readMessageHeader(reader)
	if err != nil {
		return nil, err
	}

	body := make([]byte, header.Length)
	_, err = io.ReadFull(reader, body)
	if err != nil {
		return nil, err
	}

	msg, err := makeEmptyMessage(header.MsgType)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// read message header from reader.
func readMessageHeader(reader io.Reader) (messageHeader, error) {
	msgh := messageHeader{}
	err := binary.Read(reader, binary.LittleEndian, &msgh)
	return msgh, err
}

// fill the header according to the message.
func buildMessageHeader(msg Message, len int) (*messageHeader, error) {
	header := &messageHeader{
		Magic:   0,
		MsgType: msg.MsgType(),
		Length:  uint32(len),
	}
	return header, nil
}

// make empty message according to the message type
func makeEmptyMessage(msgType MessageType) (Message, error) {
	switch msgType {
	case VERSION_TYPE:
		return &Version{}, nil
	case VERACK_TYPE:
		return &VersionAck{}, nil
	case PING_TYPE:
		return &PingMsg{}, nil
	case PONG_TYPE:
		return &PongMsg{}, nil
	case GETADDR_TYPE:
		return &AddrReq{}, nil
	case ADDR_TYPE:
		return &Addr{}, nil
	case REJECT_TYPE:
		return &RejectMsg{}, nil
	case GET_HEADERS_TYPE:
		return &BlockHeaderReq{}, nil
	case HEADERS_TYPE:
		return &BlockHeaders{}, nil
	case GET_BLOCK_TYPE:
		return &BlockReq{}, nil
	case BLOCK_TYPE:
		return &Block{}, nil
	case TX_TYPE:
		return &Transaction{}, nil
	case TRACE_TYPE:
		return &TraceMsg{}, nil
	default:
		return nil, fmt.Errorf("unknown message type %v", msgType)
	}
}
