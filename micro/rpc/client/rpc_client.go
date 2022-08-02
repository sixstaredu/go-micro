package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/rpc"
	"net/rpc/jsonrpc"
	"github.com/sixstaredu/go-micro/micro/core/debug"
	"github.com/sixstaredu/go-micro/micro/core/errors"
	"time"
)

type rpcClient struct {
	opts *dialOptions
	mp *managePool
	id int64
}
// 创建客户端
func NewClient(opt ...DialOption) (client *rpcClient) {
	// 设置配置
	opts := newDialOptions()
	for _, o := range opt {
		o.apply(opts)
	}
	// 初始化rpc
	client =  &rpcClient{
		opts: opts,
		mp: newMangePool(),
	}

	poolOpts := PoolOptions{
		Size:                opts.poolsize,
		TTl:                 opts.poolTTl,
	}

	for serverName, server := range opts.servers {
		debug.PrintDirExePos(dir+"NewClient","创建 %v 连接池", serverName)

		poolOpts.CreateConnectHandle = client.newConnect(serverName, server)
		poolOpts.Id = serverName

		pool,err := initPool(poolOpts)
		if err == ErrCreateConnHandleNotExit {
			debug.PrintErrDirExePos(dir+"NewClient", err, "创建 %v 连接池出现异常", serverName)
			continue
		}
		client.mp.Add(serverName, pool)
	}

	return
}
// 回收连接
func (c *rpcClient) ConnRelease(serverName string, conn Conn) {
	pool, ok := c.mp.Get(serverName)
	if !ok {
		return
	}
	pool.Release(conn)
}
// 根据服务名,创建新的连接
func (c *rpcClient) NewConnect(serverName string) (Conn, error) {

	pool, ok := c.mp.Get(serverName)
	if !ok {
		debug.PrintErrDirExePos(dir+"NewConnect", ErrNotServer, "获取 %v 连接池错误", serverName)
		return nil, errors.NotFound("go-micro/rpc/client/rpcClient.NewConnect", "client get %s server not found", serverName)
	}

	ctx,_ := context.WithTimeout(context.TODO(), c.opts.connTimeout)
	conn, err := pool.Get(ctx)
	if err!=nil {
		debug.PrintErrDirExePos(dir+"NewConnect", err, "从连接池中获取 %v 服务连接出现问题", serverName)
		return nil, err
	}

	return conn, err
}

func (c *rpcClient) call(ctx context.Context, req Request, resp interface{}, callOption CallOptions) (err error) {
	// 创建新的连接，获取新的连接
	conn, err := c.NewConnect(req.Service())
	defer func() {
		//if err != nil {
		//
		//}
		c.ConnRelease(req.Service(), conn)
	}()
	// 判断连接是否获取成功
	if err != nil {
		debug.PrintErrDirExePos(dir, err, "获取服务连接 %v 异常", req.Service())
		return err
	}

	// 调度服务
	return conn.Call(ctx, req, resp, callOption)
}

func (c *rpcClient) Call(ctx context.Context, req Request, resp interface{}, callOption ...CallOption) error {
	// 覆盖执行操作信息
	callOpts := c.opts.callOptions
	for _, opt := range callOption {
		opt(&callOpts)
	}
	// 是否设置超时
	d, ok := ctx.Deadline()
	if !ok {
		// no deadline so we create a new one
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, callOpts.requestTimeout)
		defer cancel()
	} else {
		opt := WithRequestTimeout(d.Sub(time.Now()))
		opt(&callOpts)
	}

	// 复制call方法
	rcall := c.call

	// 执行中间件方法
	for i := len(callOpts.callWrappers); i > 0; i-- {
		rcall = callOpts.callWrappers[i-1](rcall)
	}

	// 执行-失败重试
	retries := callOpts.retries
	ch := make(chan error, retries + 1)
	var gerr error

	for i := 0; i < retries; i++ {
		go func(i int) {
			ch <- rcall(ctx, req, resp, callOpts)
		}(i)

		select {
		case <-ctx.Done():
			return errors.Timeout("go-micro/rpc/client/rpcClient.Call", "client call server %s.%s: timeout",req.Service() , req.Method())
		case err := <- ch:
			// if the call succeeded lets bail early
			if err == nil {
				return nil
			}

			retry, rerr := callOpts.retry(ctx, req, i, err)
			if rerr != nil {
				return rerr
			}

			if !retry {
				return err
			}

			gerr = err
		}
	}

	return gerr
}

func (c *rpcClient) NewRequest(serverName string, serverMethod string, req interface{}, opts ...RequestOption) Request {
	return newRequest(serverName, serverMethod, req, opts...)
}

func (c *rpcClient) newConnect(serverName string, s *Server) CreateConnectHandle {
	return func() (Conn, error) {
		c.id ++
		// 建立连接
		client, err := c.getClient(s)

		if err != nil {
			debug.PrintErrDirExePos(dir+"newConnect", err, "创建 %v 服务出现异常", serverName)
			return &connect{
				id:  c.id,
				err: err,
			}, err
		}

		debug.PrintDirExePos(dir+"newConnect", "连接服务 %v ", serverName)

		return &connect{
			client:  client,
			id:      c.id, // 应该是利用uuid解决 ，作为作业实现
			addr:    s.Address,
			err:     nil,
			created: time.Now(),
		}, nil
	}
}
// 获取连接
func (c *rpcClient) getClient(s *Server) (client *rpc.Client,err error) {
	defer func() {
		if err != nil {
			err = errors.New("go-micro/rpc/client/rpcClient.getClient", err.Error(), 500)
		}
	}()

	if !s.openssl {
		return jsonrpc.Dial(s.Network, s.Address)
	}
	// 开启tls认证
	certBytes, err := ioutil.ReadFile(s.CartFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(certBytes)

	// 连接
	conn, err := tls.Dial("tcp", s.Address,  &tls.Config{
		RootCAs: certPool,
		ServerName: s.TlsServerName,
	})

	if err != nil {
		return nil, err
	}

	return jsonrpc.NewClient(conn), err
}

// 调度失败重试 -》 重试， 重试过程中如果存在问题，我们可以切换其他的连接
