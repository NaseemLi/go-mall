package usercenterapi

import (
	addrapi "fast_gin/api/user_center_api/addr_api"
	collectapi "fast_gin/api/user_center_api/collect_api"
	lookapi "fast_gin/api/user_center_api/look_api"
)

type UserCenterApi struct {
	LookApi    lookapi.LookApi
	CollectApi collectapi.CollectApi
	AddrApi    addrapi.AddrApi
}
