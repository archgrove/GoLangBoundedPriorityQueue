package bpq

import (
  "fmt"
  "math"
)

type entry struct {
  item interface{}
  priority int
  inUse bool
}

type bpqBoundedHeapImpl struct {
 entries []entry
 nextSlot int
 highestPriorityIndex int
}

func (bpq bpqBoundedHeapImpl) String() string {
  return fmt.Sprintf("BPQ entries: %v", bpq.entries)
}

func (e entry) String() string {
  return fmt.Sprintf("{ Value %v, priority %v }", e.item, e.priority)
}

func makeBoundedHeap(capacity int) *bpqBoundedHeapImpl {
  result := bpqBoundedHeapImpl{ make([]entry, capacity), 0, 0 }

  for i := 0; i < capacity; i++ {
    result.entries[i] = entry{ nil, 0, false }
  }

  return &result
}

func (bpq *bpqBoundedHeapImpl) Capacity() int {
  return len(bpq.entries)
}

func (bpq *bpqBoundedHeapImpl) Push(item interface{}, priority int) bool {
  if (bpq.Capacity() == bpq.nextSlot) {
    // We're full!
    if (priority < bpq.entries[bpq.highestPriorityIndex].priority) {
      // But we can knock the back entry off
      bpq.entries[bpq.highestPriorityIndex].item = item
      bpq.entries[bpq.highestPriorityIndex].priority = priority
      bpq.entries[bpq.highestPriorityIndex].inUse = true

      bpq.bubbleUpIndex(bpq.highestPriorityIndex)

      // We must now restore the highest priority index. Alas, this is an O(n)
      // operation, but it only triggers when insertion is successful
      // and the heap is full
      highestPriority := math.MinInt64
      highestIndex := 0
      for i, v := range(bpq.entries) {
        if v.priority >= highestPriority {
          highestIndex = i
          highestPriority = v.priority
        }
      }
      bpq.highestPriorityIndex = highestIndex
      
      return true
    } else {
      return false
    }
  } else {
    // We're not full; add in the result
    bpq.entries[bpq.nextSlot].item = item
    bpq.entries[bpq.nextSlot].priority = priority
    bpq.entries[bpq.nextSlot].inUse = true

    // If we're highest than the highest, we become the highest
    // We've guarnteed to not bubble up in such cases
    if priority > bpq.entries[bpq.highestPriorityIndex].priority {
      bpq.highestPriorityIndex = bpq.nextSlot;
    } else {
      bpq.bubbleUpIndex(bpq.nextSlot)
    }

    bpq.nextSlot = bpq.nextSlot + 1

    return true
  }
}

func (bpq *bpqBoundedHeapImpl) Pop() (interface{}, error) {
  // defer fmt.Printf("Post-pop %v\n", bpq)

  if bpq.nextSlot == 0 {
    return nil, NoElementsError
  }

  result := bpq.entries[0];

  // Are we moving the highest priority index
  poppingHighest := bpq.nextSlot - 1 == bpq.highestPriorityIndex

  bpq.entries[0] = bpq.entries[bpq.nextSlot - 1]
  bpq.entries[bpq.nextSlot - 1].inUse = false
  bpq.nextSlot = bpq.nextSlot - 1

  if bpq.entries[0].inUse {
    newIndex := bpq.bubbleDownIndex(0)

    if poppingHighest {
      bpq.highestPriorityIndex = newIndex
    }
  }

  return result.item, nil
}

func (bpq *bpqBoundedHeapImpl) bubbleDownIndex(index int) int {
  // If we're greater than either of our children, swap with the child
  // and bubble down from their
  lChild, rChild := childrenOfIndex(index)

  if lChild < bpq.nextSlot && rChild < bpq.nextSlot {
     if bpq.entries[lChild].priority < bpq.entries[rChild].priority {
       if bpq.entries[lChild].priority < bpq.entries[index].priority {
         bpq.swapIndices(lChild, index)
         return bpq.bubbleDownIndex(lChild)
       }
     } else {
       if bpq.entries[rChild].priority < bpq.entries[index].priority {
         bpq.swapIndices(rChild, index)
         return bpq.bubbleDownIndex(rChild)
       }
     }
  } else if lChild < bpq.nextSlot && bpq.entries[lChild].priority < bpq.entries[index].priority {
     bpq.swapIndices(lChild, index)
     return bpq.bubbleDownIndex(lChild)
  } else if rChild < bpq.nextSlot && bpq.entries[rChild].priority < bpq.entries[index].priority {
     bpq.swapIndices(rChild, index)
     return bpq.bubbleDownIndex(rChild)
  }

  return index
}

func (bpq *bpqBoundedHeapImpl) bubbleUpIndex(index int) {
  parent := parentOfIndex(index)

  if bpq.entries[parent].priority > bpq.entries[index].priority {
    bpq.swapIndices(index, parent)
    bpq.bubbleUpIndex(parent)
  }
}

func (bpq *bpqBoundedHeapImpl) swapIndices(left, right int) {
   bpq.entries[left], bpq.entries[right] = bpq.entries[right], bpq.entries[left]
}

func parentOfIndex(index int) int {
  return (index - 1) / 2
}

func childrenOfIndex(index int) (int, int) {
  return (index * 2) + 1, (index * 2) + 2
}
