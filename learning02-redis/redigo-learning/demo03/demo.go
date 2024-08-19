package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

// zpop pops a value from the ZSET key using WATCH/MULTI/EXEC commands.
func zpop(c redis.Conn, key string) (result string, err error) {

	defer func() {
		// Return connection to normal state on error.
		if err != nil {
			c.Do("DISCARD") // nolint: errcheck
		}
	}()

	// Loop until transaction is successful.
	for {
		// 监听key变动 乐观锁，key在事务之前发生变动则取消事务
		if _, err := c.Do("WATCH", key); err != nil {
			return "", err
		}

		members, err := redis.Strings(c.Do("ZRANGE", key, 0, 0))
		if err != nil {
			return "", err
		}
		if len(members) != 1 {
			return "", redis.ErrNil
		}
		// 开启事务
		if err = c.Send("MULTI"); err != nil {
			return "", err
		}
		if err = c.Send("ZREM", key, members[0]); err != nil {
			return "", err
		}
		// 提交事务
		queued, err := c.Do("EXEC")
		if err != nil {
			return "", err
		}

		if queued != nil {
			result = members[0]
			break
		}
	}

	return result, nil
}

// zpopScript 脚本
var zpopScript = redis.NewScript(1, `
    local r = redis.call('ZRANGE', KEYS[1], 0, 0)
    if r ~= nil then
        r = r[1]
        redis.call('ZREM', KEYS[1], r)
    end
    return r
`)

func main() {
	// 连接redis数据库 默认本地
	c, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		panic(err)
	}
	defer c.Close()

	// 管道批量传输指令
	for i, member := range []string{"red", "blue", "green"} {
		if err = c.Send("ZADD", "zset", i, member); err != nil {
			fmt.Println(err)
			return
		}
	}
	// 发送管道中的指令
	if _, err := c.Do(""); err != nil {
		fmt.Println(err)
		return
	}

	// Pop using WATCH/MULTI/EXEC
	// 使用事务进行pop
	v, err := zpop(c, "zset")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(v)

	// Pop using a script.
	// 使用脚本进行pop
	v, err = redis.String(zpopScript.Do(c, "zset"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(v)
}
