package connManager

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"note_app_server/config"
	"sync"
	"time"
)

var (
	NoteLikesMQConn    *kafka.Conn
	NoteCollectsMQConn *kafka.Conn
	NoteCommentsMQConn *kafka.Conn
)

func InitKafkaConn() {
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		initNoteLikesMQConn()
	}()
	go func() {
		defer wg.Done()
		initNoteCollectsMQConn()
	}()
	go func() {
		defer wg.Done()
		initCommentMQConn()
	}()
	wg.Wait()
}

// 点赞笔记相关的连接
func initNoteLikesMQConn() {
	conn, err := kafka.DialLeader(context.Background(), config.AC.Kafka.Network, config.AC.Kafka.Host+":"+config.AC.Kafka.Port, config.AC.Kafka.NoteLikes.Topic, 0)
	if err != nil {
		// 加入重试机制
		log.Fatal("failed to dial leader:", err)
	}
	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	NoteLikesMQConn = conn
}

// 收藏笔记相关的连接
func initNoteCollectsMQConn() {
	conn, err := kafka.DialLeader(context.Background(), config.AC.Kafka.Network, config.AC.Kafka.Host+":"+config.AC.Kafka.Port, config.AC.Kafka.NoteCollects.Topic, 0)
	if err != nil {
		// 加入重试机制
		log.Fatal("failed to dial leader:", err)
	}
	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	NoteCollectsMQConn = conn
}

// 评论笔记相关的连接
func initCommentMQConn() {
	conn, err := kafka.DialLeader(context.Background(), config.AC.Kafka.Network, config.AC.Kafka.Host+":"+config.AC.Kafka.Port, config.AC.Kafka.NoteComments.Topic, 0)
	if err != nil {
		// 加入重试机制
		log.Fatal("failed to dial leader:", err)
	}
	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	NoteCommentsMQConn = conn
}
