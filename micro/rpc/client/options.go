package client

import "time"

var (
	dir = "core/rpc/client/"
)

var (
	// 默认连接池大小
	DefaultPoolSize = 5
	// 默认连接池生命周期
	DefaultPoolTTl = 10 * time.Minute
	// 默认连接超时时间
	DefaultConnTimeout = 3 * time.Second
	// 默认重试次数
	DefaultRequestTimeout = 3 * time.Second
	// 默认重试次数
	DefaultRetries = 1
	// 默认重试验证方法
	DefaultRetry = RetryAlways
)

type Server struct {
	openssl bool
	CartFile string
	TlsServerName string
	Network string
	Address string
}

type dialOptions struct {
	// 需要连接的服务
	servers map[string]*Server
	// 连接池大小
	poolsize int
	// 连接生命周期
	poolTTl time.Duration
	// 连接超时
	connTimeout time.Duration
	// 调研属性
	callOptions CallOptions
}

type CallOptions struct {
	callWrappers []CallWrapper
	// Address of remote hosts
	address []string
	// 根据异常校验是否重试
	retry RetryFunc
	// 重试次数
	retries int
	// 请求超时
	requestTimeout time.Duration
}

func newDialOptions() *dialOptions {
	return &dialOptions{
		servers:     make(map[string]*Server),
		poolsize:    DefaultPoolSize,
		poolTTl:     DefaultPoolTTl,
		connTimeout: DefaultConnTimeout,
		callOptions: CallOptions{
			retry:          DefaultRetry,
			retries:        DefaultRetries,
			requestTimeout: DefaultRequestTimeout,
		},
	}
}

type DialOption interface {
	apply(*dialOptions)
}

type funcDialOption struct {
	f func(*dialOptions)
}
func newFuncDialOption(f func(*dialOptions)) *funcDialOption {
	return &funcDialOption{
		f: f,
	}
}
func (fdo *funcDialOption) apply(do *dialOptions) {
	fdo.f(do)
}

// 设置服务
func SetServer(name string, server *Server) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		if server.CartFile != "" && server.TlsServerName != "" {
			server.openssl = true
		}

		o.servers[name] = server
	})
}
// 设置连接池大小
func SetPoolSize(size int) DialOption {
	return newFuncDialOption(func(options *dialOptions) {
		options.poolsize = size
	})
}
// 设置连接生命周期
func SetPoolTTl(ttl time.Duration) DialOption {
	return newFuncDialOption(func(options *dialOptions) {
		options.poolTTl = ttl
	})
}
// 设置超时
func SetConnTimeout(timeout time.Duration) DialOption {
	return newFuncDialOption(func(options *dialOptions) {
		options.connTimeout = timeout
	})
}
// 全局设置 -- 请求超时
func RequestTimeout(timeout time.Duration) DialOption {
	return newFuncDialOption(func(options *dialOptions) {
		options.callOptions.requestTimeout = timeout
	})
}
// 重试次数 全局
func Retries(retries int) DialOption {
	return newFuncDialOption(func(options *dialOptions) {
		options.callOptions.retries = retries
	})
}
// 重试次数 全局
func Retry(fn RetryFunc) DialOption {
	return newFuncDialOption(func(options *dialOptions) {
		options.callOptions.retry = fn
	})
}

type CallOption func(options *CallOptions)

// 全局设置 -- 请求超时
func WithRequestTimeout(timeout time.Duration) CallOption {
	return func(options *CallOptions) {
		options.requestTimeout = timeout
	}
}
// 重试次数 全局
func WithRetries(retries int) CallOption {
	return (func(options *CallOptions) {
		options.retries = retries
	})
}
// 重试次数 全局
func WithRetry(fn RetryFunc) CallOption {
	return (func(options *CallOptions) {
		options.retry = fn
	})
}
func WrapCall(cw ...CallWrapper) CallOption {
	return func(o *CallOptions) {
		o.callWrappers = append(o.callWrappers, cw...)
	}
}
func WithWrapCall(cw ...CallWrapper) DialOption {
	return newFuncDialOption(func(options *dialOptions) {
		options.callOptions.callWrappers = append(options.callOptions.callWrappers, cw...)
	})
}

type requestOptiones struct {

}
type RequestOption func(optiones *requestOptiones)