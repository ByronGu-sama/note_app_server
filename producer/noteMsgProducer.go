package producer

import (
	"github.com/segmentio/kafka-go"
	"note_app_server/config/kafkaAction"
	"note_app_server/model/mqMessageModel"
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

	// 发送消息
	_, err = connManager.NoteLikesMQConn.WriteMessages(
		kafka.Message{Value: encodedMsg},
	)
	if err != nil {
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

	// 发送消息
	_, err = connManager.NoteLikesMQConn.WriteMessages(
		kafka.Message{Value: encodedMsg},
	)
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

	// 发送消息
	_, err = connManager.NoteCollectsMQConn.WriteMessages(
		kafka.Message{Value: encodedMsg},
	)
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

	// 发送消息
	_, err = connManager.NoteCollectsMQConn.WriteMessages(
		kafka.Message{Value: encodedMsg},
	)
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

	// 发送消息
	_, err = connManager.NoteCommentsMQConn.WriteMessages(
		kafka.Message{Value: encodedMsg},
	)
	if err != nil {
		return err
	}
	return nil
}
