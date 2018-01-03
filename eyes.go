package venice

import (
	"time"

	"github.com/nanunsin/bithumbb/bithumb"
)

type HawkEye struct {
	bit      *bithumb.Bithumb
	chStop   chan bool
	chPrice  chan<- int64
	cointype string
}

func NewHawkEye(cb chan<- int64, cointype string) *HawkEye {
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
	var info bithumb.TickerInfo
	h.bit.GetPrice(h.cointype, &info)
	h.chPrice <- info.Price
}

func (h *HawkEye) Stop() {
	h.chStop <- true
}

type EagleEye struct {
	bit      *bithumb.Bithumb
	chStop   chan bool
	chInfo   chan<- bithumb.TickerInfo
	cointype string
}

func NewEagleEye(cb chan<- bithumb.TickerInfo, cointype string) *EagleEye {
	instance := &EagleEye{
		bit:      bithumb.NewBithumb("test", "sec"),
		chStop:   make(chan bool),
		chInfo:   cb,
		cointype: cointype,
	}
	return instance
}

func (e *EagleEye) Scout() {
	bContinue := true

	e.sendPrice()

	ticker := time.NewTicker(time.Second * 10)
	for bContinue {
		select {
		case <-ticker.C:
			e.sendPrice()
		case <-e.chStop:
			bContinue = false
		}
	}
}

func (e *EagleEye) sendPrice() {
	var info bithumb.TickerInfo
	e.bit.GetPrice(e.cointype, &info)
	e.chInfo <- info
}

func (e *EagleEye) Stop() {
	e.chStop <- true
}
