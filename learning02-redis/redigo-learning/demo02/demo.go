package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"reflect"
)

func main() {
	// 连接redis数据库 默认本地
	c, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		panic(err)
	}
	defer c.Close()

	// 批量发送
	c.Send("del", "test1", "test2")
	c.Send("sadd", "test1", "aa", "bb", "cc")
	c.Send("lpush", "test2", 10001, 10002, 10003)
	c.Send("smembers", "test1")
	c.Send("lrange", "test2", 0, -1)
	c.Flush()
	c.Receive()                           // del
	c.Receive()                           // sadd
	c.Receive()                           // lpush
	mbrs, _ := redis.Strings(c.Receive()) // smembers
	fmt.Println(mbrs, reflect.TypeOf(mbrs))
	lsts, _ := redis.Ints(c.Receive()) // lrange
	fmt.Println(lsts, reflect.TypeOf(lsts))
}
