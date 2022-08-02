package breaker

import "fmt"

// 定义状态

type State int

const (
	StateClosed State = iota
	StateHalfOpen
	StateOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateHalfOpen:
		return "half-open"
	case StateOpen:
		return "open"
	default:
		return fmt.Sprintf("unknown state : %d",s)
	}
}

// 次数，请求量，失败量...

type Counts struct {
	// 总请求量
	Requests 				uint32
	// 总成功量
	TotalSuccesses 			uint32
	// 总失败量
	TotalFailures			uint32
	// 连续成功量
	ConsecutiveSuccesses	uint32
	// 连续失败量
	ConsecutiveFailures 	uint32
}

func (c *Counts) request() {
	c.Requests++
}
func (c *Counts) success() {
	c.TotalSuccesses++
	c.ConsecutiveSuccesses++
	// 注意，成功了就意味着失败次数就不是连续了
	c.ConsecutiveFailures = 0
}
func (c *Counts) failure() {
	c.TotalFailures++
	c.ConsecutiveFailures++
	// 注意，失败了就意味着成功次数就不是连续了
	c.ConsecutiveSuccesses = 0
}
func (c *Counts) clear() {
	c.Requests 				= 0
	c.TotalFailures			= 0
	c.TotalSuccesses		= 0
	c.ConsecutiveFailures 	= 0
	c.ConsecutiveSuccesses	= 0
}