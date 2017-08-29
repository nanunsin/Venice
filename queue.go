package main

import "sync"

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
