package main

import "fmt"

// 发布订阅案例

func subscribe() {
	// 订阅主题
	sub := rdb.Subscribe(ctx, "channel-1", "channel-2")
	//sub = rdb.PSubscribe(ctx, "channel-*") 带通配符匹配主题
	// 接收主题消息
	for ch := range sub.Channel() {
		fmt.Println(ch.Channel)
		fmt.Println(ch.Payload)
	}
}

func unSub() {
	sub := rdb.Subscribe(ctx, "channel-1", "channel-2")
	//	取消订阅
	err := sub.Unsubscribe(ctx, "channel-1", "channel-2")
	if err != nil {
		return
	}
}

func publish() {
	// 发布消息
	rdb.Publish(ctx, "channel-1", "msg...")
}
