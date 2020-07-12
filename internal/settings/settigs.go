package settings

import (
	"log"
	"time"

	"gopkg.in/ini.v1"
)

var cfg *ini.File

// App struct used to bind app
type App struct {
	Tag                 string
	QueueEnabled        bool
	SentryDNS           string
	LivepeerBroadcaster string
	LivepeerToken       string
	LivepeerMode        string

	RuntimeRootPath string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}

// AppSetting instance from  app
var AppSetting = &App{}

// Server struct used to bind store
type Server struct {
	RunMode     string
	BucketMount string
}

// ServerSetting instance from server
var ServerSetting = &Server{}

// IPFS struct used to bind store
type IPFS struct {
	Gateway string
	Origin  string
}

// IPFSSetting instance from server
var IPFSSetting = &IPFS{}

// Redis struct used to bind redis
type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

// RedisSetting instance  from redis
var RedisSetting = &Redis{}

// Influxdb struct used to bind redis
type Influxdb struct {
	Host     string
	Password string
	User     string
	DB       string
}

// InfluxdbSetting instance  from redis
var InfluxdbSetting = &Influxdb{}

// Setup initialize the configuration instance
func Setup() {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("redis", RedisSetting)
	mapTo("ipfs", IPFSSetting)
	mapTo("influxdb", InfluxdbSetting)

	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}