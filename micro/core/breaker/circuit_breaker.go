package breaker

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrOpenState 		= errors.New("circuit breaker is open")
	ErrTooManyRequest 	= errors.New("too many requests")
)

type CircuitBreaker struct {
	name 				string
	// 最大请求次数
	//maxRequests 	uint32
	//// 熔断超时的时限
	//timeout				time.Duration
	//// 自定义熔断验证方法，判断是否开启熔断
	//readyToTrip			func(counts Counts) bool
	//// 在熔断的状态发送改变的时候执行
	//onSstateChange 		func(name string, from State, to State)
	opt 				*option

	mutex 				sync.Mutex
	// 状态
	state 				State
	// 统计熔断器的次数
	counts 				Counts
	// 记录熔断时限
	expiry 				time.Time
}
// 创建熔断器
func NewBreaker(opts ...Option) *CircuitBreaker {
	cb := new(CircuitBreaker)

	opt := newOption()

	for _ , o := range opts {
		o(opt)
	}

	cb.name = opt.name
	cb.opt = opt

	return cb
}

func (cb *CircuitBreaker) Name() string {
	return cb.name
}

func (cb *CircuitBreaker) Counts() Counts {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	return cb.counts
}

func (cb *CircuitBreaker) State() State {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	return cb.currentState()
}

//DoWithAcceptable(req func() error, acceptable Acceptable) error
//
//DoWithFallback(req func() error, fallback func(err error) error) error
//
//DoWithFallbackAcceptable(req func() error, fallback func(err error) error, acceptable Acceptable) error

func (cb *CircuitBreaker) Do(req func() error) error {
	return cb.do(req, nil, defaultAcceptable)
}
func (cb *CircuitBreaker) DoWithAcceptable(req func() error, acceptable Acceptable) error {
	return cb.do(req, nil, acceptable)
}
func (cb *CircuitBreaker) DoWithFallback(req func() error, fallback func(err error) error) error {
	return cb.do(req, fallback, defaultAcceptable)
}
func (cb *CircuitBreaker) DoWithFallbackAcceptable(req func() error, fallback func(err error) error, acceptable Acceptable) error {
	return cb.do(req, fallback, acceptable)
}
// 熔断器执行请求
// req			: 正常执行的方法
// fallback		: 熔断后降级的方法
// acceptable	: 则是判断正常执行方法之后的异常是否可以通过，否则失败
func (cb *CircuitBreaker) do(req func() error, fallback func(err error) error, acceptable func(err error) bool) error {
	// 获取当前熔断的状态，并判断是否处理方法
	if err := cb.accept(); err != nil {
		// 当处于熔断状态就执行降级方法
		if fallback != nil {
			return fallback(err)
		}

		return err
	}
	// 注意要加，避免正常的执行方法req出现问题；做好panic异常处理
	defer func() {
		if e := recover(); e != nil {
			cb.failure()
			panic(e)
		}
	}()

	// 执行实际方法
	err := req()
	// 判断异常是否视为调度失败
	if acceptable(err) {
		// 成功
		cb.success()
	} else {
		// 失败
		cb.failure()
	}

	return err
}
// 判断当前的熔断器是否接受处理方法
func (cb *CircuitBreaker) accept() error {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	// 获取当前状态
	state := cb.currentState()
	// 根据熔断状态判断
	if state == StateOpen {
		return ErrOpenState
	} else if state == StateHalfOpen && cb.counts.Requests >= cb.opt.maxRequests {
		return ErrTooManyRequest
	}

	cb.counts.request()
	// 需要考虑并发问题
	return nil
}
// 请求成功后执行
func (cb *CircuitBreaker) success() {
	switch cb.state {
	case StateClosed:
		// 统计成功的次数
		cb.counts.success()
	case StateHalfOpen:
		cb.counts.success()
		if cb.counts.ConsecutiveSuccesses >= cb.opt.maxRequests {
			cb.setState(StateClosed)
		}
	}
}
// 请求失败后执行
func (cb *CircuitBreaker) failure() {
	switch cb.state {
	case StateClosed:
		// 统计失败的次数
		cb.counts.failure()
		// 判断是否开启熔断
		if cb.opt.readyToTrip(cb.counts) {
			cb.setState(StateOpen)
		}
	case StateHalfOpen:
		// 半开放状态下的失败，就直接设置为开启状态
		cb.setState(StateOpen)
	}
}
// 获取当前是否熔断
func (cb *CircuitBreaker) currentState() State {
	now := time.Now()

	switch cb.state {
	case StateClosed:
	case StateOpen:
		// 如果熔断超过时限，则设置为半开放状态
		if cb.expiry.Before(now) {
			cb.setState(StateHalfOpen)
		}
	}

	return cb.state
}
// 设置熔断
func (cb *CircuitBreaker) setState(state State) {
	if cb.state == state {
		return
	}

	// 获取熔断记录
	prev := cb.state
	cb.state = state

	cb.toNewExpiry(time.Now())

	if cb.opt.onSstateChange != nil {
		cb.opt.onSstateChange(cb.name, prev, state)
	}
}

func (cb *CircuitBreaker) toNewExpiry(now time.Time) {
	// 重置,计数
	cb.counts.clear()

	var zero time.Time

	switch cb.state {
	case StateClosed:
	case StateOpen:
		cb.expiry = now.Add(cb.opt.timeout)
	default:
		cb.expiry = zero
	}
}