package mqMessageModel

import (
	"encoding/json"
	"time"
)

type DelNoteComment struct {
	Action    int
	Cid       string
	Uid       int64
	Timestamp time.Time
}

func (msg *DelNoteComment) Decode(object []byte) error {
	return json.Unmarshal(object, msg)
}

func (msg *DelNoteComment) Encode() ([]byte, error) {
	return json.Marshal(msg)
}
