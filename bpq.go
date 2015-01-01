package bpq

import (
  "errors"
)

// BPQ represents a bounded min-priority queue
type BPQ interface {
  // Capacity returns the capacity (bounds) on the underlying queue
  Capacity() int
  // Push attempts to enqueue an item to the queue, with the given priority
  // Returns whether it succeded or not (failure is only due to a full queue,
  // where all items are lower priority)
  Push(interface{}, int) bool
  // Pop attempts to pop an item from the queue. It returns the popped item,
  // or an error. The only error it returns is NoElementsError, indicating
  // that the queue is empty
  Pop() (interface{}, error)
}

// NoElementsError indicates that a Pop was attempted on an empty queue
var NoElementsError error = errors.New("NoElementsError")

const maxRingBufferSize = 128

// BPQWithCapacity creates a new bounded priority queue with the given capacity
func BPQWithCapacity(capacity int) BPQ {
  if capacity <= maxRingBufferSize {
    return makeRingBuffer(capacity)
  } else {
    return makeBoundedHeap(capacity)
  }
}
