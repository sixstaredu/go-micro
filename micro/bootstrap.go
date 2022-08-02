package micro

import (
	"fmt"
	"time"
	"strconv"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go"

	"github.com/sixstaredu/go-micro/micro/config"
	"github.com/sixstaredu/go-micro/micro/core/debug"
	"github.com/sixstaredu/go-micro/micro/core/cache"
	"github.com/sixstaredu/go-micro/micro/core/log"
	"github.com/sixstaredu/go-micro/micro/core/model"
	"github.com/sixstaredu/go-micro/micro/core/redis"
	"github.com/sixstaredu/go-micro/micro/core/validate"
	"github.com/sixstaredu/go-micro/micro/rpc/client"
)

func Init(cfgPath string) {
	debug.SetPrintPrefix("[go-micro]")


	initConfig(cfgPath)
	initLog(Config.Log)
	initCache(Config.Cache)
	initModel(Config.Mysql)
	initRedis(Config.Redis)
	loadValidator()
	InitJaeger(Config.App.ServerName, Config.Jaeger.Address)
}

func initConfig(cfgPath string)  {
	Viper = initViper(&Config, cfgPath)
}

func initLog(cfg *log.Config)  {
	Logs = log.InitLogger(cfg)
}

func initCache(cfg *cache.Config)  {
	// 初始化缓存
	if cfg.Default == "freecache" {
		cache.CacheManager = cache.NewCache(cache.NewFreeCache(cfg))
	}
}

func loadValidator()  {
	// 初始化模型
	v, _ := binding.Validator.Engine().(*validator.Validate)
	validate.InitValidate(v, validators, "zh")
}

func initModel(cfg *model.Config)  {
	// 初识mysql
	DB = model.InitDb(cfg)
}
func initRedis(cfg *redis.Config)  {
	Redis = redis.InitRedisPool(cfg)
}
func InitRpcClient(cfg *config.RpcClient, opts ...client.DialOption) {
	if cfg == nil {
		return
	}
	// 初始化rpc client
	if len(cfg.Servers) > 0 {
		for k, v := range cfg.Servers {
			debug.DD("v = %v", v)
			opts = append(opts, client.SetServer(k, &client.Server{
				CartFile:      v.CartFile,
				TlsServerName: v.TlsServerName,
				Network:       v.Network,
				Address:       v.Address,
			}))
		}
	}
	RpcClient = client.NewClient(opts...)
}

func InitJaeger(service, address string)  {
	var err error
	if service == "" {
		service = strconv.Itoa(int(time.Now().Unix()))
	}

	cfg := jaegercfg.Configuration{
		Sampler:             &jaegercfg.SamplerConfig{
			Type:                     jaeger.SamplerTypeConst,
			Param:                    1,
		},
		Reporter:            &jaegercfg.ReporterConfig{
			LogSpans:                   true,
			// 将span发送jaeger-collector的服务中
			CollectorEndpoint:          fmt.Sprintf("http://%s/api/traces", address),
		},
	}

	Jaeger, err = cfg.InitGlobalTracer(service, jaegercfg.Logger(jaeger.StdLogger))

	if err != nil {
		panic(fmt.Sprintf("Error: connect jaeger:%v \n", err))
	}
}
