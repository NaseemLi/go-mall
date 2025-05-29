package payser

import (
	"fast_gin/global"
	"fmt"
)

type InlinePay struct {
}

func (InlinePay) Pay(no string, price int) (payUrl string, err error) {
	return fmt.Sprintf("%s?no=%s", global.Config.Pay.WebPayUrl, no), err
}
