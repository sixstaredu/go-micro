package idempotency

import (
	"github.com/garyburd/redigo/redis"
	"github.com/sixstaredu/go-micro/micro"
	"time"
)

var RedisIdempotent = new(redisIdempotent)

type redisIdempotent struct {}

func (i *redisIdempotent) TryAcquire(id string, timeout time.Duration) bool {
	rc := micro.Redis.Get()
	defer rc.Close()

	// 设置id，
	reply, err := redis.String(rc.Do("set", id, time.Now().Add(timeout).Format(TimeFormat), "ex", "86400", "nx"))
	if err != nil {
		return false
	}
	if reply == "OK" {
		// 设置成功
		return true
	}
	// 判断是否有OK

	reply, err = redis.String(rc.Do("get", id))
	// 1. 判断是否完成
	if reply == Confirmed {
		return false
	}

	replyt,err := time.ParseInLocation(TimeFormat, reply, time.Local)
	if err != nil {
		return false
	}
	if !replyt.After(time.Now()) {
		// 未超时； 说明可能存在其他程序
		return false
	}

	// 需要注意可能当前环节会存在多个进程执行
	delta := time.Now().Add(timeout).Sub(replyt)
	reply, err = redis.String(rc.Do("incrby", id, delta.String()))
	if err != nil {
		return false
	}
	newreplyt,err := time.ParseInLocation(TimeFormat, reply, time.Local)
	if err != nil {
		return false
	}

	if newreplyt.Equal(replyt.Add(delta)) {
		return true
	} else {
		rc.Do("decrby", id, delta.String())
	}


	return false
}

func (i *redisIdempotent) Comfirm(id string) {
	rc := micro.Redis.Get()
	defer rc.Close()

	rc.Do("set", id, Confirmed, "xx")
}