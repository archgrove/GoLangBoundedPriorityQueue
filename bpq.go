package bpq

import (
  "errors"
)

type BPQ interface {
  Capacity() int
  Push(interface{}, int) bool
  Pop() (interface{}, error)
}

var NoElementsError error = errors.New("NoElementsError")

const maxRingBufferSize = 128

func BPQWithCapacity(capacity int) BPQ {
  if capacity <= maxRingBufferSize {
    return makeRingBuffer(capacity)
  } else {
    return makeBoundedHeap(capacity)
  }
}
