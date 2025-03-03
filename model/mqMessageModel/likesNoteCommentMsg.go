package mqMessageModel

import (
	"encoding/json"
	"time"
)

// LikeNoteComment 点赞&取消点赞
type LikeNoteComment struct {
	Action    string
	Cid       string
	Uid       uint
	Timestamp time.Time
}

func (msg *LikeNoteComment) Decode(object []byte) error {
	return json.Unmarshal(object, msg)
}

func (msg *LikeNoteComment) Encode() ([]byte, error) {
	return json.Marshal(msg)
}
