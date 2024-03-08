package main

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

// redis的事务处理
// 1. 事务是一个单独的隔离的操作，有多个指令按顺序执行。不会被其他客户端打断
// 2. 事务是一个原子操作，事务中的命令要么全部执行，要么全不执行

func txPipeline() {

	// 开启事务
	pipeline := rdb.TxPipeline()
	// 输入指令
	pipeline.Incr(ctx, "value-1")
	pipeline.Expire(ctx, "value-1", time.Second)
	// 提交全部指令
	exec, _ := pipeline.Exec(ctx)
	fmt.Println(exec)

	//	-----上述代码等于----
	// MULTI
	// INCR value-1
	// EXPIRE　value-1
	// EXEC
}
func call(tx *redis.Tx) error {

	return nil
}
func watch() {
	// 重试3次
	for i := 0; i < 3; i++ {
		// watch 监听一个key，
		err := rdb.Watch(ctx, call, "key-1")
		if err == nil { // 操作成功
			break
		}
		if err == redis.TxFailedErr { // 操作失败 重试
			continue
		}
	}
}
