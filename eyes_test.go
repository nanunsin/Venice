package venice

import (
	"fmt"
	"testing"
)

func Test_HawkEye(t *testing.T) {
	t.Log("TestLimitList_ToSlice Start")

	chPrice := make(chan int64, 10)
	eyes := NewHawkEye(chPrice, "ETH")
	go eyes.Scout()

	for i := 0; i < 2; i++ {
		fmt.Println(<-chPrice)
	}
	eyes.Stop()
}
