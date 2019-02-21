package message

import "github.com/DSiSc/craft/types"

// BlockReq block request message
type BlockReq struct {
	HeaderHash types.Hash `json:"header_hash"`
}

func (this *BlockReq) MsgId() types.Hash {
	return EmptyHash
}

func (this *BlockReq) MsgType() MessageType {
	return GET_BLOCK_TYPE
}

func (this *BlockReq) ResponseMsgType() MessageType {
	return BLOCK_TYPE
}

// Block block message
type Block struct {
	Block *types.Block `json:"block"`
}

func (this *Block) MsgId() types.Hash {
	return this.Block.HeaderHash
}

func (this *Block) MsgType() MessageType {
	return BLOCK_TYPE
}

func (this *Block) ResponseMsgType() MessageType {
	return NIL
}
