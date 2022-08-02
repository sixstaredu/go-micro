package server

import (
	"crypto/tls"
	"net"
	"os"
	"github.com/sixstaredu/go-micro/micro/core/debug"
	"github.com/sixstaredu/go-micro/micro/core/errors"
)

var dir = "core/rpc/server/"

type RpcServer struct {
	opts serverOptions
	count int
	svr *Server
}

func NewRpcServer(opt ...ServerOption) *RpcServer {
	opts := defaultServerOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	return &RpcServer{
		opts: opts,
		svr: NewServer(opts),
	}
}
// 注册服务
func (s *RpcServer) Register(server interface{}) error {
	return s.svr.Register(server)
}
func (s *RpcServer) RegisterName(name string, rcvr interface{}) error {
	return s.svr.RegisterName(name, rcvr)
}
// 启动服务
func (s *RpcServer) Run(addr ...string) (err error) {

	defer func() {
		if err != nil {
			err = errors.New("go-micro/rpc/server/RpcServer.Run", err.Error(), 500)
		}
		debug.DE(err)
	}()

	address := resolveAddress(addr)

	debug.DD("Listening and serving TCP on %s\n", address)

	lis, err := s.listen(address)

	if err != nil {
		return
	}

	for {
		conn, err := lis.Accept()

		if err != nil {
			continue
		}

		go func(conn net.Conn) {
			s.svr.ServeCodec(NewServerCodec(conn))
		}(conn)
	}

	return
}
// 设置监听
func (s *RpcServer) listen(address string) (lis net.Listener, err error) {
	// 没有开启openssl
	if !s.opts.openssl {
		return net.Listen("tcp", address)
	}
	// 开启认证
	debug.DD("开启tls认证")

	cert, err := tls.LoadX509KeyPair("./cert/server.pem", "./cert/server.key")

	if err != nil {
		return
	}

	return tls.Listen("tcp", address,  &tls.Config{
		Certificates: []tls.Certificate{cert},
	})
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			debug.DD("Environment variable PORT=\"%s\"", port)
			return ":" + port
		}
		debug.DD("Environment variable PORT is undefined. Using port :8080 by default")
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too many parameters")
	}
}