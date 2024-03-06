package main

import "fmt"

func hSet() {

	// hash的key  字段  字段值
	rdb.HSet(ctx, "user", "userName", "liu")
	rdb.HSet(ctx, "user", "userAge", 20)
}
func hGet() {
	// 获取指定hash的指定字段值
	result, _ := rdb.HGet(ctx, "user", "userName").Result()
	fmt.Println(result)
}

func hGetAll() {
	// 获取指定hash的全部字段值
}

func hIncrBy() {
	//	 对字段的值进行累加
	rdb.HIncrBy(ctx, "user", "userAge", 2)
}

func hKeys() {

	// 获取指定hash内的全部字段
	keys := rdb.HKeys(ctx, "user")
	fmt.Println(keys)
}

func hLen() {

	// 返回hash 的字段数量
	rdb.HLen(ctx, "user")
}

func hMGet() {
	// 批量获取字段值
	hmGet := rdb.HMGet(ctx, "user", "userName", "userAge")
	fmt.Println(hmGet)

}

func hMSet() {
	data := make(map[string]interface{})
	data["id"] = 1
	data["name"] = "liu"
	// 批量设置字段值
	hmGet := rdb.HMSet(ctx, "user", data)
	fmt.Println(hmGet)
}

func hSetNx() {

	//	如果字段不存在 才设置值 （注意是字段）
}

func hDel() {

	// 删除字段
	rdb.HDel(ctx, "user", "userName", "userAge")
}

func hExists() {

	// 检查字段是否存在
	rdb.HExists(ctx, "user", "userName")
}
