package payser

type PayInterface interface {
	Pay(no string, price int) (payUrl string, err error)
}

func Pay(payType int8, no string, price int) (payUrl string, err error) {
	var payService PayInterface
	// 根据需要选择具体的支付方式
	switch payType {
	case 1:
		payService = &InlinePay{}
	case 2:
		payService = &WxPay{}
	case 3:
		payService = &AliPay{}
	}
	return payService.Pay(no, price)
}
