package tool

import (
	"github.com/sony/sonyflake"
)

var flake *sonyflake.Sonyflake

func init() {
	flake = sonyflake.NewSonyflake(sonyflake.Settings{})
}

func GenId() int64 {

	id, err := flake.NextID()
	if err != nil {
		panic(err)
	}

	return int64(id)
}
