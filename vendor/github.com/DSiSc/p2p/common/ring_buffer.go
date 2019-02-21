package common

import (
	"container/list"
	"github.com/DSiSc/craft/types"
	"sync"
)

// RingBuffer is a ring buffer implementation.
type RingBuffer struct {
	elements map[types.Hash]interface{}
	limit    int
	keyList  *list.List
	lock     sync.RWMutex
}

// NewRingBuffer create a ring buffer instance
func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		elements: make(map[types.Hash]interface{}),
		limit:    size,
		keyList:  list.New(),
	}
}

// AddElement add a element to ring buffer
func (ring *RingBuffer) AddElement(hash types.Hash, elem interface{}) {
	ring.lock.Lock()
	defer ring.lock.Unlock()
	ring.keyList.PushFront(hash)
	ring.elements[hash] = elem
	if len(ring.elements) > ring.limit {
		node := ring.keyList.Back()
		delete(ring.elements, node.Value.(types.Hash))
	}
}

// Exist check if the element is already in the ring buffer
func (ring *RingBuffer) Exist(hash types.Hash) bool {
	ring.lock.RLock()
	defer ring.lock.RUnlock()
	return ring.elements[hash] != nil
}
