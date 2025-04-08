package mqMessageModel

import (
	"encoding/json"
	"time"
)

// LikeNotes 点赞&取消点赞笔记
type LikeNotes struct {
	Action    string
	Nid       string
	Uid       int64
	Timestamp time.Time
}

func (msg *LikeNotes) Decode(object []byte) error {
	return json.Unmarshal(object, msg)
}

func (msg *LikeNotes) Encode() ([]byte, error) {
	return json.Marshal(msg)
}
