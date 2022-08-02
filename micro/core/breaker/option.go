package breaker

import "time"

var (
	defaultTimeout				 = time.Duration(60) * time.Second
	defaultMaxRequest 	uint32	 = 1
)

func defaultReadyToTrip(counts Counts) bool {
	return counts.ConsecutiveFailures > 5
}

func defaultAcceptable (err error) bool {
	return err == nil
}

type option struct {
	name				string
	// 最大请求次数
	maxRequests 		uint32
	// 熔断超时的时限
	timeout				time.Duration
	// 自定义熔断验证方法，判断是否开启熔断
	readyToTrip			func(counts Counts) bool
	// 在熔断的状态发送改变的时候执行
	onSstateChange 		func(name string, from State, to State)
}

type Option func(opt *option)

func newOption() *option {
	return &option{
		name:           "",
		maxRequests:    defaultMaxRequest,
		timeout:        defaultTimeout,
		readyToTrip:    defaultReadyToTrip,
		onSstateChange: nil,
	}
}

func WithName(name string) Option {
	return func(opt *option) {
		opt.name = name
	}
}

func WithMaxRequest(maxRequest uint32) Option {
	return func(opt *option) {
		opt.maxRequests = maxRequest
	}
}

func WithTImeout(timeout time.Duration) Option {
	return func(opt *option) {
		opt.timeout = timeout
	}
}

func WithReadyToTrip(readyToTrip func(counts Counts) bool) Option {
	return func(opt *option) {
		opt.readyToTrip = readyToTrip
	}
}

func WithOnSstateChange(onSstateChange func(name string, from State, to State)) Option {
	return func(opt *option) {
		opt.onSstateChange = onSstateChange
	}
}
