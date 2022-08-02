package client

import (
	"context"
	"github.com/sixstaredu/go-micro/micro/core/errors"
	"time"
	"net/rpc"
)

type connect struct {
	client *rpc.Client
	id int64
	addr string
	err error
	created time.Time
}

// 调用
func (c *connect) Call(ctx context.Context, req Request, resp interface{}, callOption CallOptions) error {
	ch := make(chan error, 1)

	go func() {
		//ch <- c.client.Call(req.Method(), req.Body(), resp)

		ch <- c.client.Call(req.Method(), &Message{
			Header: req.Header(),
			Body:   req.Body(),
		}, resp)
	}()

	select {
	case err := <-ch:
		if err != nil {
			// 修改
			c.err = err
			return errors.Parse(err.Error())
		}
		return nil
	case <- ctx.Done():
		return errors.Timeout("go-micro/rpc/client", "server %s.%s: timeout",req.Service(),req.Method())
	}
}
// 关闭连接
func (c *connect) Close() error {
	return c.client.Close()
}

func (c *connect) Created() time.Time {
	return c.created
}

func (c *connect) Remote() string {
	return c.addr
}

func (c *connect) Error() error {
	return c.err
}

func (c *connect) Id() int64 {
	return c.id
}
