package random

import (
	"fmt"
	"math/rand"
	"time"
)

// 生成订单号
func GenerateOrderNumber() string {
	// 获取当前日期时间作为订单号的一部分
	currentTime := time.Now().Format("200601021504")

	time.Sleep(time.Nanosecond)

	// 生成随机数作为订单号的一部分
	rand.Seed(time.Now().UnixNano())
	randomPart := fmt.Sprintf("%06d", rand.Intn(1000000))

	// 构建订单号
	orderNumber := "NO" + currentTime + randomPart

	return orderNumber
}
