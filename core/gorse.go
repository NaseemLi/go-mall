package core

import (
	"fast_gin/global"

	"github.com/sirupsen/logrus"
	"github.com/zhenghaoz/gorse/client"
)

func InitGorse() {
	gorse := global.Config.Gorse

	if !gorse.Enable {
		return
	}

	if gorse.Addr == "" {
		logrus.Fatalf("Gorse地址未配置")
		return
	}

	global.Gorse = client.NewGorseClient(gorse.Addr, gorse.ApiKey)
	logrus.Infof("Gorse连接成功,访问地址: %s", gorse.Addr)
}
