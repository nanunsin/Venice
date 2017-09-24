package venice

import (
	"fmt"
	"testing"
)

func TestLimitList_ToSlice(t *testing.T) {
	t.Log("TestLimitList_ToSlice Start")

	list := NewLimitList(20)
	list.Fix(10)
	for i := 0; i < 29; i++ {
		list.Push(i)
	}

	b := list.ToSlice()
	fmt.Println(b)
}

func TestBufferList(t *testing.T) {
	t.Log("TestBufferList Start")

	rq := NewBufferList(4)

	for i := 0; i < 10; i++ {
		rq.Add(float64(i))
		rq.Print()
		fmt.Println(rq.AverageFloat())
	}
}
