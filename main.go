package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 地址
		Password: "",               // 密码（无则留空）
		DB:       0,                // 默认数据库
		PoolSize: 10,               // 连接池大小（推荐值：CPU核心数*2+2）
	})

	// 连接测试
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	// 设置 5 秒超时
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	// 设置键值（带10秒过期）
	err := rdb.Set(ctx, "username", "alice", 10*time.Second).Err()
	if err != nil {
		log.Printf("设置失败: %v", err) // 网页7推荐用日志替代panic
	}

	// 获取值（区分键不存在与系统错误）
	val, err := rdb.Get(ctx, "username").Result()
	if errors.Is(err, redis.Nil) { // 网页4的错误处理规范
		fmt.Println("键不存在")
	} else if err != nil {
		panic(fmt.Sprintf("系统错误: %v", err))
	}
	fmt.Println("用户名:", val)

	// 设置哈希字段（支持批量）
	err1 := rdb.HSet(ctx, "user:1001", map[string]interface{}{
		"name": "Bob",
		"age":  30,
	}).Err()

	// 获取所有字段
	userData, err1 := rdb.HGetAll(ctx, "user:1001").Result()
	if err1 == nil {
		fmt.Printf("用户数据: %+v\n", userData) // map[name:Bob age:30]
	}

	// 删除单个键
	if err := rdb.Del(ctx, "username").Err(); err != nil {
		log.Printf("删除失败: %v", err)
	}

	// 批量删除（支持通配符）
	iter := rdb.Scan(ctx, 0, "temp:*", 100).Iterator()
	for iter.Next(ctx) {
		rdb.Del(ctx, iter.Val())
	}
}
