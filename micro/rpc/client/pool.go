package client

import (
	"context"
	"github.com/sixstaredu/go-micro/micro/core/debug"
	"github.com/sixstaredu/go-micro/micro/core/errors"
	"time"
)


/*
// 运用接口可以方便后续升级和维护
type Conn interface {
	// rpc调度服务方法
	Call(serverMethod string, req interface{}, resp interface{}) error
	// 关闭连接
	Close() error
	// 获取连接的创建时间
	Created() time.Time
	// 获取连接的地址
	Remote() string
	// 连接的异常记录
	Error() error
	// 连接id
	Id() string
}
 */

type Pool interface {
	// 获取连接
	Get(ctx context.Context) (Conn, error)
	// 释放连接
	Release(conn Conn)
	// 关闭连接
	Close()
}

// 定一个创建连接的方法
type CreateConnectHandle func() (Conn, error)

// 连接池的配置信息
type PoolOptions struct {
	Id string
	Size int
	TTl time.Duration
	CreateConnectHandle
}
// 管理连接池
type managePool struct {
	pools map[string]Pool
}
func newMangePool() *managePool {
	return &managePool{
		pools: make(map[string]Pool),
	}
}
func (mp *managePool) Add(tab string, pool Pool) {
	mp.pools[tab] = pool
}
func (mp *managePool) Get(tab string) (Pool, bool){
	pool, ok := mp.pools[tab]
	return pool, ok
}


// 基于chan实现的连接池
type pool struct {
	id string
	count int // 用来记录创建的连接
	size int
	ttl time.Duration
	conns chan Conn
	CreateConnectHandle
}

func initPool(options PoolOptions) (*pool, error) {
	// 连接池大小最少需要设置1个
	if options.Size <= 0 {
		return nil, errors.New("go-micro/rpc/client.initPool", options.Id+" cannot be set size to zero", 500)
	}

	p := &pool{
		id:					 options.Id,
		size:                options.Size,
		ttl:                 options.TTl,
		conns:               make(chan Conn, options.Size),
		CreateConnectHandle: options.CreateConnectHandle,
	}
	return p, p.init()
}
// 初识连接
func (p *pool) init() error {
	if p.CreateConnectHandle == nil {
		return ErrCreateConnHandleNotExit
	}

	debug.PrintDirExePos(dir+"init","连接数 %d", p.size)
	// 创建连接
	for i := 0; i < p.size; i++ {
		conn, err := p.CreateConnectHandle()
		if err != nil {
			return err
		}
		p.count++
		p.conns <- conn
	}
	return nil
}

// 获取连接
func (p *pool) Get(ctx context.Context) (Conn, error) {
	// 判断正在使用的连接，是否少于总连接数的范围

	if p.count < p.size {
		//fmt.Println("创建链接 : ", p.id)
		p.createConn()
	}

	for  {
		select {
		case conn := <- p.conns:
			if d := time.Since(conn.Created()); d > p.ttl {
				conn.Close()
				p.count--
				// 创建新的连接
				p.createConn()
				continue
			}
			return conn, nil
		case <- ctx.Done():
			return nil, errors.Timeout("go-micro/rpc/client/pool.Get", "pool %s get timeout", p.id)
		}
	}
}
// 释放连接
func (p *pool) Release(conn Conn)  {
	// 可能连接为nil
	if conn == nil {
		return
	}

	if conn.Error() == nil {
		p.conns <- conn
		return
	}
	conn.Close()
	p.count--

	p.createConn()
}

func (p *pool) createConn() {
	go func() {
		retry := 0
		// 创建的时候会出现异常所以用for
		for {
			if p.CreateConnectHandle == nil {
				return
			}
			conn, err := p.CreateConnectHandle()

			if retry > 10 {
				return
			}

			if err != nil {
				retry++
				time.Sleep(1 * time.Second)
				continue
			}
			p.count++
			p.conns <- conn

			return
		}
	}()
}
func (p *pool) Close() {
	// 作业完善它
}
