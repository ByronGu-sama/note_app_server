package producer

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"note_app_server/config/kafkaAction"
	"note_app_server/model/mqMessageModel"
	"note_app_server/model/noteModel"
	"note_app_server/producer/connManager"
	"time"
)

// LikeNote 点赞笔记
func LikeNote(uid uint, nid string) error {
	msg := &mqMessageModel.LikeNotes{
		Uid:       uid,
		Nid:       nid,
		Action:    kafkaAction.LikeNote,
		Timestamp: time.Now(),
	}

	encodedMsg, err := msg.Encode()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 发送消息
	err = connManager.NoteLikesWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(nid),
		Value: encodedMsg,
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// DislikeNote 取消点赞笔记
func DislikeNote(uid uint, nid string) error {
	msg := &mqMessageModel.LikeNotes{
		Uid:       uid,
		Nid:       nid,
		Action:    kafkaAction.DislikeNote,
		Timestamp: time.Now(),
	}

	encodedMsg, err := msg.Encode()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 发送消息
	err = connManager.NoteLikesWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(nid),
		Value: encodedMsg,
	})
	if err != nil {
		return err
	}
	return nil
}

// CollectNote 收藏笔记
func CollectNote(uid uint, nid string) error {
	msg := &mqMessageModel.CollectNotes{
		Uid:       uid,
		Nid:       nid,
		Action:    kafkaAction.CollectNote,
		Timestamp: time.Now(),
	}

	encodedMsg, err := msg.Encode()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 发送消息
	err = connManager.NoteCollectsWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(nid),
		Value: encodedMsg,
	})
	if err != nil {
		return err
	}
	return nil
}

// AbandonNote 取消收藏笔记
func AbandonNote(uid uint, nid string) error {
	msg := &mqMessageModel.CollectNotes{
		Uid:       uid,
		Nid:       nid,
		Action:    kafkaAction.AbandonNote,
		Timestamp: time.Now(),
	}

	encodedMsg, err := msg.Encode()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 发送消息
	err = connManager.NoteCollectsWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(nid),
		Value: encodedMsg,
	})
	if err != nil {
		return err
	}
	return nil
}

// DelComment 删除笔记评论
func DelComment(cid string, uid uint) error {
	msg := &mqMessageModel.DelNoteComment{
		Cid:       cid,
		Uid:       uid,
		Action:    kafkaAction.DelNoteComment,
		Timestamp: time.Now(),
	}

	encodedMsg, err := msg.Encode()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 发送消息
	err = connManager.DelNotesWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(cid),
		Value: encodedMsg,
	})
	if err != nil {
		return err
	}
	return nil
}

// LikeNoteComment 点赞笔记评论
func LikeNoteComment(uid uint, cid string) error {
	msg := &mqMessageModel.LikeNoteComment{
		Cid:       cid,
		Uid:       uid,
		Action:    kafkaAction.LikeComment,
		Timestamp: time.Now(),
	}

	encodedMsg, err := msg.Encode()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 发送消息
	err = connManager.NoteCommentsWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(cid),
		Value: encodedMsg,
	})
	if err != nil {
		return err
	}
	return nil
}

// DislikeNoteComment 取消点赞笔记评论
func DislikeNoteComment(uid uint, cid string) error {
	msg := &mqMessageModel.LikeNoteComment{
		Cid:       cid,
		Uid:       uid,
		Action:    kafkaAction.DislikeComment,
		Timestamp: time.Now(),
	}

	encodedMsg, err := msg.Encode()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 发送消息
	err = connManager.NoteCommentsWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(cid),
		Value: encodedMsg,
	})
	if err != nil {
		return err
	}
	return nil
}

// DelNote 删除笔记
func DelNote(uid uint, nid string) error {
	msg := &mqMessageModel.DelNote{
		Nid:       nid,
		Uid:       uid,
		Action:    kafkaAction.DelNote,
		Timestamp: time.Now(),
	}

	encodedMsg, err := msg.Encode()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 发送消息
	err = connManager.DelNotesWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(nid),
		Value: encodedMsg,
	})
	if err != nil {
		return err
	}
	return nil
}

// SyncToES 同步笔记
func SyncToES(note *noteModel.ESNote) error {
	msg := &mqMessageModel.SyncNoteMsg{
		Action:    kafkaAction.SyncNote,
		Note:      note,
		Timestamp: time.Now(),
	}

	encodedMsg, err := msg.Encode()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 发送消息
	err = connManager.SyncNotesWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(note.Nid),
		Value: encodedMsg,
	})
	if err != nil {
		return err
	}
	return nil
}
