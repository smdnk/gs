package main

func main() {

	s := rdb.Ping(ctx).String()

	// getKey
	val, _ := rdb.Get(ctx, "11").Result()
	// setKey
	rdb.Set(ctx, "key-1", "value", 0) // 0不过期

	// Do 可以直接执行redis原生命令
	result, _ := rdb.Do(ctx, "get", "key-1").Result()

	print(val, s, result)
}
