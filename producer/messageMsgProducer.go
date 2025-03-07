package producer

import (
	"github.com/segmentio/kafka-go"
	"log"
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

	err = connManager.SyncMessagesMQConn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		log.Fatal("failed to set timeout:", err)
	}

	// 发送消息
	_, err = connManager.SyncMessagesMQConn.WriteMessages(
		kafka.Message{Value: encodedMsg},
	)
	if err != nil {
		return err
	}
	return nil
}
