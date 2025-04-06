package producer

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"note_app_server/config/kafkaAction"
	"note_app_server/model/mqMessageModel"
	"note_app_server/model/msgModel"
	"note_app_server/producer/connManager"
	"time"
)

func SyncMessageToMongo(firstKey, secondKey uint, message *msgModel.Message) error {
	msg := &mqMessageModel.SyncMessageMsg{
		Action:    kafkaAction.SyncMessage,
		FirstKey:  firstKey,
		SecondKey: secondKey,
		Message:   message,
		Timestamp: time.Now(),
	}

	encodedMsg, err := msg.Encode()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 使用字典序更大的值作为key
	fk, tk := fmt.Sprintf("%d%d", firstKey, secondKey), fmt.Sprintf("%d%d", secondKey, firstKey)
	finallyKey := ""
	if fk > tk {
		finallyKey = fk
	} else {
		finallyKey = tk
	}

	// 发送消息
	err = connManager.SyncMessagesWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(finallyKey),
		Value: encodedMsg,
	})
	if err != nil {
		return err
	}
	return nil
}
