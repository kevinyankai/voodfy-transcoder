package models

import (
	"github.com/Voodfy/voodfy-transcoder/internal/settings"
	"github.com/go-redis/redis"
)

var db *Models

// Models struct ...
type Models struct {
	Redis *redis.Client
}

// InitDB return instance redis client
func InitDB() {
	redis := redis.NewClient(&redis.Options{
		Addr: settings.RedisSetting.Host,
	})
	db = &Models{Redis: redis}
}
