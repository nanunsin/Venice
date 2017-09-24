package venice

import (
	"container/ring"
	"fmt"
	"sync"
)

type LimitList struct {
	count    int
	lock     *sync.Mutex
	items    []interface{}
	fix, cap int
}

func NewLimitList(size int) *LimitList {
	instance := &LimitList{}

	instance.count = 0
	instance.lock = &sync.Mutex{}
	instance.items = make([]interface{}, size)
	instance.fix = size/2 + 1
	instance.cap = size

	return instance
}

func (lst *LimitList) Len() int {
	lst.lock.Lock()
	defer lst.lock.Unlock()

	return lst.count
}

func (lst *LimitList) IsEmpty() bool {
	return 0 == lst.Len()
}

func (lst *LimitList) Fix(fix int) {
	lst.lock.Lock()
	defer lst.lock.Unlock()

	lst.fix = fix
}

func (lst *LimitList) Push(item interface{}) {
	lst.lock.Lock()
	defer lst.lock.Unlock()

	lst.resize()

	lst.items[lst.count] = item
	lst.count++
}

func (lst *LimitList) ToSlice() []interface{} {
	lst.lock.Lock()
	defer lst.lock.Unlock()

	result := make([]interface{}, lst.count)
	copy(result, lst.items)

	return result
}

func (lst *LimitList) resize() {

	if lst.count >= lst.cap {

		temp := make([]interface{}, lst.cap)
		copy(temp, lst.items[(lst.cap-lst.fix)+1:])

		lst.items = temp
		lst.count = lst.fix - 1
	}
}

type RQueue struct {
	queue     *ring.Ring
	size, cnt int
}

func NewRQueue(size int) *RQueue {
	instance := &RQueue{
		queue: ring.New(size),
		size:  size,
		cnt:   0,
	}

	return instance
}

func (rq *RQueue) AddInfo(data float64) {
	rq.queue.Value = data
	rq.queue = rq.queue.Next()
	if rq.size > rq.cnt {
		rq.cnt++
	}
}

func (rq *RQueue) Len() int {
	return rq.cnt
}

func (rq *RQueue) Avg() (float64, bool) {
	if rq.size != rq.cnt {
		return 0.0, false
	}

	sum := 0.0
	rq.queue.Do(func(x interface{}) {
		sum += x.(float64)
	})

	return sum / (float64)(rq.size), true
}

type BufferList struct {
	items []interface{}
	cnt   int
}

func NewBufferList(size int) *BufferList {
	instance := &BufferList{
		cnt:   0,
		items: make([]interface{}, size),
	}
	return instance
}

func (bl *BufferList) Add(data interface{}) {
	if len(bl.items) == bl.cnt {
		bl.Clear()
	}
	bl.items[bl.cnt] = data
	bl.cnt++
}

func (bl *BufferList) Clear() {
	bl.cnt = 0
}

func (bl *BufferList) Print() {
	for i := 0; i < bl.cnt; i++ {
		fmt.Printf("%v", bl.items[i])
	}
	fmt.Println()
}

func (bl *BufferList) AverageFloat() float64 {
	if bl.cnt == 0 {
		return 0.0
	}
	sum := 0.0
	for i := 0; i < bl.cnt; i++ {
		sum += bl.items[i].(float64)
	}

	return (float64)(sum / float64(bl.cnt))
}

func (bl *BufferList) IsFull() bool {
	return bl.cnt == len(bl.items)
}
