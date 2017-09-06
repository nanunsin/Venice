package main

import (
	"fmt"
	"time"

	"github.com/nanunsin/bithumbb/bithumb"
)

type Ethory struct {
	s, l, sig      *RQueue
	histo          float64
	macd, signal   float64
	curve          float64
	bMACD, bSignal bool
}

func NewEthory() *Ethory {
	instance := &Ethory{
		s:       NewRQueue(10),
		l:       NewRQueue(26),
		sig:     NewRQueue(9),
		bMACD:   false,
		bSignal: false,
	}
	return instance
}

func (e *Ethory) AddInfo(data float64) {
	e.s.AddInfo(data)
	e.l.AddInfo(data)

	if !e.bMACD {
		if e.l.Len() >= 26 {
			e.bMACD = true
		} else {
			return
		}
	}
	rs, _ := e.MACD10()
	rl, _ := e.MACD26()

	e.macd = rs - rl

	e.sig.AddInfo(e.macd)

	if !e.bSignal {
		if e.sig.Len() >= 9 {
			e.bSignal = true
		} else {
			return
		}
	}

	e.signal, _ = e.sig.Avg()
	phisto := e.histo
	e.histo = e.signal - e.macd

	e.curve = e.histo - phisto

}

func (e *Ethory) MACD10() (float64, bool) {
	return e.s.Avg()
}

func (e *Ethory) MACD26() (float64, bool) {
	return e.l.Avg()
}

func (e *Ethory) MACDSignal() (float64, bool) {
	return e.sig.Avg()
}

func (e *Ethory) Print() {
	if e.bMACD {
		if e.bSignal {
			fmt.Printf("MACD: %.3f, Signal : %.3f, Histo: %3.f, C:%.3f\n", e.macd, e.signal, e.histo, e.curve)
		} else {
			fmt.Printf("MACD: %.3f, Signal : %.3f,\n", e.macd, e.signal)
		}
	}
}

type Bitory struct {
	list   *LimitList
	chStop chan bool
}

func NewBitory() *Bitory {
	instance := &Bitory{}
	instance.list = NewLimitList(200)
	instance.chStop = make(chan bool)

	return instance
}

func (b *Bitory) Run() {
	var info bithumb.WMP
	bContinue := true

	bit := bithumb.NewBithumb("test", "sec")
	bit.GetETHPrice(&info)
	b.list.Push(info.Price)
	ticker := time.NewTicker(time.Second * 10)

	for bContinue {
		select {
		case <-ticker.C:
			bit.GetETHPrice(&info)
			b.list.Push(info.Price)

		case <-b.chStop:
			fmt.Println("stop")
			bContinue = false
		}
	}
}

func (b *Bitory) Stop() {
	b.chStop <- true
}

func (b *Bitory) SumArray(div int) (ret float64) {
	ret = 0.0
	if div == 0 {
		return
	}

	data := b.list.ToSlice()
	dataLen := len(data)
	if div > dataLen {
		fmt.Printf("Len : %d", dataLen)
		return
	}

	for _, d := range data[dataLen-div:] {
		ret += d.(float64)
	}
	ret /= (float64)(div)
	return
}
