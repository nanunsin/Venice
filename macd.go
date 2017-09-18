package venice

import "fmt"

type MACD struct {
	macd, signal, histo, curve float64
	Price                      float64
}

type Jarvis interface {
	Update(data MACD)
}

type Bitory struct {
	s, l, sig      *RQueue
	data           MACD
	curve          float64
	bMACD, bSignal bool
	// merchant
	core  Jarvis
	bCore bool
}

func NewBitory() *Bitory {
	instance := &Bitory{
		s:       NewRQueue(10),
		l:       NewRQueue(26),
		sig:     NewRQueue(9),
		bMACD:   false,
		bSignal: false,
		bCore:   false,
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
	b.data.Price = data

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
		} else {
			fmt.Printf("MACD: %.3f, Signal : %.3f,\n", b.data.macd, b.data.signal)
		}
	}
}

func (b *Bitory) Process() {
	fmt.Println("Process!")
	if b.bMACD && b.bSignal && b.bCore {
		b.core.Update(b.data)
	}
}

func (b *Bitory) SetCore(core Jarvis) {
	b.core = core
	b.bCore = true
}

func (b *Bitory) ResetCore(core Jarvis) {
	b.bCore = false
}
