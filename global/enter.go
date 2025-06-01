package global

import (
	"fast_gin/config"

	"github.com/redis/go-redis/v9"
	"github.com/zhenghaoz/gorse/client"
	"gorm.io/gorm"
)

const Version = "0.0.2"

var (
	Config *config.Config
	DB     *gorm.DB
	Redis  *redis.Client
	Gorse  *client.GorseClient
)
