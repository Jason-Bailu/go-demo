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

	// struct结构体
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
}
