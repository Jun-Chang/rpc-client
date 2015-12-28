package service

import (
	"github.com/CyberAgent/car-golib/cache"
)

func Run() int {
	obj := []interface{}{"Seq", 1}
	seq, _ := cache.GetConn().Increment(Key(), obj, 0)

	return seq
}

func Key() []interface{} {
	return []interface{}{"rpc_test", "test"}
}
