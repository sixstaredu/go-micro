package server

import "context"
// 修改处
type HandlerFunc func(ctx context.Context, req *Request, argv ,rsp interface{}) error

type HandlerWrapper func(HandlerFunc) HandlerFunc
