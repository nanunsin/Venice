package venice

type MACDInfo struct {
	MACD   float64
	Signal float64
	Histo  float64
	Curve  float64
	Price  float64
}

type Jarvis interface {
	Update(price int64)
	Test(price int64)
	Print()
}

// Bitory is core manager
type Bitory struct {
	prices *LimitList

	// merchant
	core    []Jarvis
	coreCnt int
}

func NewBitory() *Bitory {
	instance := &Bitory{
		prices:  NewLimitList(50000),
		core:    make([]Jarvis, 10),
		coreCnt: 0,
	}
	return instance
}

func (b *Bitory) AddCore(core Jarvis) bool {
	if b.coreCnt == 10 {
		return false
	}
	coreID := b.coreCnt
	b.core[b.coreCnt] = core
	b.coreCnt++

	if b.prices.Len() > 0 {
		prices := b.prices.ToSlice()
		for _, p := range prices {
			b.core[coreID].Update(p.(int64))
		}
	}

	return true
}

func (b *Bitory) AddPrice(price int64) {
	b.prices.Push(price)
	for i := 0; i < b.coreCnt; i++ {
		b.core[i].Update(price)
	}
}

func (b *Bitory) Print() {
	for i := 0; i < b.coreCnt; i++ {
		b.core[i].Print()
	}
}

func (b *Bitory) Process() {
	b.Print()
}

/*
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

	b.data.MACD = rs - rl

	b.sig.AddInfo(b.data.MACD)

	if !b.bSignal {
		if b.sig.Len() >= 9 {
			b.bSignal = true
		} else {
			return
		}
	}

	b.data.Signal, _ = b.sig.Avg()
	phisto := b.data.Histo
	b.data.Histo = b.data.MACD - b.data.Signal
	b.data.Price = data

	b.curve = b.data.Histo - phisto
	b.data.Curve = b.curve

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
			fmt.Printf("MACD: %.3f, Signal : %.3f, Histo: %.3f, C:%.3f\n", b.data.MACD, b.data.Signal, b.data.Histo, b.curve)
		} else {
			fmt.Printf("MACD: %.3f, Signal : %.3f\n", b.data.MACD, b.data.Signal)
		}
	}
}

func (b *Bitory) Process() {
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
*/
