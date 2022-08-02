package idempotency

import (
	"context"
	"errors"
	"fmt"
	"github.com/sixstaredu/go-micro/micro/rpc/client"
	"github.com/sixstaredu/go-micro/micro/rpc/server"
	"github.com/sixstaredu/go-micro/micro/tool"
	"time"
)

type Idempotent interface {
	TryAcquire(id string, timeout time.Duration) bool
	Comfirm(id string)
}

var (
	tkey = "go-micro-task-key"
	dkey = "go-micro-dispatch-key"
)

var (
	Confirmed = "-1"
	TimeFormat = "20060102150405"
)

// 请求任务id的生成
func ContextWithVal(ctx context.Context) context.Context {
	return context.WithValue(ctx, tkey, tool.GenId())
}
// 客户端中间件
func CallWrapper(call client.CallFunc) client.CallFunc {
	return func(ctx context.Context, req client.Request, rsp interface{}, opts client.CallOptions) error {
		// 先获取任务id
		v := ctx.Value(tkey)
		// 生成任务id
		id := fmt.Sprintf("%d.%s.%s", v, req.Service(), req.Method())
		// 设置到header中
		// 发送给服务端
		req.Header().Set(dkey, id)

		return call(ctx, req, rsp, opts)
	}
}
// 服务端中间件
func NewHandlerWrapper(timeout time.Duration, i Idempotent) server.HandlerWrapper {
	return func(call server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req *server.Request, argv, rsp interface{}) error {

			// 获取id
			id := req.Header.Get(dkey)
			if id == "" {
				// 没有做幂等性的请求，直接执行
				return call(ctx, req, argv, rsp)
			}

			// 幂等性判断
			if i.TryAcquire(id, timeout) {
				err := call(ctx, req, argv, rsp)
				i.Comfirm(id)
				return err
			}

			return errors.New("存在其他程序正在执行该任务 id := "+id)
		}
	}
}