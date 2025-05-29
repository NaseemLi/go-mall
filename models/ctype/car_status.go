package ctype

type CarStatus int8

const (
	CarStatusPending CarStatus = 0 // 待下单
	CarStatusLocked  CarStatus = 2 // 已锁定（下单中，未支付）
)
