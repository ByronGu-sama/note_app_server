package mqMessageModel

import (
	"encoding/json"
	"time"
)

// CollectNotes 收藏&取消收藏笔记
type CollectNotes struct {
	Action    string
	Nid       string
	Uid       int64
	Timestamp time.Time
}

func (msg *CollectNotes) Decode(object []byte) error {
	return json.Unmarshal(object, msg)
}

func (msg *CollectNotes) Encode() ([]byte, error) {
	return json.Marshal(msg)
}
