package micro

import (
	"github.com/garyburd/redigo/redis"
	"github.com/mojocn/base64Captcha"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io"
	"github.com/sixstaredu/go-micro/micro/config"
	"github.com/sixstaredu/go-micro/micro/rpc/client"
)

var (
	Config config.Config
	Viper *viper.Viper // 后面可能会对配置文件操作，可以通过它来实现
	Logs *zap.Logger
	DB *gorm.DB
	Redis *redis.Pool

	Jaeger io.Closer

	RpcClient client.RpcClient
)

var CaptchaStore = base64Captcha.DefaultMemStore

func Close()  {
	Jaeger.Close()
	Redis.Close()
}