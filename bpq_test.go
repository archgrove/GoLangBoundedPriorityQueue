package bpq

import (
  "testing"
  "fmt"
  "math/rand"
  "sort"
)

func TestCapacity(t *testing.T) {
  queue1 := BPQWithCapacity(10)
  queue2 := BPQWithCapacity(20)
  queue3 := BPQWithCapacity(100)

  if queue1.Capacity() != 10 {
    t.Error("Capacity set to 10, was not 10")
  }

  if queue2.Capacity() != 20 {
    t.Error("Capacity set to 20, was not 20")
  }

  if queue3.Capacity() != 100 {
    t.Error("Capacity set to 100, was not 100")
  }
}

// Tests for empty pops

func TestEmptyPopWithRingBuffer(t *testing.T) {
  testEmptyPopWithBPQ(BPQWithCapacity(maxRingBufferSize), t)
}

func TestEmptyPopWithBoundedHeap(t *testing.T) {
  testEmptyPopWithBPQ(BPQWithCapacity(maxRingBufferSize + 1), t)
}

func testEmptyPopWithBPQ(queue BPQ, t *testing.T) {
  v, err := queue.Pop()

  if v != nil {
    t.Error("Expected an nil value when popping an empty queue")
  }

  if err == nil {
    t.Error("Expected an error value when popping an empty queue")
  }
}

// Test for simple push/pop

func TestSimplePushPopWithRingBuffer(t *testing.T) {
  testSimplePushPopWithQueue(BPQWithCapacity(maxRingBufferSize), t)
}

func TestSimplePushPopWithBoundedHeap(t *testing.T) {
  testSimplePushPopWithQueue(BPQWithCapacity(maxRingBufferSize + 1), t)
}

func testSimplePushPopWithQueue(queue BPQ, t *testing.T) {
  queue.Push(1, 10)
  v, err := queue.Pop()

  if err != nil {
    t.Error("Unexpected error value when popping a non-empty queue")
  }

  if v == nil {
    t.Error("Recieved an unexpected empty value")
  }

  c, suc := v.(int)

  if !suc || c != 1 {
    t.Error("Expected integer value 1 in reply")
  }
}

// Test for push/pop/pop

func TestPushDoublePopWithRingBuffer(t *testing.T) {
  testPushDoublePopWithQueue(BPQWithCapacity(maxRingBufferSize), t)
}

func TestPushDoublePopWithBoundedHeap(t *testing.T) {
  testPushDoublePopWithQueue(BPQWithCapacity(maxRingBufferSize + 1), t)
}

func testPushDoublePopWithQueue(queue BPQ, t *testing.T) {
  queue.Push(1, 10)
  queue.Pop()
  v, err := queue.Pop()
  if err == nil {
    t.Error("Expected error value when popping an empty queue")
  }

  if v != nil {
    t.Error("Did not expect a value when popping an empty queue")
  }
}

// Test for priority ordering

func TestPriorityOrderingWithRingBuffer(t *testing.T) {
  testPriorityOrderingWithQueue(BPQWithCapacity(maxRingBufferSize), t)
}

func TestPriorityOrderingWithBoundedHeap(t *testing.T) {
  testPriorityOrderingWithQueue(BPQWithCapacity(maxRingBufferSize + 1), t)
}

func testPriorityOrderingWithQueue(queue BPQ, t *testing.T) {
  queue.Push(1, 10)
  queue.Push(2, 5)
  queue.Push(3, 100)
  v1, _ := queue.Pop()
  v2, _ := queue.Pop()
  v3, _ := queue.Pop()


  if (v1 == nil || v1.(int) != 2) ||
     (v2 == nil || v2.(int) != 1) ||
     (v3 == nil || v3.(int) != 3) {
    result := fmt.Sprintf("%v, %v, %v\n", v1, v2, v3)
    t.Error("Unexpected priority ordering, got " + result)
  }
}


// Test for over-fill
// Overfill tests currently only test the ring buffer

func TestOverFill(t *testing.T) {
  queue := BPQWithCapacity(5)

  queue.Push(1, 10)
  queue.Push(2, 11)
  queue.Push(3, 20)
  queue.Push(4, 30)
  queue.Push(5, 40)
  queue.Push(6, 50)
  queue.Push(7, 5)
  queue.Push(8, 1)
  queue.Push(9, 50)
  queue.Push(10, 35)

  v1, _ := queue.Pop()
  v2, _ := queue.Pop()
  v3, _ := queue.Pop()
  v4, _ := queue.Pop()
  v5, _ := queue.Pop()
  v6, err := queue.Pop()

  // We should not get a v6
  if v6 != nil || err == nil {
    t.Error("Got more than 5 results!")
  }

  if (v1 == nil || v1.(int) != 8) ||
     (v2 == nil || v2.(int) != 7) ||
     (v3 == nil || v3.(int) != 1) ||
     (v4 == nil || v4.(int) != 2) ||
     (v4 == nil || v5.(int) != 3) {
    t.Errorf("Incorrect priority ordering, got %v, %v, %v, %v, %v",
             v1, v2, v3, v4, v5)
  }
}


//
// Psuedo-random tests (with fixed seed)
//

type Entry struct {
  value int
  priority int
}

type Entries []Entry

func (es Entries) Len() int {
  return len(es)
}

func (es Entries) Swap(i, j int) {
  es[i], es[j] = es[j], es[i]
}

func (es Entries) Less(i, j int) bool {
  return es[i].priority < es[j].priority
}

func (es Entries) ContainsPriority(priority int) bool {
  for _, v := range(es) {
    if v.priority == priority {
      return true
    }
  }

  return false
}

func TestRandomInsertAndPopWithRingBuffer(t *testing.T) {
  testRandomInsertAndPopWithQueue(BPQWithCapacity(maxRingBufferSize), 
                                  maxRingBufferSize, t)
}

func TestRandomInsertAndPopWithBoundedHeap(t *testing.T) {
  testRandomInsertAndPopWithQueue(BPQWithCapacity(100 * maxRingBufferSize),
                                  100 * maxRingBufferSize, t)
}

func testRandomInsertAndPopWithQueue(queue BPQ, max int, t *testing.T) {
  // Use a determined seed for reproducability
  rand.Seed(123456)

  // Generate a random list of entries and priorities
  
  randomItems := make(Entries, max)
  for i := 0; i < max; i++ {
    var item Entry = Entry{ rand.Int(), rand.Int() }

    // We need to ensure each item has a distinct priority; our priority sort is
    // not stable in any sense
    for {
      if randomItems.ContainsPriority(item.priority) {
        item = Entry{ rand.Int() % 100, rand.Int() % 100 }
      } else {
        break
      }
    }

    randomItems[i] = item
    queue.Push(item.value, item.priority)
  }

  sort.Sort(randomItems)

  for _, entry := range(randomItems) {
    v, _ := queue.Pop()

    if entry.value != v.(int) {
      t.Errorf("Not the right value, got %v but wanted %v", v.(int), 
               entry.value)
    }
  }
}

//
// Benchmarks
//

func BenchmarkBPQRingBuffer(b *testing.B) {
  benchmarkBPQ(BPQWithCapacity(maxRingBufferSize), b)
}

func BenchmarkBPQBoundedHeap(b *testing.B) {
  benchmarkBPQ(BPQWithCapacity(100 * maxRingBufferSize), b)
}

func benchmarkBPQ(queue BPQ, b *testing.B) {
  var max = queue.Capacity()

  for n := 0; n < b.N; n++ {
    for i := 0; i < max; i++ {
      queue.Push(i, i % 10)
    }

    for i := 0; i < max; i++ {
      queue.Pop()
    }
  }
}
