package payser

type WxPay struct {
}

func (WxPay) Pay(no string, price int) (payUrl string, err error) {
	return "https://example.com/wxpay/success", err
}
