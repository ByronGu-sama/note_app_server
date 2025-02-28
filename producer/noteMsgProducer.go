package producer

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"note_app_server/config"
	"note_app_server/model/mqMessageModel"
	"time"
)

var partitions = config.AC.Kafka.NoteLikes.Partitions
var part = 0

func LikeNote(uid uint, nid string) error {
	network := config.AC.Kafka.Network
	host := config.AC.Kafka.Host
	port := config.AC.Kafka.Port
	topic := config.AC.Kafka.NoteLikes.Topic

	// 连接至Kafka集群的Leader节点
	conn, err := kafka.DialLeader(context.Background(), network, host+":"+port, topic, part)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	// 设置发送消息的超时时间
	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return err
	}

	msg := &mqMessageModel.LikeNotes{
		Uid:       uid,
		Nid:       nid,
		Action:    "likeNote",
		Timestamp: time.Now(),
	}

	marshal, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// 发送消息
	_, err = conn.WriteMessages(
		kafka.Message{Value: marshal},
	)
	if err != nil {
		return err
	}

	// 关闭连接
	if err = conn.Close(); err != nil {
		return err
	}
	return nil
}
