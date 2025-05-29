package payser

type AliPay struct {
}

func (AliPay) Pay(no string, price int) (payUrl string, err error) {
	return "https://example.com/alipay/success", err
}
