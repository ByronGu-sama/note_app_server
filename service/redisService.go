package service

import (
	"context"
	"fmt"
	"note_app_server/global"
)

// Publish 发布消息
func Publish(c context.Context, channel string, msg string) error {
	return global.MsgRdb.Publish(c, channel, msg).Err()
}

// Subscribe 订阅
func Subscribe(c context.Context, channel string) (string, error) {
	sub := global.MsgRdb.Subscribe(c, channel)
	defer sub.Close()
	fmt.Print("Subscribe...", sub)
	msg, err := sub.ReceiveMessage(c)
	if err != nil {
		return "", err
	}
	return msg.Payload, err
}
