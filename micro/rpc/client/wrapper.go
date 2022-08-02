package client

import "context"
// CallFunc represents the individual call func
type CallFunc func(ctx context.Context, req Request, rsp interface{}, opts CallOptions) error

// CallWrapper is a low level wrapper for the CallFunc
type CallWrapper func(CallFunc) CallFunc