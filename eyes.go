package venice

import (
	"time"

	"github.com/nanunsin/bithumbb/bithumb"
)

type HawkEye struct {
	bit      *bithumb.Bithumb
	chStop   chan bool
	chPrice  chan<- float64
	cointype int
}

func NewHawkEye(cb chan<- float64, cointype int) *HawkEye {
	instance := &HawkEye{
		bit:      bithumb.NewBithumb("test", "sec"),
		chStop:   make(chan bool),
		chPrice:  cb,
		cointype: cointype,
	}
	return instance
}

func (h *HawkEye) Scout() {
	bContinue := true

	h.sendPrice()

	ticker := time.NewTicker(time.Second * 10)
	for bContinue {
		select {
		case <-ticker.C:
			h.sendPrice()
		case <-h.chStop:
			bContinue = false
		}
	}
}

func (h *HawkEye) sendPrice() {
	var info bithumb.WMP
	h.bit.GetPrice(h.cointype, &info)
	h.chPrice <- info.Price
}

func (h *HawkEye) Stop() {
	h.chStop <- true
}
