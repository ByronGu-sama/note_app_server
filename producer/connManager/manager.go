package connManager

import "C"
import (
	"github.com/segmentio/kafka-go"
	"note_app_server/config"
	"sync"
)

var (
	NoteLikesWriter    *kafka.Writer
	NoteCollectsWriter *kafka.Writer
	NoteCommentsWriter *kafka.Writer
	SyncNotesWriter    *kafka.Writer
	DelNotesWriter     *kafka.Writer
	SyncMessagesWriter *kafka.Writer
)

func InitKafkaConn() {
	wg := sync.WaitGroup{}
	wg.Add(6)
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
	go func() {
		defer wg.Done()
		initSyncNoteMQConn()
	}()
	go func() {
		defer wg.Done()
		initDelNoteMQConn()
	}()
	go func() {
		defer wg.Done()
		initSyncMessageMQConn()
	}()
	wg.Wait()
}

// 点赞笔记相关的连接
func initNoteLikesMQConn() {
	NoteLikesWriter = &kafka.Writer{
		Addr:     kafka.TCP(config.AC.Kafka.Addr),
		Topic:    config.AC.Kafka.NoteLikes.Topic,
		Balancer: &kafka.LeastBytes{},
	}
}

// 收藏笔记相关的连接
func initNoteCollectsMQConn() {
	NoteCollectsWriter = &kafka.Writer{
		Addr:     kafka.TCP(config.AC.Kafka.Addr),
		Topic:    config.AC.Kafka.NoteCollects.Topic,
		Balancer: &kafka.LeastBytes{},
	}
}

// 评论笔记相关的连接
func initCommentMQConn() {
	NoteCommentsWriter = &kafka.Writer{
		Addr:     kafka.TCP(config.AC.Kafka.Addr),
		Topic:    config.AC.Kafka.NoteComments.Topic,
		Balancer: &kafka.LeastBytes{},
	}
}

// 同步笔记相关的连接
func initSyncNoteMQConn() {
	SyncNotesWriter = &kafka.Writer{
		Addr:     kafka.TCP(config.AC.Kafka.Addr),
		Topic:    config.AC.Kafka.SyncNotes.Topic,
		Balancer: &kafka.LeastBytes{},
	}
}

// 删除笔记相关的连接
func initDelNoteMQConn() {
	DelNotesWriter = &kafka.Writer{
		Addr:     kafka.TCP(config.AC.Kafka.Addr),
		Topic:    config.AC.Kafka.DelNotes.Topic,
		Balancer: &kafka.LeastBytes{},
	}
}

// 同步聊天消息的连接
func initSyncMessageMQConn() {
	SyncMessagesWriter = &kafka.Writer{
		Addr:     kafka.TCP(config.AC.Kafka.Addr),
		Topic:    config.AC.Kafka.SyncMessages.Topic,
		Balancer: &kafka.LeastBytes{},
	}
}
