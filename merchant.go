package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/nanunsin/bithumbb/bithumb"
)

func sumArray(arr []interface{}, div int) (ret float64) {

	ret = 0.0
	if div == 0 {
		return
	}

	for i := 0; i < len(arr); i++ {
		ret += arr[i].(float64)
	}

	//log.Printf("[sumArray] %.3f\n", ret/(float64)div)

	ret /= (float64)(div)
	return
}

func main() {
	fmt.Println("Start")

	/*
		bit := bithumb.NewBithumb("test", "sec")

		average9 := ring.New(10)
		average25 := ring.New(26)
		averageMACD := ring.New(9)
	*/

	chStop := make(chan bool)

	now := time.Now()
	test := now.Add(5 * time.Second).Round(10 * time.Second)
	timer := time.NewTimer(test.Sub(now))
	<-timer.C

	list := NewLimitList(50)
	signal := NewLimitList(50)

	go func() {
		var info bithumb.WMP
		bContinue := true

		bit := bithumb.NewBithumb("test", "sec")
		bit.GetETHPrice(&info)
		list.Push(info.Price)

		ticker := time.NewTicker(time.Second * 10)
		for bContinue {
			select {
			case <-ticker.C:
				bit.GetETHPrice(&info)
				list.Push(info.Price)

				/*
					if i > 26 {
						average9 := sumArray(datas[i-9:i], 9)
						average26 := sumArray(datas[i-26:i], 26)

						fmt.Println(math.Abs(average26 - average9))
					}
				*/
				length := list.Len()
				if length > 26 {
					getdata := list.ToSlice()
					average9 := sumArray(getdata[length-9:], 9)
					average26 := sumArray(getdata[length-26:], 26)

					macd := average9 - average26
					signal.Push(macd)

					macdVals := signal.ToSlice()
					averageMACD := sumArray(macdVals, signal.Len())
					fmt.Printf("[%.0f ] %.3f | %.3f | %.3f\n", info.Price, macd, averageMACD, macd-averageMACD)
				}

			case <-chStop:
				fmt.Println("stop")
				bContinue = false
			}
		}

		ticker.Stop()

	}()

	bContinue := true
	scanner := bufio.NewScanner(os.Stdin)

	for bContinue {
		scanner.Scan()
		readtext := scanner.Text()
		fmt.Println(readtext) // Println will add back the final '\n'

		if "exit" == readtext {
			chStop <- true
			bContinue = false
		}
	}

	fmt.Println("End")
}
