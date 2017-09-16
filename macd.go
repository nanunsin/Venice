package venice

import (
	"fmt"
	"math"

	"github.com/fatih/color"
	"github.com/nanunsin/bithumbb/bithumb"
)

type MACD struct {
	macd, signal, histo, curve float64
	price                      float64
}

type Bitory struct {
	s, l, sig      *RQueue
	data           MACD
	curve          float64
	bMACD, bSignal bool
	thinker        *Thinker
}

func NewBitory() *Bitory {
	instance := &Bitory{
		s:       NewRQueue(10),
		l:       NewRQueue(26),
		sig:     NewRQueue(9),
		bMACD:   false,
		bSignal: false,
		thinker: NewThinker(),
	}
	return instance
}

func (b *Bitory) AddInfo(data float64) {
	b.s.AddInfo(data)
	b.l.AddInfo(data)

	if !b.bMACD {
		if b.l.Len() >= 26 {
			b.bMACD = true
		} else {
			return
		}
	}
	rs, _ := b.MACD10()
	rl, _ := b.MACD26()

	b.data.macd = rs - rl

	b.sig.AddInfo(b.data.macd)

	if !b.bSignal {
		if b.sig.Len() >= 9 {
			b.bSignal = true
		} else {
			return
		}
	}

	b.data.signal, _ = b.sig.Avg()
	phisto := b.data.histo
	b.data.histo = b.data.macd - b.data.signal
	b.data.price = data

	b.curve = b.data.histo - phisto
	b.data.curve = b.curve

	fmt.Printf("[%.f] ", data)

}

func (b *Bitory) MACD10() (float64, bool) {
	return b.s.Avg()
}

func (b *Bitory) MACD26() (float64, bool) {
	return b.l.Avg()
}

func (b *Bitory) MACDSignal() (float64, bool) {
	return b.sig.Avg()
}

func (b *Bitory) Print() {
	if b.bMACD {
		if b.bSignal {
			fmt.Printf("MACD: %.3f, Signal : %.3f, Histo: %.3f, C:%.3f\n", b.data.macd, b.data.signal, b.data.histo, b.curve)
			b.thinker.Add(b.data)
		} else {
			fmt.Printf("MACD: %.3f, Signal : %.3f,\n", b.data.macd, b.data.signal)
		}
	}
}

type Thinker struct {
	oldMACD MACD
	bStart  bool
	bBuy    bool
	bit     *bithumb.Bithumb
}

func NewThinker() *Thinker {
	var apikey = "fd8871256b116a150f9ef1390f909105"
	var apisecret = "2605ae4bcd5d9c2cb43b2341d68b0909"
	return &Thinker{
		bStart: false,
		bit:    bithumb.NewBithumb(apikey, apisecret),
	}
}

func (th *Thinker) Add(data MACD) {
	if !th.bStart {
		th.bStart = true
	} else {

		difMACD := data.macd - th.oldMACD.macd       // MACD 기울기
		difSignal := data.signal - th.oldMACD.signal // Signal 기울기
		bonus := 50
		if math.Abs(data.curve) > 100 {
			bonus += 50
		}
		/*
			color.New(color.FgBlue).Printf("%.3f\t", difMACD)
			color.New(color.FgCyan).Printf("%.3f\t", difSignal)

			if data.histo > 0 {
				if data.curve < -100 && th.oldMACD.curve < -100 {
					if th.bBuy {
						fmt.Printf("Sell, %.f\n", data.price)
						th.bit.SellPlaceETH(int(data.price)-bonus, 0.1)
						th.bBuy = false
					}
				} else {
					if !th.bBuy {
						fmt.Printf("Buy, %.f\n", data.price)
						th.bit.BuyPlaceETH(int(data.price)+bonus, 0.1)
						th.bBuy = true
					}
				}

			} else {
				if data.curve > 100 && th.oldMACD.curve > 100 {
					if !th.bBuy {
						fmt.Printf("Buy, %.f\n", data.price)
						th.bit.BuyPlaceETH(int(data.price)+bonus, 0.1)
						th.bBuy = true
					}
				} else {
					if th.bBuy {
						fmt.Printf("Sell, %.f\n", data.price)
						th.bit.SellPlaceETH(int(data.price)-bonus, 0.1)
						th.bBuy = false
					}
				}
			}
		*/

		if data.curve < -100 && th.oldMACD.curve < -100 {
			if th.bBuy {
				fmt.Printf("Sell, %.f\n", data.price)
				th.bit.SellPlaceETH(int(data.price)-bonus, 0.1)
				th.bBuy = false
			}
		} else if data.curve > 100 && th.oldMACD.curve > 100 {
			if !th.bBuy {
				fmt.Printf("Buy, %.f\n", data.price)
				th.bit.BuyPlaceETH(int(data.price)+bonus, 0.1)
				th.bBuy = true
			}
		} else {
			if difMACD > difSignal {
				color.New(color.FgRed).Println("Good")
				if difSignal > 0 { // 시그널이 증가 상태
					if !th.bBuy {
						fmt.Printf("Buy, %.f\n", data.price)
						th.bit.BuyPlaceETH(int(data.price)+bonus, 0.1)
						th.bBuy = true
					}
				}
			} else {
				color.New(color.FgGreen).Println("Bad")
				if data.histo < 0 {
					if th.bBuy {
						fmt.Printf("Sell, %.f\n", data.price)
						th.bit.SellPlaceETH(int(data.price)-bonus, 0.1)
						th.bBuy = false
					}
				}
			}
		}

	}

	th.oldMACD = data
}
