package venice

import (
	"bufio"
	"os"
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

func main_temp() {
	bit := NewBitory()
	chPrice := make(chan float64)
	chStop := make(chan bool)

	eye := NewHawkEye(chPrice)
	bContinue := true

	go eye.Scout()
	go func() {
		userInput := bufio.NewScanner(os.Stdin)
		for userInput.Scan() {
			inputData := userInput.Text()
			if "exit" == inputData {
				break
			}
		}
		eye.Stop()
		chStop <- true
	}()

	for bContinue {
		select {
		case <-chStop:
			bContinue = false
		case price := <-chPrice:
			bit.AddInfo(price)
			bit.Print()
		}
	}

}
