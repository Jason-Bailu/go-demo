# Redigo数据库操作

- Redigo 模块安装：go get github.com/gomodule/redigo/redis
- 相关文档：[RedigoDoc](https://pkg.go.dev/github.com/gomodule/redigo/redis#pkg-overview)

## redis连接

```go
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
	// 也可以通过用户名，密码来连接
	//c, err := redis.Dial("tcp", "localhost:6379",
	//	redis.DialUsername("username"),
	//	redis.DialPassword("password"),
	//)
	// context上下文连接
	//ctx := context.Background()
	//c, err := redis.DialContext(ctx, "tcp", ":6379")
	// 使用 URL 连接到 Redis 的远程实例。
	//c, err := redis.DialURL(os.Getenv("REDIS_URL"))
}
```

- 通过dial创建了redis客户端实例，进行交互

## redis常用操作

```go
// 执行命令 通用方法Do("command", "key", "value")
	// 不同的返回类型通过通过redis类型转换
	// integer                 int64
	// simple string           string
	// bulk string             []byte or nil if value not present.
	// array                   []interface{} or nil if value not present.
	// key 操作
	// set key
	reply, _ := c.Do("set", "redigo:key", 2)
	fmt.Println(reply) // OK
	// get key
	result, _ := redis.Int(c.Do("get", "redigo:key"))
	fmt.Println(result, reflect.TypeOf(result))
	// remove key
	reply, _ = c.Do("del", "redigo:key")
	fmt.Println(reply)

	// list 操作
	// 添加
	args := redis.Args{"test_list"}.Add("test1").Add("test2")
	c.Do("lpush", args...)
	// 获取 []string
	resultStrs, _ := redis.Strings(c.Do("lrange", "test_list", 0, -1))
	fmt.Println(resultStrs, reflect.TypeOf(resultStrs))
	// 删除
	c.Do("del", "test_list")
	// 不同变量的读取 scan方法进行读取
	c.Do("lpush", "list", "aaabb", 100)
	vals, _ := redis.Values(c.Do("lrange", "list", 0, -1))
	var name string
	var score int
	redis.Scan(vals, &score, &name)
	fmt.Println(name, ":", score)
	c.Do("del", "list")

	// struct 操作
	var p1, p2 struct {
		Name string `redis:"name"`
		Age  string `redis:"age"`
		Sex  string `redis:"sex"`
	}
	p1.Age = "18"
	p1.Name = "chaochaoyu"
	p1.Sex = "male"
	// 结构体传参
	args = redis.Args{}.Add("role:test1").AddFlat(&p1)
	c.Do("hmset", args...)
	// map传参
	m := map[string]string{
		"name": "quxiansen",
		"age":  "20",
		"sex":  "female",
	}
	args = redis.Args{}.Add("role:test2").AddFlat(m)
	c.Do("hmset", args...)
	// 读取结构体 通过ScanStruct读取
	for _, id := range []string{"role:test1", "role:test2"} {
		// HGETALL id, 获取当前key和value值
		v, _ := redis.Values(c.Do("HGETALL", id))
		_ = redis.ScanStruct(v, &p2)
		fmt.Printf("%+v\n", p2)
	}
	c.Do("del", "role:test1", "role:test2")
```

## 管道：一次发送多条命令，减少与 redis-server 之间的网络交互

```go
// 批量发送，批量接收
c.Send(cmd1, ...)
c.Send(cmd2, ...)
c.Send(cmd3, ...)
c.Flush() // 将上面的三个命令发送出去
c.Receive() // cmd1 的返回值
c.Receive() // cmd2 的返回值
c.Receive() // cmd3 的返回值
// 如果不需要关注返回值
c.Send(cmd1, ...)
c.Send(cmd2, ...)
c.Send(cmd3, ...)
c.Do("")
// 如果只关注最后一个命令的返回值
c.Send(cmd1, ...)
c.Send(cmd2, ...)
c.Do(cmd3, ...)
```

---

```go
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
```

## 事务：对于一条链接请求队列是线性执行的

```go
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
```

---

- redis的事务具有**隔离性**，但是不具备**原子性、持久化、一致性**

## 发布订阅：

```go
package main

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

// listenPubSubChannels listens for messages on Redis pubsub channels. The
// onStart function is called after the channels are subscribed. The onMessage
// function is called for each message.
func listenPubSubChannels(ctx context.Context, redisServerAddr string,
	onStart func() error,
	onMessage func(channel string, data []byte) error,
	channels ...string) error {
	// A ping is set to the server with this period to test for the health of
	// the connection and server.
	const healthCheckPeriod = time.Minute

	c, err := redis.Dial("tcp", redisServerAddr,
		// Read timeout on server should be greater than ping period.
		redis.DialReadTimeout(healthCheckPeriod+10*time.Second),
		redis.DialWriteTimeout(10*time.Second))
	if err != nil {
		return err
	}
	defer c.Close()

	psc := redis.PubSubConn{Conn: c}

	if err := psc.Subscribe(redis.Args{}.AddFlat(channels)...); err != nil {
		return err
	}

	done := make(chan error, 1)

	// Start a goroutine to receive notifications from the server.
	go func() {
		for {
			switch n := psc.Receive().(type) {
			case error:
				done <- n
				return
			case redis.Message:
				if err := onMessage(n.Channel, n.Data); err != nil {
					done <- err
					return
				}
			case redis.Subscription:
				switch n.Count {
				case len(channels):
					// Notify application when all channels are subscribed.
					if err := onStart(); err != nil {
						done <- err
						return
					}
				case 0:
					// Return from the goroutine when all channels are unsubscribed.
					done <- nil
					return
				}
			}
		}
	}()

	ticker := time.NewTicker(healthCheckPeriod)
	defer ticker.Stop()
loop:
	for {
		select {
		case <-ticker.C:
			// Send ping to test health of connection and server. If
			// corresponding pong is not received, then receive on the
			// connection will timeout and the receive goroutine will exit.
			if err = psc.Ping(""); err != nil {
				break loop
			}
		case <-ctx.Done():
			break loop
		case err := <-done:
			// Return error from the receive goroutine.
			return err
		}
	}

	// Signal the receiving goroutine to exit by unsubscribing from all channels.
	if err := psc.Unsubscribe(); err != nil {
		return err
	}

	// Wait for goroutine to complete.
	return <-done
}

func publish() {
	c, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	if _, err = c.Do("PUBLISH", "c1", "hello"); err != nil {
		fmt.Println(err)
		return
	}
	if _, err = c.Do("PUBLISH", "c2", "world"); err != nil {
		fmt.Println(err)
		return
	}
	if _, err = c.Do("PUBLISH", "c1", "goodbye"); err != nil {
		fmt.Println(err)
		return
	}
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	err := listenPubSubChannels(ctx,
		"localhost:6379",
		func() error {
			// The start callback is a good place to backfill missed
			// notifications. For the purpose of this example, a goroutine is
			// started to send notifications.
			// 这里发布了一些信息
			go publish()
			return nil
		},
		func(channel string, message []byte) error {
			fmt.Printf("channel: %s, message: %s\n", channel, message)
			// For the purpose of this example, cancel the listener's context
			// after receiving last message sent by publish().
			if string(message) == "goodbye" {
				cancel()
			}
			return nil
		},
		"c1", "c2")

	if err != nil {
		fmt.Println(err)
		return
	}
}
```

