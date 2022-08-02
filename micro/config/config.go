package config

import (
	"github.com/sixstaredu/go-micro/micro/core/cache"
	"github.com/sixstaredu/go-micro/micro/core/log"
	"github.com/sixstaredu/go-micro/micro/core/model"
	"github.com/sixstaredu/go-micro/micro/core/redis"
)

/**
 * @Com www.github.com/sixstaredu
 * @Author 六星教育-shineyork老师
 */

type Config struct {
	App
	Sms
	Smsbao
	Captche `mapstructure:"captcha"`
	Pay
	Jaeger

	*RpcServer `mapstructure:"rpc_server"`
	*RpcClient `mapstructure:"rpc_client"`
	Redis *redis.Config `mapstructure:"redis`
	Mysql *model.Config `mapstructure:"mysql"`
	Cache *cache.Config `mapstructure:"cache"`
	Log *log.Config `mapstructure:"log"`

}
