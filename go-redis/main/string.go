package main

import (
	"context"
	"fmt"
	"time"
)

func set() {
	// 设置key-value
	rdb.Set(ctx, "key-1", "1111", 0)
}

func get() {
	// 获取key的值
	result, _ := rdb.Get(ctx, "key-1").Result()
	fmt.Println(result)
}

func getSet() {

	// 设置值并返回旧值
	stringCmd := rdb.GetSet(ctx, "key-1", "111")
	fmt.Println(stringCmd)
}

func setNx() {
	// 如果key不存在 才会设置值
	rdb.SetNX(ctx, "key-1", "111", 0)
	rdb.SetNX(ctx, "key-2", "222", 0)
}

func mGet() {
	// 批量获取
	sliceCmd := rdb.MGet(ctx, "key-1", "key-2")
	fmt.Println(sliceCmd)
}

func mSet() {
	// 批量设置值
	rdb.MSet(ctx, "key-3", "333", "key-4", "444")
}

func Incr() {
	// 对指定key自增 1
	result, _ := rdb.Incr(ctx, "key-1").Result()
	fmt.Println(result)
}
func incrBy() {
	// 自增指定数值
	result, _ := rdb.IncrBy(ctx, "key-1", 2).Result()
	fmt.Println(result)
}
func incrByFloat() {
	// 自增指定数值
	result, _ := rdb.IncrByFloat(ctx, "key-1", 2.2).Result()
	fmt.Println(result)
}
func decr() {
	//自减
}

func decrBy() {

}
func del() {
	// 删除key 支持多个
	rdb.Del(context.Background(), "key1", "key2")
}

func expire() {
	// 设置过期时间
	rdb.Expire(ctx, "key1", 3*time.Second)
}
